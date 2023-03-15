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

package bridge

import (
	"encoding/json"

	"github.com/apache/incubator-devlake/core/errors"
)

type (
	CallResult struct {
		Results []byte
		Err     errors.Error
	}
	StreamResult = CallResult
	MethodStream struct {
		// send (not supported for now)
		outbound chan<- any
		// receive
		inbound <-chan *StreamResult
	}
)

func NewCallResult(results []byte, err errors.Error) *CallResult {
	return &CallResult{
		Results: results,
		Err:     err,
	}
}

func (m *CallResult) Get(target any) errors.Error {
	if m.Err != nil {
		return m.Err
	}
	err := json.Unmarshal(m.Results, &target)
	if err != nil {
		return errors.Convert(err)
	}
	return nil
}

func NewStreamResult(results []byte, err errors.Error) *StreamResult {
	return &StreamResult{
		Results: results,
		Err:     err,
	}
}

func (m *MethodStream) Receive() <-chan *StreamResult {
	return m.inbound
}
