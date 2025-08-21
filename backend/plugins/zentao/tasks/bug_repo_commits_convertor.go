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
	"github.com/spf13/cast"
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

var _ plugin.SubTaskEntryPoint = ConvertBugRepoCommits

var ConvertBugRepoCommitsMeta = plugin.SubTaskMeta{
	Name:             "convertBugRepoCommits",
	EntryPoint:       ConvertBugRepoCommits,
	EnabledByDefault: false,
	Description:      "convert Zentao bug repo commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func ConvertBugRepoCommits(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	db := taskCtx.GetDal()

	cursor, err := db.Cursor(
		dal.From(&models.ZentaoBugRepoCommit{}),
		dal.Where(`project = ? and connection_id = ?`, data.Options.ProjectId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGenerator := didgen.NewDomainIdGenerator(&models.ZentaoBug{})
	convertor, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.ZentaoBugRepoCommit{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_BUG_REPO_COMMITS_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			toolEntity := inputRow.(*models.ZentaoBugRepoCommit)
			domainEntity := &crossdomain.IssueRepoCommit{
				IssueId:   issueIdGenerator.Generate(data.Options.ConnectionId, cast.ToInt64(toolEntity.IssueId)),
				RepoUrl:   toolEntity.RepoUrl,
				CommitSha: toolEntity.CommitSha,
			}
			host, namespace, repoName, err := parseRepoUrl(domainEntity.RepoUrl)
			if err != nil {
				return nil, errors.Default.WrapRaw(err)
			}
			domainEntity.Host = host
			domainEntity.Namespace = namespace
			domainEntity.RepoName = repoName

			var results []interface{}
			results = append(results, domainEntity)
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return convertor.Execute()
}
