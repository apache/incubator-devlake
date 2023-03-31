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
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"net/url"
	"reflect"
	"strings"
)

func ConvertBugCommit(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_BUG_COMMIT_TABLE)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	logger.Info("convert workspace: %d", data.Options.WorkspaceId)
	clauses := []dal.Clause{
		dal.From(&models.TapdBugCommit{}),
		dal.Where("connection_id = ? AND workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	issueIdGen := didgen.NewDomainIdGenerator(&models.TapdBug{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.TapdBugCommit{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			toolL := inputRow.(*models.TapdBugCommit)
			results := make([]interface{}, 0, 2)
			issueCommit := &crossdomain.IssueCommit{
				IssueId:   issueIdGen.Generate(data.Options.ConnectionId, toolL.BugId),
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
					IssueId:   issueIdGen.Generate(data.Options.ConnectionId, toolL.BugId),
					RepoUrl:   repoUrl,
					CommitSha: toolL.CommitId,
					Host:      u.Host,
					Namespace: strings.Split(u.Path, `/`)[1],
					RepoName:  toolL.HookProjectName,
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

var ConvertBugCommitMeta = plugin.SubTaskMeta{
	Name:             "convertBugCommit",
	EntryPoint:       ConvertBugCommit,
	EnabledByDefault: true,
	Description:      "convert Tapd BugCommit",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}
