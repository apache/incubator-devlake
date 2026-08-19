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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"github.com/apache/incubator-devlake/helpers/unithelper"
	mocklog "github.com/apache/incubator-devlake/mocks/core/log"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/stretchr/testify/mock"
)

func Test_buildJQL(t *testing.T) {
	base := time.Date(2021, 2, 3, 4, 5, 6, 7, time.UTC)
	timeAfter := base
	add48 := base.Add(48 * time.Hour)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	type args struct {
		since    *time.Time
		location *time.Location
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test incremental",
			args: args{
				since:    &add48,
				location: loc,
			},
			want: "updated >= '2021/02/05 12:05' ORDER BY created ASC",
		},
		{
			name: "test incremental",
			args: args{
				since: &timeAfter,
			},
			want: "updated >= '2021/02/02 04:05' ORDER BY created ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildJQL(*tt.args.since, tt.args.location); got != tt.want {
				t.Errorf("buildJQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildFilterJQL(t *testing.T) {
	tests := []struct {
		name           string
		filterId       string
		extraJql       string
		incrementalJql string
		want           string
	}{
		{
			name:           "full sync with filter",
			filterId:       "12345",
			incrementalJql: "ORDER BY created ASC",
			want:           "filter = 12345 ORDER BY created ASC",
		},
		{
			name:           "incremental sync with filter",
			filterId:       "12345",
			incrementalJql: "updated >= '2021/02/05 12:05' ORDER BY created ASC",
			want:           "filter = 12345 AND updated >= '2021/02/05 12:05' ORDER BY created ASC",
		},
		{
			name:           "empty filter id falls back to incremental only",
			filterId:       "",
			incrementalJql: "ORDER BY created ASC",
			want:           "ORDER BY created ASC",
		},
		{
			name:           "empty filter id with incremental clause",
			filterId:       "",
			incrementalJql: "updated >= '2024/01/01 00:00' ORDER BY created ASC",
			want:           "updated >= '2024/01/01 00:00' ORDER BY created ASC",
		},
		{
			name:           "extra jql with filter full sync",
			filterId:       "12345",
			extraJql:       `project = "MyComponent"`,
			incrementalJql: "ORDER BY created ASC",
			want:           `filter = 12345 AND (project = "MyComponent") ORDER BY created ASC`,
		},
		{
			name:           "extra jql with filter incremental sync",
			filterId:       "12345",
			extraJql:       `project = "MyComponent"`,
			incrementalJql: "updated >= '2024/01/01 00:00' ORDER BY created ASC",
			want:           `filter = 12345 AND (project = "MyComponent") AND updated >= '2024/01/01 00:00' ORDER BY created ASC`,
		},
		{
			name:           "extra jql without filter",
			filterId:       "",
			extraJql:       `project = "MyComponent"`,
			incrementalJql: "ORDER BY created ASC",
			want:           `(project = "MyComponent") ORDER BY created ASC`,
		},
		{
			name:           "extra jql without filter, incremental sync",
			filterId:       "",
			extraJql:       `project = "MyComponent"`,
			incrementalJql: "updated >= '2024/01/01 00:00' ORDER BY created ASC",
			want:           `(project = "MyComponent") AND updated >= '2024/01/01 00:00' ORDER BY created ASC`,
		},
		{
			name:           "extra jql with OR operator is parenthesized",
			filterId:       "12345",
			extraJql:       `project = "A" OR project = "B"`,
			incrementalJql: "ORDER BY created ASC",
			want:           `filter = 12345 AND (project = "A" OR project = "B") ORDER BY created ASC`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildFilterJQL(tt.filterId, tt.extraJql, tt.incrementalJql); got != tt.want {
				t.Errorf("buildFilterJQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseIssuesResponse(t *testing.T) {
	const twoIssues = `{"issues":[{"id":"1"},{"id":"2"}]}`
	// A proxy returning an HTML error page under a 200, the case reported in #8949.
	const htmlPage = `<html><head><title>502 Bad Gateway</title></head></html>`
	// A response truncated mid-flight, so the JSON never terminates.
	const truncated = `{"issues":[{"id":"1"},{"id":`

	makeResponse := func(body string) *http.Response {
		return &http.Response{
			Body:    io.NopCloser(strings.NewReader(body)),
			Request: &http.Request{URL: &url.URL{Scheme: "https", Host: "jira.example.com", Path: "/rest/api/2/search"}},
		}
	}

	tests := []struct {
		name            string
		body            string
		skipUnparseable bool
		wantCount       int
		wantErr         bool
		// wantSkipped marks the cases that took the skip path, which must yield an
		// empty non-nil slice so the collector records "this page had no rows" rather
		// than treating the result as absent.
		wantSkipped bool
	}{
		{
			name:      "valid page is parsed",
			body:      twoIssues,
			wantCount: 2,
		},
		{
			name:            "valid page is parsed when skipping is enabled",
			body:            twoIssues,
			skipUnparseable: true,
			wantCount:       2,
		},
		{
			name:    "html body fails by default",
			body:    htmlPage,
			wantErr: true,
		},
		{
			name:            "html body is skipped when enabled",
			body:            htmlPage,
			skipUnparseable: true,
			wantCount:       0,
			wantSkipped:     true,
		},
		{
			name:    "truncated body fails by default",
			body:    truncated,
			wantErr: true,
		},
		{
			name:            "truncated body is skipped when enabled",
			body:            truncated,
			skipUnparseable: true,
			wantCount:       0,
			wantSkipped:     true,
		},
		{
			name:      "valid page without an issues key yields no issues",
			body:      `{"total":0}`,
			wantCount: 0,
		},
	}

	// unithelper.DummyLogger stubs Warn with two arguments, but the real signature is
	// Warn(err, format, a ...interface{}), which mockery records as three. Add the
	// variadic form so the skip path can log.
	newLogger := func() *mocklog.Logger {
		logger := unithelper.DummyLogger()
		logger.On("Warn", mock.Anything, mock.Anything, mock.Anything).Maybe()
		return logger
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIssuesResponse(makeResponse(tt.body), tt.skipUnparseable, newLogger())
			if (err != nil) != tt.wantErr {
				t.Errorf("parseIssuesResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("parseIssuesResponse() returned %d issues, want %d", len(got), tt.wantCount)
			}
			if tt.wantSkipped && got == nil {
				t.Error("parseIssuesResponse() returned nil for a skipped page, want an empty slice")
			}
		})
	}
}

func Test_responseUrl(t *testing.T) {
	withUrl := &http.Response{
		Request: &http.Request{URL: &url.URL{Scheme: "https", Host: "jira.example.com", Path: "/rest/api/2/search"}},
	}
	if got, want := responseUrl(withUrl), "https://jira.example.com/rest/api/2/search"; got != want {
		t.Errorf("responseUrl() = %v, want %v", got, want)
	}
	if got, want := responseUrl(&http.Response{}), "unknown url"; got != want {
		t.Errorf("responseUrl() with no request = %v, want %v", got, want)
	}
	if got, want := responseUrl(nil), "unknown url"; got != want {
		t.Errorf("responseUrl(nil) = %v, want %v", got, want)
	}
}

func Test_renderExtraJQL(t *testing.T) {
	// Task 3.2 (1): wire the third parameter to set ProjectName on both structs.
	makeData := func(boardId uint64, boardName string, projectName string) *JiraTaskData {
		return &JiraTaskData{
			Options:     &JiraOptions{BoardId: boardId, ProjectName: projectName},
			Board:       &models.JiraBoard{BoardId: boardId, Name: boardName},
			ProjectName: projectName,
		}
	}

	tests := []struct {
		name    string
		tmpl    string
		data    *JiraTaskData
		want    string
		wantErr bool
	}{
		{
			name: "static JQL passes through unchanged",
			tmpl: `project = "MyProject"`,
			data: makeData(1, "My Board", ""),
			want: `project = "MyProject"`,
		},
		{
			name: "BoardName substitution",
			tmpl: `project = "{{.BoardName}}"`,
			data: makeData(42, "Team Alpha", ""),
			want: `project = "Team Alpha"`,
		},
		{
			name: "BoardId substitution",
			tmpl: `cf[10001] = {{.BoardId}}`,
			data: makeData(99, "Some Board", ""),
			want: `cf[10001] = 99`,
		},
		{
			name: "nil Board falls back to empty BoardName",
			tmpl: `project = "{{.BoardName}}"`,
			data: &JiraTaskData{Options: &JiraOptions{BoardId: 1}, Board: nil},
			want: `project = ""`,
		},
		{
			name:    "invalid template returns error",
			tmpl:    `project = "{{.Unclosed"`,
			data:    makeData(1, "My Board", ""),
			wantErr: true,
		},
		{
			name:    "unknown field returns error (missingkey=error)",
			tmpl:    `component = "{{.Typo}}"`,
			data:    makeData(1, "My Board", ""),
			wantErr: true,
		},
		// Task 3.2 (2): ProjectName table-driven test cases.
		{
			name: "ProjectName substitution",
			tmpl: `labels = "{{.ProjectName}}"`,
			data: makeData(1, "My Board", "My DevLake Project"),
			want: `labels = "My DevLake Project"`,
		},
		{
			name: "empty ProjectName renders empty string",
			tmpl: `labels = "{{.ProjectName}}"`,
			data: makeData(1, "My Board", ""),
			want: `labels = ""`,
		},
		{
			name: "combined ProjectName and BoardName substitution",
			tmpl: `component = "{{.ProjectName}}" AND owner = "{{.BoardName}}"`,
			data: makeData(1, "My Board", "My Project"),
			want: `component = "My Project" AND owner = "My Board"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := renderExtraJQL(tt.tmpl, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("renderExtraJQL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("renderExtraJQL() = %v, want %v", got, tt.want)
			}
		})
	}

	// Task 3.2 (3): Property 1 — template variable round-trip.
	// Validates: Requirements 1.2, 1.4, 5.4
	// Tag: Feature: jira-extrajql-projectname, Property 1: template variable round-trip
	t.Run("property: triple-variable round-trip", func(t *testing.T) {
		f := func(projectName, boardName string, boardId uint64) bool {
			data := &JiraTaskData{
				Options:     &JiraOptions{BoardId: boardId, ProjectName: projectName},
				Board:       &models.JiraBoard{BoardId: boardId, Name: boardName},
				ProjectName: projectName,
			}
			got, err := renderExtraJQL(`{{.ProjectName}}-{{.BoardName}}-{{.BoardId}}`, data)
			if err != nil {
				return false
			}
			want := fmt.Sprintf("%s-%s-%d", projectName, boardName, boardId)
			return got == want
		}
		if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 1 (triple-variable round-trip) failed: %v", err)
		}
	})

	// Task 3.3: Property 2 — static JQL passthrough.
	// Validates: Requirements 4.3
	// Tag: Feature: jira-extrajql-projectname, Property 2: static JQL passthrough
	t.Run("property: static JQL passthrough", func(t *testing.T) {
		f := func(s string) bool {
			if strings.Contains(s, "{{") {
				return true // skip: not a static template
			}
			data := &JiraTaskData{
				Options: &JiraOptions{BoardId: 1},
				Board:   &models.JiraBoard{},
			}
			got, err := renderExtraJQL(s, data)
			return err == nil && got == s
		}
		if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
			t.Errorf("Property 2 (static JQL passthrough) failed: %v", err)
		}
	})
}

// Test_renderExtraJQL_Property1_RoundTrip is a top-level property-based test for
// the template variable round-trip property.
// Validates: Requirements 1.2, 1.4, 5.4
func Test_renderExtraJQL_Property1_RoundTrip(t *testing.T) {
	f := func(projectName, boardName string, boardId uint64) bool {
		data := &JiraTaskData{
			Options:     &JiraOptions{BoardId: boardId},
			Board:       &models.JiraBoard{BoardId: boardId, Name: boardName},
			ProjectName: projectName,
		}
		got, err := renderExtraJQL(`{{.ProjectName}}-{{.BoardName}}-{{.BoardId}}`, data)
		if err != nil {
			return false
		}
		want := fmt.Sprintf("%s-%s-%d", projectName, boardName, boardId)
		return got == want
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 1 (round-trip) failed: %v", err)
	}
}

// Test_renderExtraJQL_Property2_StaticPassthrough is a top-level property-based test for
// the static JQL passthrough property.
// Validates: Requirements 4.3
func Test_renderExtraJQL_Property2_StaticPassthrough(t *testing.T) {
	f := func(s string) bool {
		if strings.Contains(s, "{{") {
			return true // skip: not a static template
		}
		data := &JiraTaskData{
			Options: &JiraOptions{BoardId: 1},
			Board:   &models.JiraBoard{},
		}
		got, err := renderExtraJQL(s, data)
		return err == nil && got == s
	}
	if err := quick.Check(f, &quick.Config{MaxCount: 100}); err != nil {
		t.Errorf("Property 2 (static passthrough) failed: %v", err)
	}
}
