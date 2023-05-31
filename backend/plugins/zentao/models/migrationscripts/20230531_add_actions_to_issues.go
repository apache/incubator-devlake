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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"gorm.io/datatypes"
)

var _ plugin.MigrationScript = (*addActionsToIssues)(nil)

type Bug20230531 struct {
	Actions datatypes.JSON `json:"actions"`
}

func (Bug20230531) TableName() string {
	return "_tool_zentao_bugs"
}

type Task20230531 struct {
	Actions datatypes.JSON `json:"actions"`
}

func (Task20230531) TableName() string {
	return "_tool_zentao_tasks"
}

type Story20230531 struct {
	Actions datatypes.JSON `json:"actions"`
}

func (Story20230531) TableName() string {
	return "_tool_zentao_stories"
}

type addActionsToIssues struct{}

func (script *addActionsToIssues) Up(basicRes context.BasicRes) errors.Error {

	return migrationhelper.AutoMigrateTables(
		basicRes,
		&Bug20230531{},
		&Task20230531{},
		&Story20230531{},
	)
}

func (*addActionsToIssues) Version() uint64 {
	return 20230531000001
}

func (*addActionsToIssues) Name() string {
	return "add actions to bugs, stories and tasks"
}
