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

import "net/http"

// Supported error types
var (
	Default    = Type{meta: "default"}
	SubtaskErr = Type{meta: "subtask"}
	NotFound   = Type{httpCode: http.StatusNotFound, meta: "not-found"}
	Internal   = Type{httpCode: http.StatusInternalServerError, meta: "internal"}
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
	}
)

// New constructs a new Error instance with this message
func (t *Type) New(message string, opts ...Option) Error {
	return newCrdbError(t, nil, message, opts...)
}

// Wrap constructs a new Error instance with this message and wraps the passed in error
func (t *Type) Wrap(err error, message string, opts ...Option) Error {
	return newCrdbError(t, err, message, opts...)
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
