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

package services

import (
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/server/services"
	"github.com/apache/incubator-devlake/test/helper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputePipelineStatus(t *testing.T) {
	client := helper.StartDevLakeServer(t, nil)
	db := client.GetDal()
	// insert fake tasks to database
	pipeline := &models.Pipeline{
		TotalTasks: 3,
	}
	err := db.Create(pipeline)
	assert.Nil(t, err)
	assert.NotZero(t, pipeline.ID)

	task_row1_col1 := &models.Task{
		PipelineId:  pipeline.ID,
		PipelineRow: 1,
		PipelineCol: 1,
		Plugin:      "github",
		Status:      models.TASK_COMPLETED,
	}
	err = db.Create(task_row1_col1)
	assert.Nil(t, err)
	assert.NotZero(t, task_row1_col1.ID)

	task_row1_col2 := &models.Task{
		PipelineId:  pipeline.ID,
		PipelineRow: 1,
		PipelineCol: 2,
		Plugin:      "gitext",
		Status:      models.TASK_FAILED,
	}
	err = db.Create(task_row1_col2)
	assert.Nil(t, err)
	assert.NotZero(t, task_row1_col2.ID)

	task_row2_col1 := &models.Task{
		PipelineId:  pipeline.ID,
		PipelineRow: 2,
		PipelineCol: 1,
		Plugin:      "refdiff",
		Status:      models.TASK_CREATED,
	}
	err = db.Create(task_row2_col1)
	assert.Nil(t, err)
	assert.NotZero(t, task_row2_col1.ID)

	// pipeline.status == "failed" if SkipOnFailed=false and any tasks failed
	status, err := services.ComputePipelineStatus(pipeline, false)
	if !assert.Nil(t, err) {
		println(err.Messages().Format())
	}
	assert.Equal(t, models.TASK_FAILED, status)

	// pipeline.status == "completed" if all latest tasks were succeeded
	task_row1_col2.Status = models.TASK_COMPLETED
	err = db.Update(task_row1_col2)
	assert.Nil(t, err)
	task_row2_col1.Status = models.TASK_COMPLETED
	err = db.Update(task_row2_col1)
	assert.Nil(t, err)
	status, err = services.ComputePipelineStatus(pipeline, false)
	if !assert.Nil(t, err) {
		println(err.Messages().Format())
	}
	assert.Equal(t, models.TASK_COMPLETED, status)

	pipeline.SkipOnFail = true
	err = db.Update(pipeline)
	assert.Nil(t, err)
	status, err = services.ComputePipelineStatus(pipeline, false)
	assert.Nil(t, err)
	assert.Equal(t, models.TASK_COMPLETED, status)

	// pipeline.status == "partial" if SkipOnFail=true and some were succeeded while others not
	task_row1_col1.Status = models.TASK_FAILED
	err = db.Update(task_row1_col1)
	assert.Nil(t, err)
	status, err = services.ComputePipelineStatus(pipeline, false)
	assert.Nil(t, err)
	assert.Equal(t, models.TASK_PARTIAL, status)

	// pipeline.status == "failed" is SkipOnFail=true and all tasks were fail
	task_row1_col1.Status = models.TASK_FAILED
	err = db.Update(task_row1_col1)
	assert.Nil(t, err)
	task_row1_col2.Status = models.TASK_FAILED
	err = db.Update(task_row1_col2)
	assert.Nil(t, err)
	task_row2_col1.Status = models.TASK_FAILED
	err = db.Update(task_row2_col1)
	assert.Nil(t, err)
	status, err = services.ComputePipelineStatus(pipeline, false)
	assert.Nil(t, err)
	assert.Equal(t, models.TASK_FAILED, status)

	// pipeline.status == "completed" if all failed tasks were reran successfully
	task_row1_col1_rerun := &models.Task{
		PipelineId:  pipeline.ID,
		PipelineRow: 1,
		PipelineCol: 1,
		Plugin:      "github",
		Status:      models.TASK_COMPLETED,
	}
	err = db.Create(task_row1_col1_rerun)
	assert.Nil(t, err)
	assert.NotZero(t, task_row1_col1_rerun.ID)

	task_row1_col2_rerun := &models.Task{
		PipelineId:  pipeline.ID,
		PipelineRow: 1,
		PipelineCol: 2,
		Plugin:      "gitext",
		Status:      models.TASK_COMPLETED,
	}
	err = db.Create(task_row1_col2_rerun)
	assert.Nil(t, err)
	assert.NotZero(t, task_row1_col2_rerun.ID)

	task_row2_col1_rerun := &models.Task{
		PipelineId:  pipeline.ID,
		PipelineRow: 2,
		PipelineCol: 1,
		Plugin:      "refdiff",
		Status:      models.TASK_COMPLETED,
	}
	err = db.Create(task_row2_col1_rerun)
	assert.Nil(t, err)
	assert.NotZero(t, task_row2_col1.ID)

	status, err = services.ComputePipelineStatus(pipeline, false)
	assert.Nil(t, err)
	assert.Equal(t, models.TASK_COMPLETED, status)

	// pipeline.status == "partial" if there were failed task in reran tasks
	task_row1_col1_rerun.Status = models.TASK_CANCELLED
	err = db.Update(task_row1_col1_rerun)
	assert.Nil(t, err)
	status, err = services.ComputePipelineStatus(pipeline, false)
	assert.Nil(t, err)
	assert.Equal(t, models.TASK_PARTIAL, status)
}
