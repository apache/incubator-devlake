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

package tasks

import (
	"encoding/json"
	"testing"

	"github.com/apache/devlake/core/models/domainlayer/ticket"
	"github.com/apache/devlake/plugins/clickup/models"
)

func TestParseSprintDate(t *testing.T) {
	tests := []struct {
		name  string
		token string
		want  string // "" means nil (no reliable date)
	}{
		{"Toto M/D/YY", "7/6/26", "2026-07-06"},
		{"Toto M/D/YY day>12", "7/19/26", "2026-07-19"},
		{"Blockchain D/M/YY day>12", "29/6/26", "2026-06-29"},
		{"Blockchain D/M/YY", "26/7/26", "2026-07-26"},
		{"ambiguous both<=12 defaults M/D", "3/4/26", "2026-03-04"},
		{"no year -> nil", "7/6", ""},
		{"empty -> nil", "", ""},
		{"invalid month both>12 -> nil", "13/13/26", ""},
		{"non-numeric -> nil", "a/b/c", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSprintDate(tt.token)
			if tt.want == "" {
				if got != nil {
					t.Fatalf("expected nil, got %v", got)
				}
				return
			}
			if got == nil {
				t.Fatalf("expected %s, got nil", tt.want)
			}
			if g := got.Format("2006-01-02"); g != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, g)
			}
		})
	}
}

func TestSprintDetector(t *testing.T) {
	det, err := newSprintDetector(nil) // nil -> default pattern
	if err != nil {
		t.Fatalf("newSprintDetector: %v", err)
	}
	tests := []struct {
		name       string
		listName   string
		wantSprint bool
		wantStart  string
	}{
		{"toto sprint with dates", "v4.3.0  Sprint 40 (7/6/26 - 7/19/26)", true, "2026-07-06"},
		{"blockchain sprint D/M", "Sprint 25 (29/6/26 - 26/7/26)", true, "2026-06-29"},
		{"sprint no year -> no dates", "Sprint 31 (7/6 - 7/19)", true, ""},
		{"backlog is not a sprint", "Backlog", false, ""},
		{"qa bugs is not a sprint", "QA Bugs", false, ""},
		{"devops is not a sprint", "DevOps", false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := det.detect(tt.listName)
			if tt.wantSprint != (got != nil) {
				t.Fatalf("wantSprint=%v, got=%v", tt.wantSprint, got != nil)
			}
			if got == nil {
				return
			}
			if got.name != tt.listName {
				t.Fatalf("sprint name = %q, want %q", got.name, tt.listName)
			}
			if tt.wantStart == "" {
				if got.start != nil {
					t.Fatalf("expected nil start, got %v", got.start)
				}
			} else if got.start == nil || got.start.Format("2006-01-02") != tt.wantStart {
				t.Fatalf("start = %v, want %s", got.start, tt.wantStart)
			}
		})
	}
}

func TestStatusFromType(t *testing.T) {
	tests := []struct {
		statusType string
		want       string
	}{
		{"open", ticket.TODO},
		{"unstarted", ticket.TODO},
		{"custom", ticket.IN_PROGRESS},
		{"done", ticket.DONE},
		{"closed", ticket.DONE},
		{"CLOSED", ticket.DONE}, // case-insensitive
		{"mystery", ticket.OTHER},
	}
	for _, tt := range tests {
		t.Run(tt.statusType, func(t *testing.T) {
			if got := statusFromType(tt.statusType); got != tt.want {
				t.Fatalf("statusFromType(%q) = %q, want %q", tt.statusType, got, tt.want)
			}
		})
	}
}

func TestStatusMapperOverride(t *testing.T) {
	// Toto tags "to do" and "re-open" as ClickUp type "custom", which the
	// type-default would file as IN_PROGRESS. The scope-config override must win.
	sc := &models.ClickUpScopeConfig{
		IssueStatusTodo:       []string{"to do", "re-open"},
		IssueStatusInProgress: []string{"in development"},
		IssueStatusDone:       []string{"deployed", "Closed"},
	}
	m := newStatusMapper(sc)
	tests := []struct {
		rawStatus  string
		statusType string
		want       string
	}{
		{"to do", "custom", ticket.TODO},           // override beats type default
		{"re-open", "custom", ticket.TODO},         // override beats type default
		{"in development", "custom", ticket.IN_PROGRESS},
		{"deployed", "done", ticket.DONE},
		{"CLOSED", "closed", ticket.DONE},          // override match is case-insensitive
		{"in code review", "custom", ticket.IN_PROGRESS}, // not listed -> type default
	}
	for _, tt := range tests {
		t.Run(tt.rawStatus, func(t *testing.T) {
			if got := m.statusOf(tt.rawStatus, tt.statusType); got != tt.want {
				t.Fatalf("statusOf(%q,%q) = %q, want %q", tt.rawStatus, tt.statusType, got, tt.want)
			}
		})
	}
}

func TestStatusMapperNoConfigFallsBackToType(t *testing.T) {
	m := newStatusMapper(nil)
	if got := m.statusOf("anything", "open"); got != ticket.TODO {
		t.Fatalf("expected TODO from type, got %q", got)
	}
}

func TestIssueTypeMatcher(t *testing.T) {
	sc := &models.ClickUpScopeConfig{
		IssueTypeIncident:    "(?i)incident",
		IssueTypeBug:         "(?i)^bug$",
		IssueTypeRequirement: "(?i)(feature|story)",
	}
	m, err := newIssueTypeMatcher(sc)
	if err != nil {
		t.Fatalf("newIssueTypeMatcher: %v", err)
	}
	tests := []struct {
		taskType string
		want     string
	}{
		{"incident", ticket.INCIDENT},
		{"Bug", ticket.BUG},
		{"feature", ticket.REQUIREMENT},
		{"", ticket.REQUIREMENT}, // no match -> default REQUIREMENT
	}
	for _, tt := range tests {
		t.Run(tt.taskType, func(t *testing.T) {
			if got := m.typeOf(tt.taskType); got != tt.want {
				t.Fatalf("typeOf(%q) = %q, want %q", tt.taskType, got, tt.want)
			}
		})
	}
}

func TestIssueTypeMatcherDefaultBug(t *testing.T) {
	// nil config still classifies ClickUp's built-in "Bug" via defaultBugPattern.
	m, err := newIssueTypeMatcher(nil)
	if err != nil {
		t.Fatalf("newIssueTypeMatcher: %v", err)
	}
	if got := m.typeOf("bug"); got != ticket.BUG {
		t.Fatalf("typeOf(bug) = %q, want BUG", got)
	}
	if got := m.typeOf("anything-else"); got != ticket.REQUIREMENT {
		t.Fatalf("typeOf(anything-else) = %q, want REQUIREMENT", got)
	}
}

func TestStoryPointOf(t *testing.T) {
	fp := func(f float64) *float64 { return &f }
	native := &ClickUpApiTask{Points: fp(8)}
	// default (no field configured) uses native points
	if got := storyPointOf(native, nil); got == nil || *got != 8 {
		t.Fatalf("native points: got %v, want 8", got)
	}
	// configured custom field wins
	custom := &ClickUpApiTask{
		Points: fp(8),
		CustomFields: []clickUpCustomField{
			{Name: "LOE", Value: json.RawMessage(`"13"`)},
		},
	}
	sc := &models.ClickUpScopeConfig{StoryPointField: "loe"} // case-insensitive
	if got := storyPointOf(custom, sc); got == nil || *got != 13 {
		t.Fatalf("custom field: got %v, want 13", got)
	}
	// configured field absent on task -> nil (does not fall back to native)
	if got := storyPointOf(native, sc); got != nil {
		t.Fatalf("missing custom field: got %v, want nil", got)
	}
}

func TestNumericValue(t *testing.T) {
	tests := []struct {
		raw  string
		want *float64
	}{
		{`5`, ptr(5)},
		{`"8"`, ptr(8)},
		{`3.5`, ptr(3.5)},
		{``, nil},
		{`null`, nil},
		{`"abc"`, nil},
	}
	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			got := numericValue(json.RawMessage(tt.raw))
			if (tt.want == nil) != (got == nil) {
				t.Fatalf("raw %q: got %v, want %v", tt.raw, got, tt.want)
			}
			if tt.want != nil && *got != *tt.want {
				t.Fatalf("raw %q: got %v, want %v", tt.raw, *got, *tt.want)
			}
		})
	}
}

func ptr(f float64) *float64 { return &f }
