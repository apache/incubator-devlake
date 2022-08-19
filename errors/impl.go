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
	"strings"
)
import cerror "github.com/cockroachdb/errors"

type (
	crdbErrorImpl struct {
		wrappedRaw error
		wrapped    *crdbErrorImpl
		userMsg    string
		msg        string
		t          *Type
	}
)

func (e *crdbErrorImpl) Error() string {
	return fmt.Sprintf("%+v", e.wrappedRaw)
}

func (e *crdbErrorImpl) Message() string {
	return e.wrappedRaw.Error()
}

func (e *crdbErrorImpl) UserMessage() string {
	return strings.Join(e.getUserMessages(), "\ncaused by: ")
}

func (e *crdbErrorImpl) Unwrap() error {
	if e.wrapped != nil {
		return e.wrapped
	}
	return cerror.Cause(e.wrappedRaw)
}

func (e *crdbErrorImpl) GetType() Type {
	return *e.t
}

func (e *crdbErrorImpl) As(t Type) Error {
	err := e
	for {
		if *err.t == t {
			return e
		}
		lakeErr := AsLakeErrorType(err.Unwrap())
		if lakeErr == nil {
			return nil
		}
		err = lakeErr.(*crdbErrorImpl)
	}
}

func (e *crdbErrorImpl) getUserMessages() []string {
	msgs := []string{}
	err := e
	ok := false
	for {
		if err.userMsg != "" {
			msgs = append(msgs, err.userMsg)
		}
		unwrapped := err.Unwrap()
		if unwrapped == nil {
			break
		}
		err, ok = unwrapped.(*crdbErrorImpl)
		if !ok {
			// don't append the message if the error is "external"
			break
		}
	}
	return msgs
}

func newCrdbError(t *Type, err error, message string, opts ...Option) *crdbErrorImpl {
	cfg := &options{}
	for _, opt := range opts {
		opt(cfg)
	}
	errType := *t
	var wrappedErr *crdbErrorImpl
	if cast, ok := err.(*crdbErrorImpl); ok {
		err = cast.wrappedRaw
		wrappedErr = cast
		if *t == Default { // inherit wrapped error's type
			errType = cast.GetType()
		}
	}
	impl := &crdbErrorImpl{
		wrappedRaw: cerror.WrapWithDepth(1, err, message),
		wrapped:    wrappedErr,
		msg:        message,
		userMsg:    cfg.userMsg,
		t:          &errType,
	}
	if cfg.asUserMsg {
		impl.userMsg = impl.msg
	}
	return impl
}

var _ Error = (*crdbErrorImpl)(nil)
