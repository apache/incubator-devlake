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
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/issue_linker/tasks"
)

// make sure interface is implemented
var _ interface {
	plugin.PluginMeta
	plugin.PluginTask
	plugin.PluginApi
	plugin.PluginModel
	plugin.PluginMetric
} = (*IssueLinker)(nil)

type IssueLinker struct{}

func (i IssueLinker) Description() string {
	return "Links `pull_requests` and `issues` via evidence in the pull request metadata."
}

func (i IssueLinker) Name() string {
	return "issue_linker"
}

func (i IssueLinker) RequiredDataEntities() (data []map[string]interface{}, err errors.Error) {
	return []map[string]interface{}{}, nil
}

func (i IssueLinker) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&crossdomain.PullRequestIssue{},
	}
}

func (i IssueLinker) IsProjectMetric() bool {
	return false
}

func (i IssueLinker) RunAfter() ([]string, errors.Error) {
	return []string{}, nil // REVIEW
}

func (i IssueLinker) Settings() interface{} {
	return nil
}

func (i IssueLinker) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.LinkIssuesMeta,
	}
}

func (i IssueLinker) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {

	var op tasks.IssueLinkerOptions
	err := helper.Decode(options, &op, nil)
	if err != nil {
		return nil, err
	}

	return tasks.IssueLinkerTaskData{
		Options: &op,
	}, nil // REVIEW
}

// RootPkgPath information lost when compiled as plugin(.so)
func (i IssueLinker) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/issue_linker"
}

func (i IssueLinker) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return nil
}
