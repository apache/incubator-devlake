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
	"encoding/json"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
)

var _ core.MigrationScript = (*addSubtaskToTaskTable)(nil)

type tasks20220601 struct {
	Subtasks json.RawMessage `json:"subtasks"`
}

func (tasks20220601) TableName() string {
	return "_devlake_tasks"
}

type addSubtaskToTaskTable struct{}

func (*addSubtaskToTaskTable) Up(basicRes core.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&tasks20220601{})
}

func (*addSubtaskToTaskTable) Version() uint64 {
	return 20220601000005
}

func (*addSubtaskToTaskTable) Name() string {
	return "add column `subtasks` at _devlake_tasks"
}
