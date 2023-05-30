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

package plugin

type GenericPipelineTask[T any] struct {
	Plugin   string   `json:"plugin" binding:"required"`
	Subtasks []string `json:"subtasks"`
	Options  T        `json:"options"`
}

type GenericPipelineStage[T any] []*GenericPipelineTask[T]
type GenericPipelinePlan[T any] []GenericPipelineStage[T]

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
