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

package api

import (
	"github.com/apache/incubator-devlake/core/models"
	"reflect"
)

func reflectField(obj any, fieldName string) reflect.Value {
	obj = models.UnwrapObject(obj)
	return reflectValue(obj).FieldByName(fieldName)
}

func hasField(obj any, fieldName string) bool {
	obj = models.UnwrapObject(obj)
	_, ok := reflectType(obj).FieldByName(fieldName)
	return ok
}

func reflectValue(obj any) reflect.Value {
	obj = models.UnwrapObject(obj)
	val := reflect.ValueOf(obj)
	kind := val.Kind()
	for kind == reflect.Ptr || kind == reflect.Interface {
		val = val.Elem()
		kind = val.Kind()
	}
	return val
}

func reflectType(obj any) reflect.Type {
	obj = models.UnwrapObject(obj)
	typ := reflect.TypeOf(obj)
	kind := typ.Kind()
	for kind == reflect.Ptr {
		typ = typ.Elem()
		kind = typ.Kind()
	}
	return typ
}
