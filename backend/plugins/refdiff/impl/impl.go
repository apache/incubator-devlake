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

package impl

import (
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/refdiff/models"
	"github.com/apache/incubator-devlake/plugins/refdiff/tasks"
)

// make sure interface is implemented
var _ interface {
	plugin.PluginMeta
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginMetric
} = (*RefDiff)(nil)

type RefDiff struct{}

func (p RefDiff) Description() string {
	return "Calculate commits diff for specified ref pairs based on `commits` and `commit_parents` tables"
}

func (p RefDiff) Name() string {
	return "refdiff"
}

func (p RefDiff) RequiredDataEntities() (data []map[string]interface{}, err errors.Error) {
	return []map[string]interface{}{}, nil
}

func (p RefDiff) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.FinishedCommitsDiff{},
	}
}

func (p RefDiff) IsProjectMetric() bool {
	return false
}

func (p RefDiff) RunAfter() ([]string, errors.Error) {
	return []string{}, nil
}

func (p RefDiff) Settings() interface{} {
	return nil
}

func (p RefDiff) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CalculateCommitsDiffMeta,
		tasks.CalculateIssuesDiffMeta,
		tasks.CalculatePrCherryPickMeta,
		tasks.CalculateDeploymentCommitsDiffMeta,
	}
}

func (p RefDiff) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op models.RefdiffOptions
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}

	db := taskCtx.GetDal()
	tagsPattern := op.TagsPattern
	tagsLimit := op.TagsLimit
	tagsOrder := op.TagsOrder

	rs, err := tasks.CalculateTagPattern(db, tagsPattern, tagsLimit, tagsOrder)
	if err != nil {
		return nil, err
	}
	op.AllPairs, err = tasks.CalculateCommitPairs(db, op.RepoId, op.Pairs, rs)
	if err != nil {
		return nil, err
	}

	return &tasks.RefdiffTaskData{
		Options: &op,
	}, nil
}

// RootPkgPath information lost when compiled as plugin(.so)
func (p RefDiff) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/refdiff"
}

func (p RefDiff) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return nil
}
