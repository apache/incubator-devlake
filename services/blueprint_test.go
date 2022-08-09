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
	"encoding/json"
	"testing"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/stretchr/testify/assert"
)

func TestMergePipelineTasks(t *testing.T) {
	plan1 := core.PipelinePlan{
		{
			{Plugin: "github"},
			{Plugin: "gitlab"},
		},
		{
			{Plugin: "gitextractor1"},
			{Plugin: "gitextractor2"},
		},
	}

	plan2 := core.PipelinePlan{
		{
			{Plugin: "jira"},
		},
	}

	plan3 := core.PipelinePlan{
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

	assert.Equal(t, plan1, MergePipelinePlans(plan1))
	assert.Equal(t, plan2, MergePipelinePlans(plan2))
	assert.Equal(
		t,
		core.PipelinePlan{
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
		MergePipelinePlans(plan1, plan2),
	)
	assert.Equal(
		t,
		core.PipelinePlan{
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
		MergePipelinePlans(plan1, plan2, plan3),
	)
}

func TestFormatPipelinePlans(t *testing.T) {
	beforePlan2 := json.RawMessage(`[[{"plugin":"github"},{"plugin":"gitlab"}],[{"plugin":"gitextractor1"},{"plugin":"gitextractor2"}]]`)

	mainPlan := core.PipelinePlan{
		{
			{Plugin: "jira"},
		},
	}

	afterPlan2 := json.RawMessage(`[[{"plugin":"jenkins"}],[{"plugin":"jenkins"}]]`)

	result1, err1 := FormatPipelinePlans(nil, mainPlan, nil)
	assert.Nil(t, err1)
	assert.Equal(t, mainPlan, result1)

	result2, err2 := FormatPipelinePlans(beforePlan2, mainPlan, afterPlan2)
	assert.Nil(t, err2)
	assert.Equal(t, core.PipelinePlan{
		{
			{Plugin: "github"},
			{Plugin: "gitlab"},
		},
		{
			{Plugin: "gitextractor1"},
			{Plugin: "gitextractor2"},
		},
		{
			{Plugin: "jira"},
		},
		{
			{Plugin: "jenkins"},
		},
		{
			{Plugin: "jenkins"},
		},
	}, result2)

	result3, err3 := FormatPipelinePlans(json.RawMessage("[]"), mainPlan, json.RawMessage("[]"))
	assert.Nil(t, err3)
	assert.Equal(t, mainPlan, result3)
}
