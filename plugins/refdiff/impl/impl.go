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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/refdiff/tasks"
)

// make sure interface is implemented
var _ core.PluginMeta = (*RefDiff)(nil)
var _ core.PluginTask = (*RefDiff)(nil)
var _ core.PluginApi = (*RefDiff)(nil)
var _ core.PluginModel = (*RefDiff)(nil)
var _ core.PluginMetric = (*RefDiff)(nil)

type RefDiff struct{}

// PluginEntry is a variable exported for Framework to search and load
var PluginEntry RefDiff //nolint

func (plugin RefDiff) Description() string {
	return "Calculate commits diff for specified ref pairs based on `commits` and `commit_parents` tables"
}

func (plugin RefDiff) RequiredDataEntities() (data []map[string]interface{}, err errors.Error) {
	return []map[string]interface{}{}, nil
}

func (plugin RefDiff) GetTablesInfo() []core.Tabler {
	return []core.Tabler{}
}

func (plugin RefDiff) IsProjectMetric() bool {
	return false
}

func (plugin RefDiff) RunAfter() ([]string, errors.Error) {
	return []string{}, nil
}

func (plugin RefDiff) Settings() interface{} {
	return nil
}

func (plugin RefDiff) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CalculateCommitsDiffMeta,
		tasks.CalculateIssuesDiffMeta,
		tasks.CalculatePrCherryPickMeta,
		tasks.CalculateProjectDeploymentCommitsDiffMeta,
	}
}

func (plugin RefDiff) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	var op tasks.RefdiffOptions
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

// PkgPath information lost when compiled as plugin(.so)
func (plugin RefDiff) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/refdiff"
}

func (plugin RefDiff) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return nil
}
