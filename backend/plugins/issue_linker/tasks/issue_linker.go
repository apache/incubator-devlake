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
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

var LinkIssuesMeta = plugin.SubTaskMeta{
	Name:             "Link Pull Requests with Issues",
	EntryPoint:       LinkIssues,
	EnabledByDefault: true,
	Description:      "", // TODO
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
	DependencyTables: []string{
		code.PullRequest{}.TableName(), // cursor
		ticket.Issue{}.TableName()},
	ProductTables: []string{
		crossdomain.PullRequestIssue{}.TableName(),
	},
}

type Config struct {
	IssueRegex string `mapstructure:"IssueRegex" json:"issueRegex"`
}
type IssueLinkerOptions struct {
	ScopeConfig *Config
}
type IssueLinkerTaskData struct {
	Options *IssueLinkerOptions
}

func LinkIssues(taskCtx plugin.SubTaskContext) (err errors.Error) {
	db := taskCtx.GetDal()

	data := taskCtx.GetData().(*IssueLinkerTaskData)

	rawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx:    taskCtx,
		Params: data,
	}

	issuePattern := data.Options.ScopeConfig.IssueRegex
	mrIssueRegex, err := errors.Convert01(regexp.Compile(issuePattern))
	if err != nil {
		return errors.Default.Wrap(err, "regexp compile failed")
	}

	cursor, err := db.Cursor(dal.From(&code.PullRequest{}))
	if err != nil {
		return err
	}
	defer cursor.Close()

	// iterate all rows
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(code.PullRequest{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			pullRequest := inputRow.(*code.PullRequest)

			//find the issue in the body
			issueNumberStr := ""

			if mrIssueRegex != nil {
				issueNumberStr = mrIssueRegex.FindString(pullRequest.Description)
			}
			//find the issue in the title
			if issueNumberStr == "" {
				issueNumberStr = mrIssueRegex.FindString(pullRequest.Title)
			}

			if issueNumberStr == "" {
				return nil, nil
			}

			issueNumberStr = strings.ReplaceAll(issueNumberStr, "#", "")
			issueNumberStr = strings.TrimSpace(issueNumberStr)

			issue := &ticket.Issue{}

			//change the issueNumberStr to int, if cannot be changed, just continue
			issueNumber, numFormatErr := strconv.Atoi(issueNumberStr)
			if numFormatErr != nil {
				return nil, nil
			}
			err = db.All(
				issue,
				dal.Where("issue_key = ?",
					issueNumber),
				dal.Limit(1),
			)
			if err != nil {
				return nil, err
			}
			fmt.Println("found one:", issueNumberStr, issue)
			pullRequestIssue := &crossdomain.PullRequestIssue{
				PullRequestId:  pullRequest.Id,
				IssueId:        issue.Id,
				PullRequestKey: pullRequest.PullRequestKey,
				IssueKey:       issueNumber,
			}

			return []interface{}{pullRequestIssue}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
