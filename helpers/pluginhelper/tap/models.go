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
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
	"gorm.io/datatypes"
	"time"
)

// abstract tap types
type (
	// Record the fields embedded in a singer-tap record. The specifics of the record are tap-implementation specific.
	Record[R any] struct {
		Type          string    `json:"type"`
		Stream        string    `json:"stream"`
		TimeExtracted time.Time `json:"time_extracted"`
		Record        R         `json:"record"`
	}
	// State the fields embedded in a singer-tap state. The specifics of the value are tap-implementation specific.
	State struct {
		Type  string         `json:"type"`
		Value map[string]any `json:"value"`
	}

	// RawState The raw-database version of State
	RawState struct {
		archived.GenericModel[string]
		Type  string
		Value datatypes.JSON
	}

	// RawOutput raw data from a tap. One of these fields can ever be non-nil
	RawOutput[R any] struct {
		state  *State
		record *Record[R]
	}
)

// TableName the table name
func (*RawState) TableName() string {
	return "_devlake_collector_state"
}

// NewTapState creates a new tap State variable
func NewTapState(v map[string]any) *State {
	return &State{
		Type:  "STATE",
		Value: v,
	}
}

// FromState converts State to RawState
func FromState(t *State) *RawState {
	b, err := json.Marshal(t.Value)
	if err != nil {
		panic(err)
	}
	return &RawState{
		Type:  t.Type,
		Value: b,
	}
}

// ToState converts RawState to State
func ToState(raw *RawState) *State {
	val := new(map[string]any)
	err := json.Unmarshal(raw.Value, val)
	if err != nil {
		panic(err)
	}
	return &State{
		Type:  raw.Type,
		Value: *val,
	}
}

// NewRawTapOutput construct for RawOutput. The src is the raw data coming from the tap
func NewRawTapOutput[R any](src json.RawMessage) (*RawOutput[R], errors.Error) {
	srcMap := map[string]any{}
	err := convert(src, &srcMap)
	if err != nil {
		return nil, err
	}
	ret := &RawOutput[R]{}
	srcType := srcMap["type"]
	if srcType == "STATE" {
		state := State{}
		if err = convert(src, &state); err != nil {
			return nil, err
		}
		ret.state = &state
	} else if srcType == "RECORD" {
		record := Record[R]{}
		if err = convert(src, &record); err != nil {
			return nil, err
		}
		ret.record = &record
	}
	return ret, nil
}

// AsTapState tries to convert the map object to a State. Returns false if it can't be done.
func (r *RawOutput[R]) AsTapState() (*State, bool) {
	return r.state, r.state != nil
}

// AsTapRecord tries to convert the map object to a Record. Returns false if it can't be done.
func (r *RawOutput[R]) AsTapRecord() (*Record[R], bool) {
	return r.record, r.record != nil
}

func convert(src json.RawMessage, dest any) errors.Error {
	if err := json.Unmarshal(src, dest); err != nil {
		return errors.Default.Wrap(err, "error converting type")
	}
	return nil
}

var _ core.Tabler = (*RawState)(nil)
