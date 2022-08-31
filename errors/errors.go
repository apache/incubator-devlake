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

type (
	requiredSupertype interface {
		error
		Unwrap() error
	}
	// Error The interface that all internally managed errors should adhere to.
	Error interface {
		requiredSupertype
		// Message the message associated with this Error.
		Message() string
		// UserMessage the message associated with this Error appropriated for end users.
		UserMessage() string
		// GetType gets the Type of this error
		GetType() *Type
		// As Attempts to cast this Error to the requested Type, and returns nil if it can't.
		As(*Type) Error
		// GetData returns the data associated with this Error (may be nil)
		GetData() any
	}
)

// AsLakeErrorType attempts to cast err to Error, otherwise returns nil
func AsLakeErrorType(err error) Error {
	if cast, ok := err.(Error); ok {
		return cast
	}
	return nil
}

var _ error = (Error)(nil)
var _ requiredSupertype = (Error)(nil)
