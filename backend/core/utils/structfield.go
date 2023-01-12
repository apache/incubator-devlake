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
	"reflect"
)

// WalkFields get the field data by tag
func WalkFields(t reflect.Type, filter func(field *reflect.StructField) bool) (f []reflect.StructField) {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if filter == nil {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Type.Kind() == reflect.Struct {
				f = append(f, WalkFields(field.Type, filter)...)
			} else {
				f = append(f, field)
			}
		}
	} else {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if filter(&field) {
				f = append(f, field)
			} else if field.Type.Kind() == reflect.Struct {
				f = append(f, WalkFields(field.Type, filter)...)
			}
		}
	}

	return f
}
