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
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/log"
)

// noopLogger is a Logger stub used by tests. It records the last Warn call so
// tests can assert that skip-mode logged an informative warning.
type noopLogger struct {
	warnCalls int
	lastWarn  string
}

func (l *noopLogger) IsLevelEnabled(log.LogLevel) bool { return false }
func (l *noopLogger) Printf(string, ...interface{})    {}
func (l *noopLogger) Log(log.LogLevel, string, ...interface{}) {}
func (l *noopLogger) Debug(string, ...interface{}) {}
func (l *noopLogger) Info(string, ...interface{})  {}
func (l *noopLogger) Warn(_ error, format string, a ...interface{}) {
	l.warnCalls++
	l.lastWarn = format
	_ = a
}
func (l *noopLogger) Error(error, string, ...interface{}) {}
func (l *noopLogger) Nested(string) log.Logger           { return l }
func (l *noopLogger) GetConfig() *log.LoggerConfig       { return &log.LoggerConfig{} }
func (l *noopLogger) SetStream(*log.LoggerStreamConfig)  {}

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildFilterJQL(tt.filterId, tt.incrementalJql); got != tt.want {
				t.Errorf("buildFilterJQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_parseIssuesResponse verifies that the shared jira:collectIssues response
// parser: (a) still returns the parsed issues on the happy path; (b) still
// surfaces the parse error by default when JIRA_SKIP_UNPARSEABLE_ISSUES is not
// set — preserving pre-#8949 behaviour; (c) logs a warning and returns an empty
// slice + nil error when JIRA_SKIP_UNPARSEABLE_ISSUES=true.
func Test_parseIssuesResponse(t *testing.T) {
	newRes := func(body string) *http.Response {
		return &http.Response{Body: io.NopCloser(bytes.NewBufferString(body))}
	}

	t.Run("happy path returns issues", func(t *testing.T) {
		logger := &noopLogger{}
		got, err := parseIssuesResponse(logger)(newRes(`{"issues":[{"id":"1"},{"id":"2"}]}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("want 2 issues, got %d", len(got))
		}
		if logger.warnCalls != 0 {
			t.Fatalf("did not expect a warn call on happy path, got %d", logger.warnCalls)
		}
	})

	t.Run("default returns error on unparseable body", func(t *testing.T) {
		t.Setenv(skipUnparseableIssuesEnvVar, "")
		logger := &noopLogger{}
		_, err := parseIssuesResponse(logger)(newRes(`<html>bad gateway maybe</html>`))
		if err == nil {
			t.Fatal("expected an error when body is unparseable and skip flag is unset")
		}
		if logger.warnCalls != 0 {
			t.Fatalf("did not expect a warn call in default mode, got %d", logger.warnCalls)
		}
	})

	t.Run("skip flag logs warning and returns empty slice", func(t *testing.T) {
		t.Setenv(skipUnparseableIssuesEnvVar, "true")
		logger := &noopLogger{}
		got, err := parseIssuesResponse(logger)(newRes(`<html>bad gateway maybe</html>`))
		if err != nil {
			t.Fatalf("expected no error when skip flag is set, got: %v", err)
		}
		if got == nil || len(got) != 0 {
			t.Fatalf("expected empty (non-nil) slice on skip, got %v", got)
		}
		if logger.warnCalls != 1 {
			t.Fatalf("expected one warn call, got %d", logger.warnCalls)
		}
	})

	t.Run("skip flag is case-insensitive", func(t *testing.T) {
		t.Setenv(skipUnparseableIssuesEnvVar, "TRUE")
		logger := &noopLogger{}
		_, err := parseIssuesResponse(logger)(newRes(`<html>bad gateway maybe</html>`))
		if err != nil {
			t.Fatalf("expected no error with TRUE, got: %v", err)
		}
	})

	t.Run("skip flag off with 'false' surfaces error", func(t *testing.T) {
		t.Setenv(skipUnparseableIssuesEnvVar, "false")
		logger := &noopLogger{}
		_, err := parseIssuesResponse(logger)(newRes(`{`))
		if err == nil {
			t.Fatal("expected error when flag is 'false'")
		}
	})
}
