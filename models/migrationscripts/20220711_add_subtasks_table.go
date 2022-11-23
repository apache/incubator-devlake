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
	"time"

	"github.com/apache/incubator-devlake/errors"
	commonArchived "github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/plugins/core"
)

var _ core.MigrationScript = (*addSubtasksTable)(nil)

type addSubtasksTable struct {
}

type subtasks20220711 struct {
	commonArchived.Model
	TaskID       uint64     `json:"task_id" gorm:"index"`
	SubtaskName  string     `json:"name" gorm:"column:name;index"`
	Number       int        `json:"number"`
	BeganAt      *time.Time `json:"beganAt"`
	FinishedAt   *time.Time `json:"finishedAt" gorm:"index"`
	SpentSeconds int64      `json:"spentSeconds"`
}

func (subtasks20220711) TableName() string {
	return "_devlake_subtasks"
}

func (addSubtasksTable) Up(basicRes core.BasicRes) errors.Error {
	return basicRes.GetDal().AutoMigrate(&subtasks20220711{})
}

func (addSubtasksTable) Version() uint64 {
	return 20220711000001
}

func (addSubtasksTable) Name() string {
	return "create subtask schema"
}
