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
 * selector.go: Predicate applier
 */

package dispatcher

import (
	"errors"
	"fmt"
	"github.com/jwriteclub/weblog/predicate"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var baseTimestamp = time.Now()

type Selector struct {
	q chan logrus.Entry
	predicate predicate.BoolOp
	d *Dispatcher
	t *time.Ticker
	m *sync.RWMutex
	reg bool
}

func NewSelector(expression string, dispatcher *Dispatcher) (ret *Selector, err error) {
	ret = &Selector{}
	ret.q = make(chan logrus.Entry, chanBuffer)
	ret.d = dispatcher
	ret.m = &sync.RWMutex{}
	ret.d.Register(ret)
	ret.reg = true
	err = ret.Select(expression)
	return
}

func (s *Selector) Select(expression string) (err error) {
	var op interface{}
	op, err = predicate.Parse("selector", []byte(expression))
	if err != nil {
		return
	}
	if _, ok := op.(predicate.BoolOp); !ok {
		err = errors.New("Invalid BoolOp from predicate")
		return
	}
	defer s.m.Unlock()
	s.m.Lock()
	fmt.Printf("selector: %#v\n", op)
	p, ok := (op.(predicate.BoolOp))
	if !ok {
		panic("unable to convert predicate")
	}
	s.predicate = p
	return
}

func (s *Selector) true(e *logrus.Entry) bool {
	defer s.m.RUnlock()
	s.m.RLock()
	return s.predicate.True(e)
}

func (s *Selector) MaybeRead() (e *logrus.Entry) {
	if !s.reg {
		return nil
	}
	found := false
	for !found {
		select {
		case ent := <-s.q:
			if s.true(&ent) {
				found = true
			}
			e = &ent
			break
		default:
			e = nil
			return
			break
		}
	}
	return
}

func (s *Selector) Stop() {
	s.d.Unregister(s)
	s.reg = false
}

func (s *Selector) BaseTime() time.Time {
	return baseTimestamp
}