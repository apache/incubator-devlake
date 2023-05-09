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
	"reflect"
	"regexp"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

var ConvertIssueRepoCommitsMeta = plugin.SubTaskMeta{
	Name:             "convertIssueRepoCommits",
	EntryPoint:       ConvertIssueRepoCommits,
	EnabledByDefault: false,
	Description:      "convert Jira issue repo commits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

// ConvertIssueRepoCommits is to extract issue_repo_commits from jira_issue_commits, nothing difference with
// issue_commits but added a RepoUrl. This task is needed by EE group.
func ConvertIssueRepoCommits(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDal()
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("convert issue repo commits")
	var commitRepoUrlRegex *regexp.Regexp
	commitRepoUrlPattern := `(.*)\-\/commit`
	commitRepoUrlRegex, err := errors.Convert01(regexp.Compile(commitRepoUrlPattern))
	if err != nil {
		return errors.Default.Wrap(err, "regexp Compile commitRepoUrlPattern failed")
	}
	var commitRepoUrlRegexps []*regexp.Regexp
	if tr := data.Options.TransformationRules; tr != nil {
		for _, s := range tr.RemotelinkRepoPattern {
			pattern, e := regexp.Compile(s)
			if e != nil {
				return errors.Convert(e)
			}
			commitRepoUrlRegexps = append(commitRepoUrlRegexps, pattern)
		}
	}

	clauses := []dal.Clause{
		dal.Select("jic.*"),
		dal.From("_tool_jira_issue_commits jic"),
		dal.Join(`left join _tool_jira_board_issues jbi on (
			jbi.connection_id = jic.connection_id
			AND jbi.issue_id = jic.issue_id
		)`),
		dal.Where("jbi.connection_id = ? AND jbi.board_id = ?", connectionId, boardId),
		dal.Orderby("jbi.connection_id, jbi.issue_id"),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGenerator := didgen.NewDomainIdGenerator(&models.JiraIssue{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_REMOTELINK_TABLE,
		},
		InputRowType: reflect.TypeOf(models.JiraIssueCommit{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			issueCommit := inputRow.(*models.JiraIssueCommit)
			item := &crossdomain.IssueRepoCommit{
				IssueId:   issueIdGenerator.Generate(connectionId, issueCommit.IssueId),
				CommitSha: issueCommit.CommitSha,
			}
			if commitRepoUrlRegex != nil {
				groups := commitRepoUrlRegex.FindStringSubmatch(issueCommit.CommitUrl)
				if len(groups) > 1 {
					item.RepoUrl = groups[1]
				}
			}
			api.RefineIssueRepoCommit(item, commitRepoUrlRegexps, issueCommit.CommitUrl)
			return []interface{}{item}, nil
		},
	})
	if err != nil {
		return errors.Convert(err)
	}

	return converter.Execute()
}
