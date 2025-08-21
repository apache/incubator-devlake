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
	"net/url"
	"reflect"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

func ConvertTaskCommit(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_COMMIT_TABLE)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("convert workspace: %d", data.Options.WorkspaceId)

	clauses := []dal.Clause{
		dal.From(&models.TapdTaskCommit{}),
		dal.Where("connection_id = ? AND workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	issueIdGen := didgen.NewDomainIdGenerator(&models.TapdTask{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.TapdTaskCommit{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			toolL := inputRow.(*models.TapdTaskCommit)
			results := make([]interface{}, 0, 2)
			issueCommit := &crossdomain.IssueCommit{
				IssueId:   issueIdGen.Generate(data.Options.ConnectionId, toolL.TaskId),
				CommitSha: toolL.CommitId,
			}
			results = append(results, issueCommit)
			if toolL.WebURL != `` {
				u, err := errors.Convert01(url.Parse(toolL.WebURL))
				if err != nil {
					return nil, err
				}
				repoUrl := toolL.WebURL
				if !strings.HasSuffix(repoUrl, `.git`) {
					repoUrl = repoUrl + `.git`
				}
				issueRepoCommit := &crossdomain.IssueRepoCommit{
					IssueId:   issueIdGen.Generate(data.Options.ConnectionId, toolL.TaskId),
					RepoUrl:   repoUrl,
					CommitSha: toolL.CommitId,
					Host:      u.Hostname(),
					Namespace: getRepoNamespaceFromUrlPath(u.Path),
					RepoName:  getRepoNameFromUrlPath(u.Path),
				}
				results = append(results, issueRepoCommit)
			}

			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var ConvertTaskCommitMeta = plugin.SubTaskMeta{
	Name:             "convertTaskCommit",
	EntryPoint:       ConvertTaskCommit,
	EnabledByDefault: true,
	Description:      "convert Tapd TaskCommit",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}
