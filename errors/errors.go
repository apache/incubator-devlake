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

package errors

import (
	"fmt"
	"github.com/joomcode/errorx"
	"net/http"
)

// Supported error types
var (
	Default  = Type{code: 0, meta: "default"}
	NotFound = Type{code: http.StatusNotFound, meta: "not-found"}
	Internal = Type{code: http.StatusInternalServerError, meta: "internal"}
)

var errorxNamespace = errorx.NewNamespace("lake")
var errorxTypes = make(map[Type]*errorx.Type)

var _ error = (*Error)(nil)
var _ requiredSupertype = (*Error)(nil)

type (
	Type struct {
		code int
		meta string
	}
	Error struct {
		err *errorx.Error
		t   *Type
	}
	requiredSupertype interface {
		Unwrap() error
	}
)

func (e *Error) Error() string {
	return fmt.Sprintf("%+v", e.err)
}

func (e *Error) Message() string {
	return e.err.Message()
}

func (e *Error) Messages() string {
	return e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err.Unwrap()
}

func (e *Error) GetType() Type {
	return *e.t
}

func (e *Error) As(t Type) *Error {
	err := e
	for {
		if *err.t == t {
			return e
		}
		err = AsLakeErrorType(err.Unwrap())
		if err == nil {
			return nil
		}
	}
}

func (t *Type) New(message string) *Error {
	return &Error{
		err: t.getErrorxType().New(message),
		t:   t,
	}
}

func (t *Type) Wrap(err error, message string) *Error {
	errType := *t
	if cast, ok := err.(*Error); ok {
		err = cast.err
		if *t == Default { // inherit wrapped error's type
			errType = cast.GetType()
		}
	}
	return &Error{
		err: t.getErrorxType().Wrap(err, message),
		t:   &errType,
	}
}

func (t *Type) getErrorxType() *errorx.Type {
	val, ok := errorxTypes[*t]
	if !ok {
		val = errorxNamespace.NewType(t.meta)
		errorxTypes[*t] = val
	}
	return val
}

func AsLakeErrorType(err error) *Error {
	if cast, ok := err.(*Error); ok {
		return cast
	}
	return nil
}
