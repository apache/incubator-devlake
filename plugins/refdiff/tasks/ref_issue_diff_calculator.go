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

	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

func CalculateIssuesDiff(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*RefdiffTaskData)
	repoId := data.Options.RepoId
	pairs := data.Options.Pairs
	db := taskCtx.GetDb()
	// use to calculate progress
	pairList := make([][2]string, len(pairs))
	for i, pair := range pairs {
		pairList[i] = [2]string{fmt.Sprintf("%s:%s", repoId, pair.NewRef), fmt.Sprintf("%s:%s", repoId, pair.OldRef)}
	}
	cursor, err := db.Table("refs_commits_diffs").
		Joins(
			`left join (  
        select pull_request_id as id, commit_sha from pull_request_commits 
			left join pull_requests p on pull_request_commits.pull_request_id = p.id
			where p.base_repo_id = ?
			 union  
			select id, merge_commit_sha as commit_sha from pull_requests where base_repo_id = ?) _combine_pr 
			on _combine_pr.commit_sha = refs_commits_diffs.commit_sha`, repoId, repoId).
		Joins("left join pull_request_issues on pull_request_issues.pull_request_id = _combine_pr.id").
		Joins("left join refs on refs.commit_sha = refs_commits_diffs.new_ref_commit_sha").
		Order("refs_commits_diffs.new_ref_id ASC").
		Where("refs.repo_id = ? and pull_request_issues.issue_number > 0 and (refs_commits_diffs.new_ref_id, refs_commits_diffs.old_ref_id) in ?",
			repoId, pairList).
		Select(`refs_commits_diffs.new_ref_commit_sha as new_ref_commit_sha, refs_commits_diffs.old_ref_commit_sha as old_ref_commit_sha, 
			pull_request_issues.issue_id as issue_id, pull_request_issues.issue_number as issue_number, 
			refs_commits_diffs.new_ref_id as new_ref_id, refs_commits_diffs.old_ref_id as old_ref_id`).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(crossdomain.RefsIssuesDiffs{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: "refs_commits_diffs",
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			refPairIssue := inputRow.(*crossdomain.RefsIssuesDiffs)
			return []interface{}{
				refPairIssue,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var CalculateIssuesDiffMeta = core.SubTaskMeta{
	Name:             "calculateIssuesDiff",
	EntryPoint:       CalculateIssuesDiff,
	EnabledByDefault: true,
	Description:      "Calculate diff issues between refs",
}
