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
 * types_test.go: Static tests of the type system
 */

package predicate

import (
	"testing"
	"reflect"
	"fmt"
	"github.com/sirupsen/logrus"
)

func TestBoolOp(t *testing.T) {
	assertBoolOp := func(i interface{}, name string) {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, ok := i.(BoolOp)
			if !ok {
				fmt.Printf("%s [%s] is not a BoolOp", name, reflect.TypeOf(i))
				t.Fail()
			}
		})
	}

	assertBoolOp(OpAnd{}, "OpAnd")
	assertBoolOp(OpOr{}, "OpOr")
	assertBoolOp(OpPrefix{}, "OpPrefix")
	assertBoolOp(OpHasField{}, "OpHasField")
	assertBoolOp(OpNot{}, "OpNot")
	assertBoolOp(OpEquals{}, "OpEquals")
	assertBoolOp(OpGreater{}, "OpGreater")
	assertBoolOp(OpLess{}, "OpLess")
	assertBoolOp(OpTrue{}, "OpTrue")
	assertBoolOp(OpFalse{}, "OpFalse")
}

func TestValuable(t *testing.T) {
	assertValuable := func(i interface{}, name string) {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			_, ok := i.(Valueable)
			if !ok {
				fmt.Printf("%s [%s] is not a Valuable", name, reflect.TypeOf(i))
				t.Fail()
			}
		})
	}

	assertValuable(Val{}, "Val")
	assertValuable(LogLevel{}, "LogLevel")
	assertValuable(OpField{}, "OpField")
}

func TestOpAnd_True(t *testing.T) {
	a := OpAnd{OpTrue{}, OpTrue{}}
	if !a.True(nil) {
		t.Fail()
	}

	b := OpAnd{OpTrue{}, OpFalse{}}
	if b.True(nil) {
		t.Fail()
	}

	c := OpAnd{OpFalse{}, OpTrue{}}
	if c.True(nil) {
		t.Fail()
	}

	d := OpAnd{OpFalse{}, OpFalse{}}
	if d.True(nil) {
		t.Fail()
	}
}

func TestOpOr_True(t *testing.T) {
	a := OpOr{OpTrue{}, OpTrue{}}
	if !a.True(nil) {
		t.Fail()
	}

	b := OpOr{OpTrue{}, OpFalse{}}
	if !b.True(nil) {
		t.Fail()
	}

	c := OpOr{OpFalse{}, OpTrue{}}
	if !c.True(nil) {
		t.Fail()
	}

	d := OpOr{OpFalse{}, OpFalse{}}
	if d.True(nil) {
		t.Fail()
	}
}

func TestOpNot_True(t *testing.T) {
	a := OpNot{OpTrue{}}
	if a.True(nil) {
		t.Fail()
	}

	b := OpNot{OpFalse{}}
	if !b.True(nil) {
		t.Fail()
	}
}

func TestOpPrefix_True(t *testing.T) {
	e := &logrus.Entry{Data: make(logrus.Fields)}
	e = e.WithField("prefix", "hello")

	a := OpPrefix{"hello"}
	if !a.True(e) {
		t.Fail()
	}

	b := OpPrefix{"world"}
	if b.True(e) {
		t.Fail()
	}
}

func TestOpHasField_True(t *testing.T) {
	e := &logrus.Entry{Data: make(logrus.Fields)}
	e = e.WithField("hello", "hello")

	a := OpHasField{"hello"}
	if !a.True(e) {
		t.Fail()
	}

	b := OpHasField{"world"}
	if b.True(e) {
		t.Fail()
	}
}

func TestOpField_toVal(t *testing.T) {
	tests := []tst{
		{"1.0", true, Val{typ: ValTypeFloat, flt: 1.0}},
		{"1", true, Val{typ: ValTypeInt, itg: 1}},
		{"-1", true, Val{typ: ValTypeInt, itg: -1}},
		{"true", true, Val{typ:ValTypeBool, bl: true}},
		{"false", true, Val{typ:ValTypeBool, bl: false}},
		{"nil", true, Val{typ:ValTypeNil}},
		{"null", true, Val{typ:ValTypeNil}},
		{"True", true, Val{typ:ValTypeBool, bl: true}},
		{"False", true, Val{typ:ValTypeBool, bl: false}},
		{"NIL", true, Val{typ:ValTypeNil}},
		{"NuLl", true, Val{typ:ValTypeNil}},
		{"world", true, Val{typ:ValTypeString, str: "world"}},
	}


	a := OpField{"a"}
	tv := func (test tst) func(*testing.T) {
		return func(t *testing.T) {
			e := (&logrus.Entry{Data: make(logrus.Fields)}).WithField("a", test.input)
			if a.toVal(e) != test.output {
				t.Fail()
			}
		}
	}

	for _, test := range tests {
		t.Run(test.input, tv(test))
	}
}

// Private types for testing
