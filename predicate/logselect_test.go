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
 * logselect_test.go: Test of the select grammar
 */

package predicate

import (
	"testing"
	"fmt"
	"reflect"
	"github.com/sirupsen/logrus"
)

type tst struct {
	input string
	success bool
	output interface{}
}

func doTest(s tst, opts ...Option) (func(t *testing.T)) {
	return func(t *testing.T) {
		//t.Parallel()
		fmt.Printf("Testing \"%s\"\n", s.input)
		iface, err := Parse("test", []byte(s.input), opts...)
		if err != nil && s.success {
			fmt.Printf("Error: %s\n", err.Error())
			t.Fail()
			return
		} else if err == nil && !s.success {
			fmt.Printf("Expecting an error and got none\n")
			t.Fail()
			return
		} else if err != nil && !s.success {
			fmt.Printf("Successfully got an error\n")
			return
		}
		typ := "nil"
		if iface != nil {
			typ = reflect.TypeOf(iface).String()
		}
		if reflect.TypeOf(iface) != reflect.TypeOf(s.output) {
			fmt.Printf("Type mismatch, expected %s but got %s\n", reflect.TypeOf(s.output), reflect.TypeOf(iface))
			t.Fail()
			return
		}
		if iface != s.output {
			fmt.Printf("Value mismatch expected %v but got %v\n", s.output, iface)
			t.Fail()
			return
		}
		fmt.Printf("Got %#v -> %s\n", iface, typ)
	}
}

func TestParse_Number(t *testing.T) {

	var tests = []tst{
		{"0", true, int64(0)},
		{"1", true, int64(1)},
		{"-1", true, int64(-1)},
		{"0.", true, float64(0)},
		{"1.", true, float64(1)},
		{"-1.", true, float64(-1)},
		{"0.0", false, nil},
		{"-0", false, nil},
		{"-0.0", false, nil},
		{"0.1", true, float64(0.1)},
		{"1.1", true, float64(1.1)},
		{"-1.1", true, float64(-1.1)},
		{"11", true, int64(11)},
		{"-11", true, int64(-11)},
		{"1234567890.0", true, float64(1234567890.0)},
		{"1234567890", true, int64(1234567890)},
	}

	for _, s := range tests {
		t.Run(s.input, doTest(s, Entrypoint("Number")))
	}
}

func TestParse_String(t *testing.T) {
	var tests = []tst{
		{"''", true, ""},
		{"\"\"", true, ""},
		{"'hello'", true, "hello"},
		{"\"world\"", true, "world"},
		{"'π'", true, "π"},
		{"'it\\'s π day'", true, "it's π day" },
		{"\"it's π day\"", true, "it's π day"},
		{"'\"we\\'ll see about that \"'", true, "\"we'll see about that \""},
		{"\"\\\"we'll see about that \\\"\"", true, "\"we'll see about that \""},
	}

	for _, s := range tests {
		t.Run(s.input, doTest(s, Entrypoint("String")))
	}
}

func TestParse_Ident(t *testing.T) {
	var tests = []tst{
		{"hello", true, "hello"},
		{"hello-world", true, "hello-world"},
		{"h0la", true, "h0la"},
	}

	for _, s := range tests {
		t.Run(s.input, doTest(s, Entrypoint("Ident")))
	}
}

func TestParse_OpBool(t *testing.T) {
	var tests = []tst{
		{"Prefix(hello-world)", true, OpPrefix{"hello-world"}},
		{"HasField(hi-hi)", true, OpHasField{"hi-hi"}},
		{"Prefix('π day')", true, OpPrefix{"π day"}},
	}

	for _, s := range tests {
		t.Run(s.input, doTest(s))
	}
}

func TestParse_BoolAndOr(t *testing.T) {
	var tests = []tst {
		{"Prefix(hello) && HasField('world')", true, OpAnd{OpPrefix{"hello"}, OpHasField{"world"}}},
		{"Prefix(hello) && HasField('world') && HasField(\"worker\")", true, OpAnd{OpPrefix{"hello"}, OpAnd{OpHasField{"world"}, OpHasField{"worker"}}}},
		{"Prefix(hello) || HasField('world')", true, OpOr{OpPrefix{"hello"}, OpHasField{"world"}}},
		{"Prefix(hello) && HasField('world') || HasField(\"worker\")", true, OpOr{OpAnd{OpPrefix{"hello"}, OpHasField{"world"}}, OpHasField{"worker"}}},
		{"Prefix('1') && HasField('2') && HasField('3') && Prefix('4')", true, OpAnd{OpPrefix{"1"}, OpAnd{OpHasField{"2"}, OpAnd{OpHasField{"3"}, OpPrefix{"4"}}}}},
	}

	for _, s := range tests {
		t.Run(s.input, doTest(s))
	}
}

func TestParse_BoolNot(t *testing.T) {
	var tests = []tst {
		{"!Prefix(hello)", true, OpNot{OpPrefix{"hello"}}},
		{"!Prefix(hello) && HasField('world') && HasField(\"worker\")", true, OpAnd{OpNot{OpPrefix{"hello"}}, OpAnd{OpHasField{"world"}, OpHasField{"worker"}}}},
		{"!(Prefix(hello) || HasField('world'))", true, OpNot{OpOr{OpPrefix{"hello"}, OpHasField{"world"}}}},
		{"Prefix(hello) && HasField('world') || !HasField(\"worker\")", true, OpOr{OpAnd{OpPrefix{"hello"}, OpHasField{"world"}}, OpNot{OpHasField{"worker"}}}},
		{"Prefix('1') && HasField('2') && !(HasField('3') && Prefix('4'))", true, OpAnd{OpPrefix{"1"}, OpAnd{OpHasField{"2"}, OpNot{OpAnd{OpHasField{"3"}, OpPrefix{"4"}}}}}},
	}

	for _, s := range tests {
		t.Run(s.input, doTest(s))
	}
}

func TestParse_Value(t *testing.T) {
	var tests = []tst {
		{"1", true, Val{typ:ValTypeInt, itg:1}},
		{"1.0", true, Val{typ:ValTypeFloat, flt:1.0}},
		{"0", true, Val{typ:ValTypeInt, itg:0}},
		{"0.0", false, nil},
		{"'hello'", true, Val{typ:ValTypeString, str:"hello"}},
		{"hello", true, Val{typ: ValTypeString, str: "hello"}},
		{"true", true, Val{typ: ValTypeBool, bl: true}},
		{"false",true, Val{typ: ValTypeBool, bl: false}},
		{"nil", true, Val{typ: ValTypeNil}},
		{"null", true, Val{typ: ValTypeNil}},
		{"panic", true, LogLevel{int64(logrus.PanicLevel)}},
		{"WARN", true, LogLevel{int64(logrus.WarnLevel)}},
	}

	for _, s := range tests {
		t.Run(s.input, doTest(s, Entrypoint("Value")))
	}
}

func TestParse_EmptyString(t *testing.T) {
	t.Run("<empty string>", doTest(tst{"", true, OpTrue{}}))
}

func TestParse_Comparison(t *testing.T) {
	var tests = []tst {
		{"1 == 1", true, OpEquals{Val{typ:ValTypeInt, itg:1}, Val{typ: ValTypeInt, itg: 1}}},
		{"Field('hello') == 'world'", true, OpEquals{OpField{"hello"}, Val{typ: ValTypeString, str: "world"}}},
		{"'world' != Field('hello')", true, OpNot{OpEquals{Val{typ: ValTypeString, str: "world"}, OpField{"hello"}}}},
		{"HasField(hello) && Field('hello') == 'world'", true, OpAnd{OpHasField{"hello"}, OpEquals{OpField{"hello"}, Val{typ: ValTypeString, str: "world"}}}},
		{"(HasField(hello)) && Field('hello') == 'world'", true, OpAnd{OpHasField{"hello"}, OpEquals{OpField{"hello"}, Val{typ: ValTypeString, str: "world"}}}},
		{"HasField(hello) && (Field('hello') == 'world')", true, OpAnd{OpHasField{"hello"}, OpEquals{OpField{"hello"}, Val{typ: ValTypeString, str: "world"}}}},
		{"HasField(hello) && !(Field('hello') == 'world')", true, OpAnd{OpHasField{"hello"}, OpNot{OpEquals{OpField{"hello"}, Val{typ: ValTypeString, str: "world"}}}}},
		{"Field('hello') > 'world'", true, OpGreater{OpField{"hello"}, Val{typ: ValTypeString, str: "world"}}},
		{"Field('hello') < 'world'", true, OpLess{OpField{"hello"}, Val{typ: ValTypeString, str: "world"}}},
		{"Field('hello') >= 'world'", true, OpOr{OpEquals{OpField{"hello"}, Val{typ: ValTypeString, str: "world"}},OpGreater{OpField{"hello"}, Val{typ: ValTypeString, str: "world"}}}},
		{"Field('hello') <= 'world'", true, OpOr{OpEquals{OpField{"hello"}, Val{typ: ValTypeString, str: "world"}}, OpLess{OpField{"hello"}, Val{typ: ValTypeString, str: "world"}}}},
		{"Field('hello')== 'world'", true, OpEquals{OpField{"hello"}, Val{typ: ValTypeString, str: "world"}}},

	}

	for _, s := range tests {
		t.Run(s.input, doTest(s))
	}
}