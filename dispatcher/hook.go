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
 * hook.go: Logrus hook plugin
 */

package dispatcher

import "github.com/sirupsen/logrus"

type DispatcherHook struct {
	d *Dispatcher
}

func (h DispatcherHook) Levels() []logrus.Level  {
	return logrus.AllLevels
}

func (h DispatcherHook) Fire(entry *logrus.Entry) error {
	h.d.q<-*entry
	return nil
}