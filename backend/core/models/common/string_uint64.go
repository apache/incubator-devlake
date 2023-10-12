/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"reflect"
)

type StringUint64 struct {
	v uint64
	t string
}

func NewStringUint64(value uint64) *StringUint64 {
	return &StringUint64{
		v: value,
		t: "uint64",
	}
}

func NewStringUint64FromAny(value interface{}) *StringUint64 {
	return &StringUint64{
		v: cast.ToUint64(value),
		t: "string",
	}
}

func (f *StringUint64) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

func (f *StringUint64) Uint64() uint64 {
	if f == nil {
		return 0
	}
	return f.v
}

func (f *StringUint64) String() string {
	//if f.t == "string" {
	//	return fmt.Sprintf("\"%v\"", f.v)
	//}
	return fmt.Sprintf("%v", f.v)
}

func (f *StringUint64) UnmarshalJSON(data []byte) error {
	var i interface{}
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	switch i.(type) {
	case uint64:
		f.t = "uint64"
	case string:
		f.t = "string"
	}
	value, err := cast.ToUint64E(i)
	if err != nil {
		return err
	}
	f.v = value
	return nil
}

func (f *StringUint64) Value() (driver.Value, error) {
	if f == nil {
		return nil, nil
	}
	return f.v, nil
}

func (f *StringUint64) Scan(v interface{}) error {
	switch value := v.(type) {
	case uint8, uint16, uint32, uint64, int8, int16, int32, int64, int, uint:
		*f = StringUint64{
			v: cast.ToUint64(value),
			t: "uint64",
		}
	case string:
		*f = StringUint64{
			v: cast.ToUint64(value),
			t: "string",
		}
	default:
		return fmt.Errorf("[StringUint64] %+v is an unknown type, type: %v with value: %v", v, reflect.TypeOf(v), value)
	}
	return nil
}
