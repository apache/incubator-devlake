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

type GenericPipelineTask[T any] struct {
	Plugin   string   `json:"plugin" binding:"required"`
	Subtasks []string `json:"subtasks"`
	Options  T        `json:"options"`
}

// PipelineTask represents a smallest unit of execution inside a PipelinePlan
type PipelineTask GenericPipelineTask[map[string]interface{}]

// PipelineStage consist of multiple PipelineTasks, they will be executed in parallel
type PipelineStage []*PipelineTask

// PipelinePlan consist of multiple PipelineStages, they will be executed in sequential order
type PipelinePlan []PipelineStage

// IsEmpty checks if a PipelinePlan is empty
func (plan PipelinePlan) IsEmpty() bool {
	if len(plan) == 0 {
		return true
	}
	for _, stage := range plan {
		if len(stage) > 0 {
			return false
		}
	}
	return true
}

type Pipeline struct {
	common.Model
	Name          string       `json:"name" gorm:"index"`
	BlueprintId   uint64       `json:"blueprintId"`
	Plan          PipelinePlan `json:"plan" gorm:"serializer:encdec"`
	TotalTasks    int          `json:"totalTasks"`
	FinishedTasks int          `json:"finishedTasks"`
	BeganAt       *time.Time   `json:"beganAt"`
	FinishedAt    *time.Time   `json:"finishedAt" gorm:"index"`
	Status        string       `json:"status"`
	Message       string       `json:"message"`
	ErrorName     string       `json:"errorName"`
	SpentSeconds  int          `json:"spentSeconds"`
	Stage         int          `json:"stage"`
	Labels        []string     `json:"labels" gorm:"-"`
	SyncPolicy    `gorm:"embedded"`
}

// We use a 2D array because the request body must be an array of a set of tasks
// to be executed concurrently, while each set is to be executed sequentially.
type NewPipeline struct {
	Name        string       `json:"name"`
	Plan        PipelinePlan `json:"plan" swaggertype:"array,string" example:"please check api /pipelines/<PLUGIN_NAME>/pipeline-plan"`
	Labels      []string     `json:"labels"`
	BlueprintId uint64
	SyncPolicy  `gorm:"embedded"`
}

func (Pipeline) TableName() string {
	return "_devlake_pipelines"
}

type DbPipelineLabel struct {
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	PipelineId uint64    `json:"pipeline_id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"primaryKey;index"`
}

func (DbPipelineLabel) TableName() string {
	return "_devlake_pipeline_labels"
}
