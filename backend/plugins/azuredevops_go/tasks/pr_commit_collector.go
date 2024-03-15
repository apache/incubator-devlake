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

package tasks

import (
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"reflect"

	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"

	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
)

func init() {
	RegisterSubtaskMeta(&CollectApiPullRequestCommitsMeta)
}

const RawPrCommitTable = "azuredevops_go_api_pull_request_commits"

var CollectApiPullRequestCommitsMeta = plugin.SubTaskMeta{
	Name:             "collectApiPullRequestCommits",
	EntryPoint:       CollectApiPullRequestCommits,
	EnabledByDefault: true,
	Description:      "Collect PullRequestCommits data from Azure DevOps API.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS, plugin.DOMAIN_TYPE_CODE_REVIEW},
	DependencyTables: []string{models.AzuredevopsPullRequest{}.TableName()},
	ProductTables:    []string{RawPrCommitTable},
}

type SimplePr struct {
	AzuredevopsId int
}

func CollectApiPullRequestCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawPrCommitTable)
	db := taskCtx.GetDal()

	cursor, err := db.Cursor(
		dal.Select("azuredevops_id"),
		dal.From(models.AzuredevopsPullRequest{}.TableName()),
		dal.Where("repository_id = ? and connection_id=?", data.Options.RepositoryId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimplePr{}))
	if err != nil {
		return err
	}

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs:    *rawDataSubTaskArgs,
		ApiClient:             data.ApiClient,
		PageSize:              100,
		Input:                 iterator,
		Incremental:           false,
		UrlTemplate:           "{{ .Params.OrganizationId }}/{{ .Params.ProjectId }}/_apis/git/repositories/{{ .Params.RepositoryId }}/pullRequests/{{ .Input.AzuredevopsId }}/commits?api-version=7.1",
		Query:                 BuildPaginator(true),
		ResponseParser:        ParseRawMessageFromValue,
		GetNextPageCustomData: ExtractContToken,
		AfterResponse:         change203To401,
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}
