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

	"github.com/apache/incubator-devlake/core/models/common"
)

const (
	TASK_CREATED   = "TASK_CREATED"
	TASK_RERUN     = "TASK_RERUN"
	TASK_RESUME    = "TASK_RESUME"
	TASK_RUNNING   = "TASK_RUNNING"
	TASK_COMPLETED = "TASK_COMPLETED"
	TASK_FAILED    = "TASK_FAILED"
	TASK_CANCELLED = "TASK_CANCELLED"
	TASK_PARTIAL   = "TASK_PARTIAL"
)

var PendingTaskStatus = []string{TASK_CREATED, TASK_RERUN, TASK_RUNNING}

type TaskProgressDetail struct {
	TotalSubTasks        int    `json:"totalSubTasks"`
	FinishedSubTasks     int    `json:"finishedSubTasks"`
	TotalRecords         int    `json:"totalRecords"`
	FinishedRecords      int    `json:"finishedRecords"`
	SubTaskName          string `json:"subTaskName"`
	SubTaskNumber        int    `json:"subTaskNumber"`
	CollectSubtaskNumber int    `json:"collectSubtaskNumber"`
	OtherSubtaskNumber   int    `json:"otherSubtaskNumber"`
}

type NewTask struct {
	// Plugin name
	*PipelineTask
	PipelineId  uint64 `json:"-"`
	PipelineRow int    `json:"-"`
	PipelineCol int    `json:"-"`
	IsRerun     bool   `json:"-"`
}

type Task struct {
	common.Model
	Plugin         string                 `json:"plugin" gorm:"index"`
	Subtasks       []string               `json:"subtasks" gorm:"type:json;serializer:json"`
	Options        map[string]interface{} `json:"options" gorm:"serializer:encdec"`
	Status         string                 `json:"status"`
	Message        string                 `json:"message"`
	ErrorName      string                 `json:"errorName"`
	Progress       float32                `json:"progress"`
	ProgressDetail *TaskProgressDetail    `json:"progressDetail" gorm:"-"`

	FailedSubTask string     `json:"failedSubTask"`
	PipelineId    uint64     `json:"pipelineId" gorm:"index"`
	PipelineRow   int        `json:"pipelineRow"`
	PipelineCol   int        `json:"pipelineCol"`
	BeganAt       *time.Time `json:"beganAt"`
	FinishedAt    *time.Time `json:"finishedAt" gorm:"index"`
	SpentSeconds  int        `json:"spentSeconds"`
}

func (Task) TableName() string {
	return "_devlake_tasks"
}

type Subtask struct {
	common.Model
	TaskID          uint64     `json:"task_id" gorm:"index"`
	Name            string     `json:"name" gorm:"index"`
	Number          int        `json:"number"`
	BeganAt         *time.Time `json:"beganAt"`
	FinishedAt      *time.Time `json:"finishedAt" gorm:"index"`
	SpentSeconds    int64      `json:"spentSeconds"`
	FinishedRecords int        `json:"finishedRecords"`
	Sequence        int        `json:"sequence"`
	IsCollector     bool       `json:"isCollector"`
	IsFailed        bool       `json:"isFailed"`
	Message         string     `json:"message"`
}

func (Subtask) TableName() string {
	return "_devlake_subtasks"
}

type SubtaskDetails struct {
	ID              uint64     `json:"id"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	TaskID          uint64     `json:"taskId"`
	Name            string     `json:"name"`
	Number          int        `json:"number"`
	BeganAt         *time.Time `json:"beganAt"`
	FinishedAt      *time.Time `json:"finishedAt"`
	SpentSeconds    int64      `json:"spentSeconds"`
	FinishedRecords int        `json:"finishedRecords"`
	Sequence        int        `json:"sequence"`
	IsCollector     bool       `json:"isCollector"`
	IsFailed        bool       `json:"isFailed"`
	Message         string     `json:"message"`
}

type SubtasksInfo struct {
	ID                uint64            `json:"id"`
	PipelineID        uint64            `json:"pipelineId"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
	BeganAt           *time.Time        `json:"beganAt"`
	FinishedAt        *time.Time        `json:"finishedAt"`
	Plugin            string            `json:"plugin"`
	Options           any               `json:"options"`
	Status            string            `json:"status"`
	FailedSubTask     string            `json:"failedSubTask"`
	Message           string            `json:"message"`
	ErrorName         string            `json:"errorName"`
	SpentSeconds      int               `json:"spentSeconds"`
	SubtaskDetails    []*SubtaskDetails `json:"subtaskDetails"`
	TotalTransform    int64             `json:"totalTransform"`
	FinishedTransform int64             `json:"finishedTransform"`
}

type SubTasksOuput struct {
	SubtasksInfo   []SubtasksInfo `json:"subtasks"`
	CompletionRate float64        `json:"completionRate"`
	Status         string         `json:"status"`
	Count          int64          `json:"count"`
}
