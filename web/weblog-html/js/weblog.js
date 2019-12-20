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
 * weblog.js: Runtime javascript
 */

$(document).ready(function() {
   console.log("Setting up websocket connection");
   var socket = new WebSocket(location.href.replace('http://', 'ws://').replace('https://', 'wss://') + 'socket');
   var output =  $("#weblog-log-output");
   var disconnectedMessage = $("#weblog-message-disconnected");
   var basetime = new Date();
   var body = $("body");
   var search = $("#weblog-button-search").click(function () {
       var msg = "{\"type\": \"selector\", \"selector\": \""+$("#weblog-input-query").val()+"\"}";
       console.log(msg);
       socket.send(msg);
   });
   socket.onopen = function(ev) {
       console.log("WS opened");
       body.removeClass("weblog-connecting").removeClass("weblog-disconnected").addClass("weblog-connected");
       socket.onmessage = function(msg){
           var data = $.parseJSON(msg.data);
           if (typeof(data.type) === "undefined") {
               console.log("Invalid data");
               console.log(data);
               return
           }
           switch (data.type) {
               case "clear":
                   console.log("clearing");
                   output.empty();
                   break;
               case "basetime":
                   console.log("Setting basetime");
                   basetime = new Date(data.basetime);
                   console.log(basetime);
                   break;
               case "log":
                   var line = $("<tr class='weblog-log-line weblog-level-" + data.log.level + "'></tr>");
                   line.append("<td class='weblog-log-line-date-offset'>[" + Math.round(((new Date(data.log.time)) - basetime ) / 1000) + "]</td>");
                   line.append("<td class='weblog-log-line-level'>" + data.log.level + "</td>");
                   if (typeof(data.log.fields.prefix) !== "undefined") {
                       line.append("<td class='weblog-log-line-prefix'>" + data.log.fields.prefix + "</td>");
                   } else {
                       ling.append("<td class='weblog-log-line-prefix weblog-log-line-prefix-none'>[none]</td>")
                   }
                   line.append("<td class='weblog-log-line-message'>" + data.log.message + "</td>");
                   $("#weblog-log-output").prepend(line);

           }
       };
    };
    socket.onerror = function(ev) {
       console.log("WS error");
       console.log(ev);
       body.removeClass("weblog-connecting").removeClass("weblog-connected").addClass("weblog-disconnected");
    };
    socket.onclose = function(ev) {
        console.log("WS closed");
        console.log(ev);
        body.removeClass("weblog-connecting").removeClass("weblog-connected").addClass("weblog-disconnected");
    };
    window.shutdownSocket = function() {
        socket.close();
    }
});