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

type domainTextColumnCall struct {
	tableName  string
	columnName string
	columnType string
}

type domainTextColumnRecordingDal struct {
	dal.Dal
	calls []domainTextColumnCall
}

func (d *domainTextColumnRecordingDal) ModifyColumnType(tableName, columnName, columnType string) errors.Error {
	d.calls = append(d.calls, domainTextColumnCall{tableName, columnName, columnType})
	return nil
}

type domainTextColumnBasicRes struct {
	context.BasicRes
	database dal.Dal
}

func (r *domainTextColumnBasicRes) GetDal() dal.Dal {
	return r.database
}

func TestExpandDomainTextColumns(t *testing.T) {
	database := new(domainTextColumnRecordingDal)
	script := new(expandDomainTextColumns)

	if err := script.Up(&domainTextColumnBasicRes{database: database}); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	want := []domainTextColumnCall{
		{"cicd_tasks", "name", "text"},
		{"cicd_scopes", "name", "text"},
		{"cicd_releases", "name", "text"},
		{"cicd_releases", "display_title", "text"},
		{"cicd_deployment_commits", "name", "text"},
		{"cicd_deployment_commits", "subtask_name", "text"},
		{"cicd_deployment_commits", "ref_name", "text"},
		{"incidents", "component", "text"},
	}
	if !reflect.DeepEqual(database.calls, want) {
		t.Fatalf("ModifyColumnType calls = %#v, want %#v", database.calls, want)
	}
	if script.Version() != 20260819000001 {
		t.Fatalf("Version() = %d, want 20260819000001", script.Version())
	}
	if script.Name() != "expand domain text columns" {
		t.Fatalf("Name() = %q, want %q", script.Name(), "expand domain text columns")
	}

	for _, registeredScript := range All() {
		if registeredScript.Version() == script.Version() {
			return
		}
	}
	t.Fatalf("migration version %d is not registered", script.Version())
}
