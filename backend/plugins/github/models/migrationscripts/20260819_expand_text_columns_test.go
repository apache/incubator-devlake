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

package migrationscripts

import (
	"reflect"
	"testing"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
)

type githubTextColumnCall struct {
	tableName  string
	columnName string
	columnType string
}

type githubTextColumnRecordingDal struct {
	dal.Dal
	calls []githubTextColumnCall
}

func (d *githubTextColumnRecordingDal) ModifyColumnType(tableName, columnName, columnType string) errors.Error {
	d.calls = append(d.calls, githubTextColumnCall{tableName, columnName, columnType})
	return nil
}

type githubTextColumnBasicRes struct {
	context.BasicRes
	database dal.Dal
}

func (r *githubTextColumnBasicRes) GetDal() dal.Dal {
	return r.database
}

func TestExpandGithubTextColumns(t *testing.T) {
	database := new(githubTextColumnRecordingDal)
	script := new(expandGithubTextColumns)

	if err := script.Up(&githubTextColumnBasicRes{database: database}); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	want := []githubTextColumnCall{
		{"_tool_github_jobs", "name", "text"},
		{"_tool_github_jobs", "runner_name", "text"},
		{"_tool_github_jobs", "environment", "text"},
		{"_tool_github_runs", "head_branch", "text"},
		{"_tool_github_runs", "path", "text"},
		{"_tool_github_pull_requests", "head_ref", "text"},
		{"_tool_github_pull_requests", "base_ref", "text"},
		{"_tool_github_pull_requests", "author_name", "text"},
		{"_tool_github_pull_requests", "merged_by_name", "text"},
		{"_tool_github_deployments", "environment", "text"},
		{"_tool_github_deployments", "ref_name", "text"},
		{"_tool_github_deployments", "description", "text"},
		{"_tool_github_accounts", "company", "text"},
		{"_tool_github_accounts", "name", "text"},
	}
	if !reflect.DeepEqual(database.calls, want) {
		t.Fatalf("ModifyColumnType calls = %#v, want %#v", database.calls, want)
	}
	if script.Version() != 20260819000000 {
		t.Fatalf("Version() = %d, want 20260819000000", script.Version())
	}
	if script.Name() != "expand GitHub text columns" {
		t.Fatalf("Name() = %q, want %q", script.Name(), "expand GitHub text columns")
	}

	for _, registeredScript := range All() {
		if registeredScript.Version() == script.Version() {
			return
		}
	}
	t.Fatalf("migration version %d is not registered", script.Version())
}
