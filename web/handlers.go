/**
 * Weblog
 *
 *    Copyright 2019 Christopher O'Connell
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * For any questions, please contact jwriteclub@gmail.com
 *
 * handlers.go: HTTP (+WebSocket) Handler
 */

package web

import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/jwriteclub/weblog/dispatcher"
	"fmt"
	"strings"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader {
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func NewWeblogHandler(d *dispatcher.Dispatcher) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("weblog: Got error %s\n", err.Error())
			return
		}

		run := true

		hello := make(map[string]string)
		hello["hello"] = "world"
		conn.WriteJSON(hello)

		s, err := dispatcher.NewSelector("", d)
		newSelector := true
		defer s.Stop()
		if err != nil {
			fmt.Printf("weblog: error creating selector: %s\n", err.Error())
			dat := make(map[string]string)
			dat["type"] = "error"
			dat["message"] = "unable to parse selector"
			dat["error"] = err.Error()
			conn.WriteJSON(dat)
			_ = conn.Close()
			return
		}

		mutex := &sync.RWMutex{}

		go func () {
			for run {
				mp := make(map[string]string)
				err := conn.ReadJSON(&mp)
				if err != nil {
					fmt.Printf("Got an error from the read channel: %s\n", err.Error())
					run = false
					continue
				}
				fmt.Printf("%#v\n", mp)
				if t, ok := mp["type"]; ok && t == "selector" {
					if sel, ok := mp["selector"]; ok {
						selector, err := dispatcher.NewSelector(sel, d)
						if err != nil {
							fmt.Printf("weblog: error creating selector: %s\n", err.Error())
							continue
						}
						mutex.Lock()
						s.Stop()
						s = selector
						newSelector = true
						mutex.Unlock()
					}
				}
			}
		}()

		for run {
			didSomething := false
			mutex.Lock()
			if newSelector {
				err := conn.WriteJSON(map[string]interface{}{"type": "clear"})
				if err != nil {
					fmt.Printf("weblog: Got error %s\n", err.Error())
					goto nsdone
				}
				err = conn.WriteJSON(map[string]interface{}{"type": "basetime", "basetime": int64(s.BaseTime().UnixNano() / 1000000)})
				if err != nil {
					fmt.Printf("weblog: Got error %s\n", err.Error())
					goto nsdone
				}
				newSelector = false
				didSomething = true
			}
		nsdone:
			mutex.Unlock()
			if !run {
				continue
			}

			e := s.MaybeRead()
			if e != nil {
				dat := make(map[string]interface{})
				dat["type"] = "log"
				dat["log"] = map[string]interface{} {
					"time": int64(e.Time.UnixNano() / 1000000),
					"level": strings.ToLower(e.Level.String()),
					"fields": e.Data,
					"message": fmt.Sprintf("%v\n", e.Message),
				}
				err := conn.WriteJSON(dat)
				if err != nil {
					fmt.Printf("weblog: Got error %s\n", err.Error())
					run = false
					continue
				}
				didSomething = false
			}

			if run && !didSomething {
				time.Sleep(time.Millisecond*10)
			}
		}
	}
}