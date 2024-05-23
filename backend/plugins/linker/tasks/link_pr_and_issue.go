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
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

var LinkPrToIssueMeta = plugin.SubTaskMeta{
	Name:             "LinkPrToIssue",
	EntryPoint:       LinkPrToIssue,
	EnabledByDefault: true,
	Description:      "Try to link pull requests to issues, according to pull requests' title and description",
	DependencyTables: []string{code.PullRequest{}.TableName(), ticket.Issue{}.TableName()},
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE, plugin.DOMAIN_TYPE_TICKET, plugin.DOMAIN_TYPE_CROSS},
	ProductTables:    []string{crossdomain.PullRequestIssue{}.TableName()},
}

func normalizeIssueKey(issueNumberStr string) string {
	issueNumberStr = strings.ReplaceAll(issueNumberStr, "#", "")
	issueNumberStr = strings.TrimSpace(issueNumberStr)
	return issueNumberStr
}

func LinkPrToIssue(taskCtx plugin.SubTaskContext) errors.Error {

	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*LinkerTaskData)

	rawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx:    taskCtx,
		Params: data,
		//Options:Link
		Table: code.PullRequest{}.TableName(),
	}

	//issuePattern := data.Options.ScopeConfig.IssueRegex
	issuePattern := ""
	issueRegex, err := errors.Convert01(regexp.Compile(issuePattern))
	if err != nil {
		return errors.Default.Wrap(err, "regexp compile failed")
	}

	cursor, err := db.Cursor(dal.From(&code.PullRequest{}))
	if err != nil {
		return err
	}

	defer cursor.Close()

	// iterate all rows
	enricher, err := api.NewDataEnricher(api.DataEnricherArgs[code.PullRequest]{
		Ctx:   taskCtx,
		Name:  code.PullRequest{}.TableName(),
		Input: cursor,
		Enrich: func(pullRequest *code.PullRequest) ([]interface{}, errors.Error) {

			issueNumberStr := ""
			for _, text := range []string{pullRequest.Title, pullRequest.Description} {
				issueNumberStr = issueRegex.FindString(text)
				if issueNumberStr != "" {
					break
				}
			}
			issueNumberStr = normalizeIssueKey(issueNumberStr)
			if issueNumberStr == "" {
				return nil, nil
			}

			issue := &ticket.Issue{}

			db.F

			err = db.All(
				issue,
				dal.Where("issue_key = ?",
					issueNumber),
				dal.Limit(1),
			)
			if err != nil {
				return nil, err
			}

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

	return enricher.Execute()
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(code.PullRequest{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			pullRequest := inputRow.(*code.PullRequest)

			//find the issue in the body
			issueNumberStr := ""

			if issueRegex != nil {
				issueNumberStr = issueRegex.FindString(pullRequest.Description)
			}
			//find the issue in the title
			if issueNumberStr == "" {
				issueNumberStr = issueRegex.FindString(pullRequest.Title)
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
