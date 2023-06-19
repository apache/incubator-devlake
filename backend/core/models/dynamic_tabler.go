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

package models

import (
	"encoding/json"
	"reflect"

	"github.com/apache/incubator-devlake/core/errors"

	"github.com/apache/incubator-devlake/core/dal"
)

// DynamicTabler is a core.Tabler that wraps a runtime (anonymously) generated data-model. Due to limitations of
// reflection in Go and the GORM framework, the underlying model and the table have to be explicitly passed into dal.Dal's API
// via Unwrap() and TableName()
type DynamicTabler interface {
	dal.Tabler
	NewValue() any
	New() DynamicTabler
	NewSlice() DynamicTabler
	From(src any) errors.Error
	To(target any) errors.Error
	Unwrap() any
	UnwrapPtr() *any
	UnwrapSlice() []any
}

// DynamicTablerImpl the implementation of DynamicTabler
type DynamicTablerImpl struct {
	objType reflect.Type
	wrapped any
	table   string
}

func NewDynamicTabler(tableName string, objType reflect.Type) DynamicTabler {
	return &DynamicTablerImpl{
		objType: objType,
		table:   tableName,
	}
}

func (d *DynamicTablerImpl) NewValue() any {
	return reflect.New(d.objType).Interface()
}

func (d *DynamicTablerImpl) New() DynamicTabler {
	return &DynamicTablerImpl{
		objType: d.objType,
		wrapped: d.NewValue(),
		table:   d.table,
	}
}

func (d *DynamicTablerImpl) NewSlice() DynamicTabler {
	sliceType := reflect.SliceOf(d.objType)
	return &DynamicTablerImpl{
		objType: sliceType,
		wrapped: reflect.New(sliceType).Interface(),
		table:   d.table,
	}
}

func (d *DynamicTablerImpl) From(src any) errors.Error {
	b, err := json.Marshal(src)
	if err != nil {
		return errors.Convert(err)
	}
	return errors.Convert(json.Unmarshal(b, d.wrapped))
}

func (d *DynamicTablerImpl) To(target any) errors.Error {
	b, err := json.Marshal(d.wrapped)
	if err != nil {
		return errors.Convert(err)
	}
	return errors.Convert(json.Unmarshal(b, target))
}

func (d *DynamicTablerImpl) Unwrap() any {
	return d.wrapped
}

func (d *DynamicTablerImpl) UnwrapPtr() *any {
	return &d.wrapped
}

func (d *DynamicTablerImpl) UnwrapSlice() []any {
	var arr []any
	slice := reflect.ValueOf(d.wrapped).Elem()
	for i := 0; i < slice.Len(); i++ {
		arr = append(arr, slice.Index(i).Interface())
	}
	return arr
}

func (d *DynamicTablerImpl) TableName() string {
	return d.table
}

var _ DynamicTabler = (*DynamicTablerImpl)(nil)

// UnwrapObject if the actual object is wrapped in some proxy, it unwinds and returns it, otherwise this is idempotent
func UnwrapObject(ifc any) any {
	if dynamic, ok := ifc.(DynamicTabler); ok {
		return dynamic.Unwrap()
	}
	return ifc
}
