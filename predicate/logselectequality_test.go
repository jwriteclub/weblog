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
 * logselectequality_test.go: Bool op equality test
 */

package predicate

import (
	"testing"
	"fmt"
)

func TestVal_EqualsSimple(t *testing.T) {
	var tests = []tst {
		{"", true, nil},
		{"1 == 1", true, nil},
		{"1 == 0", false, nil},
	}

	equals := func(test tst) {
		t.Run(test.input, func(t *testing.T) {
			op, err := Parse("test", []byte(test.input))
			if err != nil {
				fmt.Printf("Unable to parse %s: %s\n", test.input, err.Error())
				t.Fail()
				return
			}
			bo, ok := op.(BoolOp)
			if !ok {
				fmt.Printf("Got back invalid op\n")
				t.Fail()
				return
			}
			fmt.Printf("%#v\n", bo)
			val := bo.True(nil)
			if val != test.success {
				fmt.Printf("Got mismatch between expected %v and actual %v\n", test.success, val)
				t.Fail()
				return
			}
		})
	}

	for _, test := range tests {
		equals(test)
	}

}