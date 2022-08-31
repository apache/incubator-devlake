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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"reflect"
	"regexp"

	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

var ConvertIssueRepoCommitsMeta = core.SubTaskMeta{
	Name:             "convertIssueRepoCommits",
	EntryPoint:       ConvertIssueRepoCommits,
	EnabledByDefault: false,
	Description:      "convert Jira issue repo commits",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

// ConvertIssueRepoCommits is to extract issue_repo_commits from jira_issue_commits, nothing difference with
// issue_commits but added a RepoUrl. This task is needed by EE group.
func ConvertIssueRepoCommits(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDal()
	connectionId := data.Options.ConnectionId
	boardId := data.Options.BoardId
	logger := taskCtx.GetLogger()
	logger.Info("convert issue repo commits")
	var commitRepoUrlRegex *regexp.Regexp
	var err error
	commitRepoUrlPattern := `(.*)\-\/commit`
	commitRepoUrlRegex, err = regexp.Compile(commitRepoUrlPattern)
	if err != nil {
		return errors.Default.Wrap(err, "regexp Compile commitRepoUrlPattern failed")
	}

	clause := []dal.Clause{
		dal.From("_tool_jira_issue_commits jic"),
		dal.Join(`left join _tool_jira_board_issues jbi on (
			jbi.connection_id = jic.connection_id
			AND jbi.issue_id = jic.issue_id
		)`),
		dal.Select("jic.*"),
		dal.Where("jbi.connection_id = ? AND jbi.board_id = ?", connectionId, boardId),
		dal.Orderby("jbi.connection_id, jbi.issue_id"),
	}
	cursor, err := db.Cursor(clause...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGenerator := didgen.NewDomainIdGenerator(&models.JiraIssue{})
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: connectionId,
				BoardId:      boardId,
			},
			Table: RAW_REMOTELINK_TABLE,
		},
		InputRowType: reflect.TypeOf(models.JiraIssueCommit{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			var result []interface{}
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
			result = append(result, item)
			return result, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
