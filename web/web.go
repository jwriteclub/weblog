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
 * web.go: Server setup hook
 */

package web

import (
	"github.com/gorilla/mux"
	"github.com/jwriteclub/weblog/dispatcher"
	"github.com/markbates/pkger"
	"net/http"
)

func Setup(prefix string, rtr *mux.Router, d *dispatcher.Dispatcher) {
	rtr = rtr.PathPrefix(prefix).Subrouter()
	rtr.HandleFunc("/socket", NewWeblogHandler(d)).Name("weblog_socket")
	rtr.PathPrefix("/").Handler(http.StripPrefix(prefix, http.FileServer(pkger.Dir("/web/weblog-html")))).Name("weblog_panel")
}