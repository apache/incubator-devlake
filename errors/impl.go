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
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)
import cerror "github.com/cockroachdb/errors"

type (
	crdbErrorImpl struct {
		wrappedRaw error
		wrapped    *crdbErrorImpl
		msg        string
		data       interface{}
		t          *Type
	}
)

var enableStacktraces = false

func init() {
	enable, exists := os.LookupEnv("ENABLE_STACKTRACE")
	if !exists {
		return
	}
	enableStacktraces, _ = strconv.ParseBool(enable)
}

func (e *crdbErrorImpl) Error() string {
	//crdb spits out a bunch of excess strings, so do some cleanup
	rawMsg := fmt.Sprintf("%+v", e.wrappedRaw)
	parts := strings.Split(rawMsg, "\n(1) ")
	if len(parts) == 1 {
		return parts[0]
	}
	return parts[1]
}

func (e *crdbErrorImpl) Message() string {
	return strings.Join(e.getMessages(func(err *crdbErrorImpl) string {
		if err.msg == "" {
			return ""
		}
		code := ""
		if err.t.httpCode != 0 {
			code = fmt.Sprintf("(%d)", err.t.httpCode)
		}
		return err.msg + " " + code
	}), "\ncaused by: ")
}

func (e *crdbErrorImpl) Unwrap() error {
	if e.wrapped != nil {
		return e.wrapped
	}
	return cerror.Cause(e.wrappedRaw)
}

func (e *crdbErrorImpl) GetType() *Type {
	return e.t
}

func (e *crdbErrorImpl) GetData() interface{} {
	return e.data
}

func (e *crdbErrorImpl) As(t *Type) Error {
	err := e
	for {
		if err.t == t {
			return e
		}
		lakeErr := AsLakeErrorType(err.Unwrap())
		if lakeErr == nil {
			return nil
		}
		err = lakeErr.(*crdbErrorImpl)
	}
}

func (e *crdbErrorImpl) getMessages(getMessage func(*crdbErrorImpl) string) []string {
	msgs := []string{}
	err := e
	ok := false
	for {
		msg := getMessage(err)
		if msg != "" {
			msgs = append(msgs, msg)
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
	errType := t
	var wrappedErr *crdbErrorImpl
	var wrappedRaw error
	rawMessage := message
	cfg.stackOffset += 2
	if err == nil {
		if enableStacktraces {
			wrappedRaw = cerror.NewWithDepth(int(cfg.stackOffset), rawMessage)
		} else {
			wrappedRaw = errors.New(message)
		}
	} else {
		if cast, ok := err.(*crdbErrorImpl); ok {
			err = cast.wrappedRaw
			wrappedErr = cast
			if t == Default { // inherit wrapped error's type
				errType = cast.GetType()
			}
		}
		if enableStacktraces {
			wrappedRaw = cerror.WrapWithDepth(int(cfg.stackOffset), err, rawMessage)
		} else {
			wrappedRaw = cerror.WithDetail(err, rawMessage)
		}
	}
	impl := &crdbErrorImpl{
		wrappedRaw: wrappedRaw,
		wrapped:    wrappedErr,
		msg:        rawMessage,
		data:       cfg.data,
		t:          errType,
	}
	return impl
}

var _ Error = (*crdbErrorImpl)(nil)
