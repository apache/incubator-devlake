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

package tap

import (
	"encoding/json"
	"github.com/apache/incubator-devlake/errors"
)

// SingerOutput raw data from a tap. One of these fields can ever be non-nil
type SingerOutput struct {
	state  *State
	record *Record[json.RawMessage]
}

// NewSingerTapOutput construct for Output. The src is the raw data coming from the tap
func NewSingerTapOutput(src json.RawMessage) (Output[json.RawMessage], errors.Error) {
	srcMap := map[string]any{}
	err := convert(src, &srcMap)
	if err != nil {
		return nil, err
	}
	ret := &SingerOutput{}
	srcType := srcMap["type"]
	if srcType == "STATE" {
		state := State{}
		if err = convert(src, &state); err != nil {
			return nil, err
		}
		ret.state = &state
	} else if srcType == "RECORD" {
		record := Record[json.RawMessage]{}
		if err = convert(src, &record); err != nil {
			return nil, err
		}
		ret.record = &record
	}
	return ret, nil
}

// AsTapState tries to convert the map object to a State. Returns false if it can't be done.
func (r *SingerOutput) AsTapState() (*State, bool) {
	return r.state, r.state != nil
}

// AsTapRecord tries to convert the map object to a Record. Returns false if it can't be done.
func (r *SingerOutput) AsTapRecord() (*Record[json.RawMessage], bool) {
	return r.record, r.record != nil
}

func convert(src json.RawMessage, dest any) errors.Error {
	if err := json.Unmarshal(src, dest); err != nil {
		return errors.Default.Wrap(err, "error converting type")
	}
	return nil
}

var _ Output[json.RawMessage] = (*SingerOutput)(nil)
