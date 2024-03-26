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
)

var _ plugin.MigrationScript = (*addSubtaskField)(nil)

type subtask20240322 struct {
	FinishedRecords int    `json:"finishedRecords"`
	Sequence        int    `json:"sequence"`
	IsCollector     bool   `json:"isCollector"`
	IsFailed        bool   `json:"isFailed"`
	Message         string `json:"message"`
}

func (subtask20240322) TableName() string {
	return "_devlake_subtasks"
}

type addSubtaskField struct{}

func (*addSubtaskField) Up(basicRes context.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(subtask20240322{})
}

func (*addSubtaskField) Version() uint64 {
	return 20240322111247
}

func (*addSubtaskField) Name() string {
	return "add some fields to _devlake_subtasks table"
}
