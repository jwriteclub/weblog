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
 * types.go: Predicate types
 */

package predicate

import (
	"github.com/sirupsen/logrus"
	"strings"
	"strconv"
	"fmt"
)

type BoolOp interface {
	True(e *logrus.Entry) bool
}

type OpPrefix struct {
	prefix string
}
func (p OpPrefix) True(e *logrus.Entry) bool {
	pfx, ok := e.Data["prefix"]
	if !ok {
		return false
	}
	return pfx == p.prefix
}

type OpHasField struct {
	field string
}
func (h OpHasField) True(e *logrus.Entry) bool {
	_, ok := e.Data[h.field]
	return ok
}

type OpAnd struct {
	left, right BoolOp
}
func (a OpAnd) True(e *logrus.Entry) bool {
	l := a.left.True(e)
	if !l {
		return false
	}
	return a.right.True(e)
}

type OpOr struct {
	left, right BoolOp
}
func (o OpOr) True(e *logrus.Entry) bool {
	l := o.left.True(e)
	if l {
		return true
	}
	return o.right.True(e)
}

type OpNot struct {
	inner BoolOp
}
func (n OpNot) True(e *logrus.Entry) bool {
	return !n.inner.True(e)
}

type OpEquals struct {
	left, right Valueable
}
func (o OpEquals) True(e *logrus.Entry) bool {
	return o.left.Equals(o.right, e)
}

type OpGreater struct {
	left, right Valueable
}
func (g OpGreater) True(e *logrus.Entry) bool {
	// TODO: Complete this
	return false
}

type OpLess struct {
	left, right Valueable
}
func (l OpLess) True(e *logrus.Entry) bool {
	// TODO: Complete this
	return false
}

type OpTrue struct {}
func (o OpTrue) True(e *logrus.Entry) bool {
	return true
}
type OpFalse struct{}
func (o  OpFalse) True(e *logrus.Entry) bool {
	return false
}

type ValType int
const (
	ValTypeString ValType = iota
	ValTypeFloat
	ValTypeInt
	ValTypeBool
	ValTypeNil
)

type Valueable interface {
	Type(e *logrus.Entry) ValType
	Equals(v Valueable, e *logrus.Entry) bool
	GetVal(e *logrus.Entry) interface{}
}

type Val struct {
	typ ValType
	str string
	flt float64
	itg int64
	bl bool
}
func (v Val) Type(e *logrus.Entry) ValType {
	return v.typ
}
func (v Val) Equals(o Valueable, e *logrus.Entry) bool {
	if v.typ != o.Type(e) {
		return false
	}
	switch v.typ {
	case ValTypeString:
		return v.str == o.GetVal(e).(string)
	case ValTypeFloat:
		return v.flt == o.GetVal(e).(float64)
	case ValTypeInt:
		return v.itg == o.GetVal(e).(int64)
	default:
		panic("impossible value type")
	}
}
func (v Val) GetVal(e *logrus.Entry) interface{} {
	switch v.typ {
	case ValTypeString:
		return v.str
	case ValTypeFloat:
		return v.flt
	case ValTypeInt:
		return v.itg
	default:
		panic("impossible value type")
	}
}

type LogLevel struct {
	v int64
}
func newLogLevel(name string) (interface{}, error) {
	lvl, err := logrus.ParseLevel(strings.ToLower(name))
	if err != nil {
		return nil, err
	}
	return LogLevel{int64(lvl)}, nil
}
func (l LogLevel) Type(e *logrus.Entry) ValType {
	return ValTypeInt
}
func (l LogLevel) Equals(o Valueable, e *logrus.Entry) bool {
	if o.Type(e) != ValTypeInt {
		return false
	}
	return l.v == o.GetVal(e).(int64)
}
func (l LogLevel) GetVal(e *logrus.Entry) interface{} {
	return l.v
}

// TODO: Finish these for real
type OpField struct {
	name string
}
func (f OpField) Type(e *logrus.Entry) ValType {
	return ValTypeString
}
func (f OpField) Equals(o Valueable, e *logrus.Entry) bool {
	return false
}
func (f OpField) GetVal(e *logrus.Entry) interface{} {
	return nil
}
func (f OpField) toVal(e *logrus.Entry) Val {
	var ok bool
	var s interface{}
	var str string
	s, ok = e.Data[f.name]
	if !ok {
		return Val{typ: ValTypeNil}
	}
	str, ok = s.(string)
	if !ok {
		var stg fmt.Stringer
		stg, ok = s.(fmt.Stringer)
		if !ok {
			return Val{typ: ValTypeNil}
		}
		str = stg.String()
	}
	itg, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return Val{typ: ValTypeInt, itg: itg}
	}
	flt, err := strconv.ParseFloat(str, 64)
	if err == nil {
		return Val{typ: ValTypeFloat, flt: flt}
	}
	sl := strings.ToLower(str)
	if sl == "true" {
		return Val{typ: ValTypeBool, bl: true}
	}
	if sl == "false" {
		return Val{typ: ValTypeBool, bl: false}
	}
	if sl == "null" || sl == "nil" {
		return Val{typ: ValTypeNil}
	}
	return Val{typ: ValTypeString, str: str}
}