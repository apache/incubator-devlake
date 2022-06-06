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

package models

import (
	"time"

	"github.com/apache/incubator-devlake/models/common"
	"gorm.io/datatypes"
)

const (
	TASK_CREATED   = "TASK_CREATED"
	TASK_RUNNING   = "TASK_RUNNING"
	TASK_COMPLETED = "TASK_COMPLETED"
	TASK_FAILED    = "TASK_FAILED"
)

type TaskProgressDetail struct {
	TotalSubTasks    int    `json:"totalSubTasks"`
	FinishedSubTasks int    `json:"finishedSubTasks"`
	TotalRecords     int    `json:"totalRecords"`
	FinishedRecords  int    `json:"finishedRecords"`
	SubTaskName      string `json:"subTaskName"`
	SubTaskNumber    int    `json:"subTaskNumber"`
}

type Task struct {
	common.Model
	Plugin         string              `json:"plugin" gorm:"index"`
	Subtasks       datatypes.JSON      `json:"subtasks"`
	Options        datatypes.JSON      `json:"options"`
	Status         string              `json:"status"`
	Message        string              `json:"message"`
	Progress       float32             `json:"progress"`
	ProgressDetail *TaskProgressDetail `json:"progressDetail" gorm:"-"`

	FailedSubTask string     `json:"failedSubTask"`
	PipelineId    uint64     `json:"pipelineId" gorm:"index"`
	PipelineRow   int        `json:"pipelineRow"`
	PipelineCol   int        `json:"pipelineCol"`
	BeganAt       *time.Time `json:"beganAt"`
	FinishedAt    *time.Time `json:"finishedAt" gorm:"index"`
	SpentSeconds  int        `json:"spentSeconds"`
}

type NewTask struct {
	// Plugin name
	Plugin      string                 `json:"plugin" binding:"required"`
	Subtasks    []string               `json:"subtasks"`
	Options     map[string]interface{} `json:"options"`
	PipelineId  uint64                 `json:"-"`
	PipelineRow int                    `json:"-"`
	PipelineCol int                    `json:"-"`
}

func (Task) TableName() string {
	return "_devlake_tasks"
}
