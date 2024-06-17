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

package utils

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/apache/incubator-devlake/core/errors"
)

type JsonObject = map[string]any
type JsonArray = []any

func GetProperty[T any](object JsonObject, key string) (T, errors.Error) {
	property, ok := object[key]
	if !ok {
		return *new(T), errors.Default.New(fmt.Sprintf("Missing property \"%s\"", key))
	}
	return Convert[T](property)
}

func GetItem[T any](array JsonArray, index int) (T, errors.Error) {
	if index < 0 || index >= len(array) {
		return *new(T), errors.Default.New(fmt.Sprintf("Index %d out of range", index))
	}
	return Convert[T](array[index])
}

// Convert converts value to type T. If value is a slice, it converts each element of the slice to type T.
// Does not support nested slices.
func Convert[T any](value any) (T, errors.Error) {
	var t T
	tType := reflect.TypeOf(t)
	if tType.Kind() == reflect.Slice {
		valueSlice, ok := value.([]any)
		if !ok {
			return t, errors.Default.New("Value is not a slice")
		}
		elemType := tType.Elem()
		result := reflect.MakeSlice(tType, 0, len(valueSlice))
		for i, v := range valueSlice {
			value := reflect.ValueOf(v)
			if elemType.AssignableTo(reflect.TypeOf(v)) {
				elem := value.Convert(elemType)
				result = reflect.Append(result, elem)
			} else {
				return t, errors.Default.New(fmt.Sprintf("Element %d is not of type %s", i, elemType.Name()))
			}
		}
		return result.Interface().(T), nil
	} else {
		result, ok := value.(T)
		if !ok {
			return t, errors.Default.New(fmt.Sprintf("Value is not of type %T", t))
		}
		return result, nil
	}
}

func ToJsonString(x any) string {
	b, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}
	return string(b)
}
