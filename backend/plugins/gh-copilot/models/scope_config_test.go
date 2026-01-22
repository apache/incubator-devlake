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
	"testing"
	"time"
)

func TestGhCopilotScopeConfig_TableName(t *testing.T) {
	sc := GhCopilotScopeConfig{}
	expected := "_tool_copilot_scope_configs"
	if sc.TableName() != expected {
		t.Errorf("TableName() = %v, want %v", sc.TableName(), expected)
	}
}

func TestGhCopilotScopeConfig_BeforeSave_BaselinePeriodDays(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"zero defaults to 90", 0, 90},
		{"negative defaults to 90", -10, 90},
		{"below minimum (5) defaults to 90", 5, 90},
		{"below minimum (6) defaults to 90", 6, 90},
		{"minimum valid (7) unchanged", 7, 7},
		{"typical value (30) unchanged", 30, 30},
		{"default value (90) unchanged", 90, 90},
		{"high value (180) unchanged", 180, 180},
		{"maximum valid (365) unchanged", 365, 365},
		{"above maximum (400) capped to 365", 400, 365},
		{"far above maximum (1000) capped to 365", 1000, 365},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &GhCopilotScopeConfig{
				BaselinePeriodDays: tt.input,
			}
			err := sc.BeforeSave()
			if err != nil {
				t.Errorf("BeforeSave() error = %v, want nil", err)
			}
			if sc.BaselinePeriodDays != tt.expected {
				t.Errorf("BaselinePeriodDays = %v, want %v", sc.BaselinePeriodDays, tt.expected)
			}
		})
	}
}

func TestGhCopilotScopeConfig_BeforeSave_PreservesOtherFields(t *testing.T) {
	now := time.Now()
	sc := &GhCopilotScopeConfig{
		ImplementationDate: &now,
		BaselinePeriodDays: 60,
	}

	err := sc.BeforeSave()
	if err != nil {
		t.Errorf("BeforeSave() error = %v, want nil", err)
	}

	// ImplementationDate should be preserved
	if sc.ImplementationDate == nil || !sc.ImplementationDate.Equal(now) {
		t.Error("ImplementationDate was modified unexpectedly")
	}

	// BaselinePeriodDays should be unchanged (valid value)
	if sc.BaselinePeriodDays != 60 {
		t.Errorf("BaselinePeriodDays = %v, want 60", sc.BaselinePeriodDays)
	}
}

func TestGhCopilotScopeConfig_GetConnectionId(t *testing.T) {
	sc := GhCopilotScopeConfig{}
	sc.ConnectionId = 42

	if sc.GetConnectionId() != 42 {
		t.Errorf("GetConnectionId() = %v, want 42", sc.GetConnectionId())
	}
}
