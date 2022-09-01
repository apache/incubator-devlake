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
	"net/http"
	"strings"
)

// Supported error types
var (
	Default      = register(nil)
	SubtaskErr   = register(&Type{meta: "subtask"})
	NotFound     = register(&Type{httpCode: http.StatusNotFound, meta: "not-found"})
	BadInput     = register(&Type{httpCode: http.StatusBadRequest, meta: "bad-input"})
	Unauthorized = register(&Type{httpCode: http.StatusUnauthorized, meta: "unauthorized"})
	Forbidden    = register(&Type{httpCode: http.StatusForbidden, meta: "forbidden"})
	Internal     = register(&Type{httpCode: http.StatusInternalServerError, meta: "internal"})
	Timeout      = register(&Type{httpCode: http.StatusGatewayTimeout, meta: "timeout"})

	//cached values
	typesByHttpCode = map[int]*Type{}
)

type (
	// Type error are constructed from these, and they contain metadata about them.
	Type struct {
		meta string
		// below are optional fields
		httpCode int
	}

	// Option add customized properties to the Error
	Option func(*options)

	options struct {
		userMsg   string
		asUserMsg bool
		data      interface{}
	}
)

func HttpStatus(code int) *Type {
	t, ok := typesByHttpCode[code]
	if !ok {
		t = Internal
	}
	return t
}

// New constructs a new Error instance with this message
func (t *Type) New(message string, opts ...Option) Error {
	return newCrdbError(t, nil, message, opts...)
}

// Wrap constructs a new Error instance with this message and wraps the passed in error
func (t *Type) Wrap(err error, message string, opts ...Option) Error {
	return newCrdbError(t, err, message, opts...)
}

// WrapRaw constructs a new Error instance that directly wraps this error with no additional context
func (t *Type) WrapRaw(err error) Error {
	msg := ""
	lakeErr := AsLakeErrorType(err)
	if lakeErr != nil {
		msg = "" // there's nothing new to add
	} else {
		msg = err.Error()
	}
	return newCrdbError(t, err, msg)
}

// Combine constructs a new Error from combining multiple errors. Stacktrace info for each of the errors will not be present in the result.
func (t *Type) Combine(errs []error, msg string, opts ...Option) Error {
	msgs := []string{}
	for _, e := range errs {
		if le := AsLakeErrorType(e); le != nil {
			if msg0 := le.Message(); msg0 != "" {
				msgs = append(msgs, le.Message())
			}
		} else {
			msgs = append(msgs, e.Error())
		}
	}
	effectiveMsg := strings.Join(msgs, "\n=====================\n")
	effectiveMsg = "\t" + strings.ReplaceAll(effectiveMsg, "\n", "\n\t")
	return newCrdbError(t, nil, fmt.Sprintf("%s\ncombined messages: \n{\n%s\n}", msg, effectiveMsg), opts...)
}

// GetHttpCode gets the associated Http code with this Type, if explicitly set, otherwise http.StatusInternalServerError
func (t *Type) GetHttpCode() int {
	if t.httpCode == 0 {
		return http.StatusInternalServerError
	}
	return t.httpCode
}

// UserMessage add a user-friendly message to the Error
func UserMessage(msg string) Option {
	return func(opts *options) {
		opts.userMsg = msg
	}
}

// AsUserMessage use the ordinary message as the user-friendly message of the Error
func AsUserMessage() Option {
	return func(opts *options) {
		opts.asUserMsg = true
	}
}

// WithData associate data with this Error
func WithData(data interface{}) Option {
	return func(opts *options) {
		opts.data = data
	}
}

func register(t *Type) *Type {
	if t == nil {
		t = &Type{meta: "default"}
		typesByHttpCode[t.httpCode] = t
	} else if t.httpCode != 0 {
		typesByHttpCode[t.httpCode] = t
	}
	return t
}
