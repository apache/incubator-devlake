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

	"github.com/apache/incubator-devlake/models"
	"github.com/stretchr/testify/assert"
)

func TestMergePipelineTasks(t *testing.T) {
	plan1 := models.PipelinePlan{
		[]*models.NewTask{
			{Plugin: "github"},
			{Plugin: "gitlab"},
		},
		[]*models.NewTask{
			{Plugin: "gitextractor1"},
			{Plugin: "gitextractor2"},
		},
	}

	plan2 := models.PipelinePlan{
		[]*models.NewTask{
			{Plugin: "jira"},
		},
	}

	plan3 := models.PipelinePlan{
		[]*models.NewTask{
			{Plugin: "jenkins"},
		},
		[]*models.NewTask{
			{Plugin: "jenkins"},
		},
		[]*models.NewTask{
			{Plugin: "jenkins"},
		},
	}

	assert.Equal(t, plan1, MergePipelinePlans(plan1))
	assert.Equal(t, plan2, MergePipelinePlans(plan2))
	assert.Equal(
		t,
		models.PipelinePlan{
			[]*models.NewTask{
				{Plugin: "github"},
				{Plugin: "gitlab"},
				{Plugin: "jira"},
			},
			[]*models.NewTask{
				{Plugin: "gitextractor1"},
				{Plugin: "gitextractor2"},
			},
		},
		MergePipelinePlans(plan1, plan2),
	)
	assert.Equal(
		t,
		models.PipelinePlan{
			[]*models.NewTask{
				{Plugin: "github"},
				{Plugin: "gitlab"},
				{Plugin: "jira"},
				{Plugin: "jenkins"},
			},
			[]*models.NewTask{
				{Plugin: "gitextractor1"},
				{Plugin: "gitextractor2"},
				{Plugin: "jenkins"},
			},
			[]*models.NewTask{
				{Plugin: "jenkins"},
			},
		},
		MergePipelinePlans(plan1, plan2, plan3),
	)
}
