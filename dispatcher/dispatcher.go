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
 * dispatcher.go: Log entry dispatcher and selector applier
 */

package dispatcher

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Dispatcher struct {
	q chan logrus.Entry
	stop chan bool
	wg *sync.WaitGroup
	history [logBuffer]logrus.Entry
	selectors []*Selector
	register chan *Selector
	unregister chan *Selector
	ptr int
	curr int
}

const chanBuffer = 64
const logBuffer = 64

func NewDispatcher() (ret *Dispatcher) {
	ret = &Dispatcher{}
	ret.q = make(chan logrus.Entry, chanBuffer)
	ret.stop = make(chan bool, 0)
	ret.wg = &sync.WaitGroup{}
	ret.history = [logBuffer]logrus.Entry{}
	ret.selectors = make([]*Selector, 0)
	ret.register = make(chan *Selector, 0)
	ret.unregister = make(chan *Selector, 0)
	ret.ptr = 0
	ret.curr = 0
	ret.history[ret.ptr] = logrus.Entry{Data:logrus.Fields{"prefix": "weblog-dispatcher", "event": "started"}, Message: "Weblog dispatcher started", Level:logrus.InfoLevel, Time: time.Now()}
	go ret.dispatch()
	return
}

func (r *Dispatcher) dispatch() {
	defer r.wg.Done()
	r.wg.Add(1)
	run := true

	for run {
		select {
		case s:=<-r.unregister:
			fmt.Printf("dispatcher: unregistered a selector")
			sel := make([]*Selector, 0)
			for _, c := range r.selectors {
				if s == c {
					continue
				}
				sel = append(sel, c)
			}
			r.selectors = sel
			break
		case <-r.stop:
			run = false
			break
		case e := <- r.q:
//			fmt.Printf("dispatcher: got a log entry\n")
			r.curr += 1
			r.curr %= logBuffer
			if r.curr == r.ptr {
				r.ptr += 1
				r.ptr %= logBuffer
			}
			r.history[r.curr] = e
			for _, s := range r.selectors {
				// Try to write the queue, but if it's full, don't deadlock
				select {
				case s.q<- e:
					break
				default:
					break
				}
			}
			break
		case s := <- r.register:
			fmt.Printf("dispatcher: registered a selector\n")
			r.selectors = append(r.selectors, s)
			p := r.ptr
			for p != r.curr {
				select {
				case s.q <- r.history[p]:
					break
				default:
					break
				}
				p += 1
				p %= logBuffer
			}
			// Add the current one as well, or else selectors catching up miss out on one entry
			select {
			case s.q <- r.history[r.curr]:
				break
			default:
				break
			}
			break
		}
	}
}

func (r *Dispatcher) Register(s *Selector) {
	r.register<-s
}

func (r *Dispatcher) Unregister(s *Selector) {
	r.unregister<-s
}

func (r *Dispatcher) Hook() logrus.Hook {
	return DispatcherHook{r}
}

func (r *Dispatcher) Stop() {
	defer r.wg.Wait()
	for _, s := range r.selectors {
		r.Unregister(s)
	}
	r.stop <- true
}
