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

package e2e

import (
	"testing"

	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/plugins/org/impl"
	"github.com/apache/incubator-devlake/plugins/org/tasks"
)

type scope struct {
	id        string
	tableName string
}

func (s scope) ScopeId() string {
	return s.id
}

func (s scope) ScopeName() string {
	panic("implement me")
}

func (s scope) TableName() string {
	return s.tableName
}

func TestProjectMappingDataFlow(t *testing.T) {
	dataflowTester := e2ehelper.NewDataFlowTester(t, "org", impl.Org{})
	scopes := []plugin.Scope{scope{
		id:        "bitbucket:BitbucketRepo:4:thenicetgp/lake",
		tableName: "boards",
	}, scope{
		id:        "github:GithubRepo:1:1",
		tableName: "repos",
	}}
	taskData := &tasks.TaskData{
		Options: &tasks.Options{
			ProjectMappings: []tasks.ProjectMapping{tasks.NewProjectMapping("my_project", scopes)},
		},
	}

	// import raw data table
	dataflowTester.FlushTabler(&crossdomain.ProjectMapping{})

	dataflowTester.Subtask(tasks.SetProjectMappingMeta, taskData)
	dataflowTester.VerifyTable(
		crossdomain.ProjectMapping{},
		"./snapshot_tables/project_mapping.csv",
		e2ehelper.ColumnWithRawData(
			"project_name",
			"table",
			"row_id",
		),
	)
}
