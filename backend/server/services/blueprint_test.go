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
	"testing"

	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/stretchr/testify/assert"
)

func TestParallelizePipelineTasks(t *testing.T) {
	plan1 := coreModels.PipelinePlan{
		{
			{Plugin: "github"},
			{Plugin: "gitlab"},
		},
		{
			{Plugin: "gitextractor1"},
			{Plugin: "gitextractor2"},
		},
	}

	plan2 := coreModels.PipelinePlan{
		{
			{Plugin: "jira"},
		},
	}

	plan3 := coreModels.PipelinePlan{
		{
			{Plugin: "jenkins"},
		},
		{
			{Plugin: "jenkins"},
		},
		{
			{Plugin: "jenkins"},
		},
	}

	assert.Equal(t, plan1, ParallelizePipelinePlans(plan1))
	assert.Equal(t, plan2, ParallelizePipelinePlans(plan2))
	assert.Equal(
		t,
		coreModels.PipelinePlan{
			{
				{Plugin: "github"},
				{Plugin: "gitlab"},
				{Plugin: "jira"},
			},
			{
				{Plugin: "gitextractor1"},
				{Plugin: "gitextractor2"},
			},
		},
		ParallelizePipelinePlans(plan1, plan2),
	)
	assert.Equal(
		t,
		coreModels.PipelinePlan{
			{
				{Plugin: "github"},
				{Plugin: "gitlab"},
				{Plugin: "jira"},
				{Plugin: "jenkins"},
			},
			{
				{Plugin: "gitextractor1"},
				{Plugin: "gitextractor2"},
				{Plugin: "jenkins"},
			},
			{
				{Plugin: "jenkins"},
			},
		},
		ParallelizePipelinePlans(plan1, plan2, plan3),
	)
}

func TestRemoveCollectorTasks(t *testing.T) {
	plan1 := coreModels.PipelinePlan{
		{
			{
				Plugin: "github",
				Subtasks: []string{
					"CollectApiPipelines",
					"ExtractApiPipelines",
					"collectApiPipelineDetails",
					"extractApiPipelineDetails",
					"collectApiJobs",
					"extractApiJobs",
					"collectAccounts",
					"extractAccounts",
					"ConvertAccounts",
					"convertApiProject",
					"convertPipelines",
					"convertPipelineCommits",
					"convertJobs",
				},
			},
		},
		{
			{
				Plugin:   "starrocks",
				Subtasks: []string{},
			},
		},
	}
	assert.Equal(t, coreModels.PipelinePlan{
		{
			{
				Plugin: "github",
				Subtasks: []string{
					"ExtractApiPipelines",
					"extractApiPipelineDetails",
					"extractApiJobs",
					"extractAccounts",
					"ConvertAccounts",
					"convertApiProject",
					"convertPipelines",
					"convertPipelineCommits",
					"convertJobs",
				},
			},
		},
		{
			{
				Plugin:   "starrocks",
				Subtasks: []string{},
			},
		},
	}, removeCollectorTasks(plan1))
}
