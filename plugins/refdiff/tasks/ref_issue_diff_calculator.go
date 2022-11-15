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
	"github.com/apache/incubator-devlake/errors"
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
)

// CaculatePairList Calculate the pair list both from Options.Pairs and TagPattern
func CaculatePairList(taskCtx core.SubTaskContext) (RefPairLists, errors.Error) {
	data := taskCtx.GetData().(*RefdiffTaskData)
	repoId := data.Options.RepoId
	pairs := data.Options.AllPairs

	pairList := make(RefPairLists, 0, len(pairs))

	for _, pair := range pairs {
		pairList = append(pairList, RefPairList{fmt.Sprintf("%s:%s", repoId, pair[2]), fmt.Sprintf("%s:%s", repoId, pair[3])})
	}

	return pairList, nil
}

func CalculateIssuesDiff(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*RefdiffTaskData)
	repoId := data.Options.RepoId
	db := taskCtx.GetDal()
	// use to calculate progress
	pairList, err := CaculatePairList(taskCtx)
	if err != nil {
		return err
	}
	// commits_diffs join refs => ref_commits_diffs, join pull_request_issues => ref_issues_diffs
	cursor, err := db.Cursor(
		dal.From("commits_diffs"),
		dal.Join(
			`left join (  
        select pull_request_id as id, commit_sha from pull_request_commits 
			left join pull_requests p on pull_request_commits.pull_request_id = p.id
			where p.base_repo_id = ?
			 union  
			select id, merge_commit_sha as commit_sha from pull_requests where base_repo_id = ?) _combine_pr 
			on _combine_pr.commit_sha = commits_diffs.commit_sha`, repoId, repoId),
		dal.Join("left join pull_request_issues on pull_request_issues.pull_request_id = _combine_pr.id"),
		dal.Join("left join refs new_refs on new_refs.commit_sha = commits_diffs.new_commit_sha"),
		dal.Join("left join refs old_refs on old_refs.commit_sha = commits_diffs.old_commit_sha"),
		dal.Orderby("new_refs.id ASC"),
		dal.Where("new_refs.repo_id = ? and pull_request_issues.issue_key > 0 and (new_refs.id, old_refs.id) in ?",
			repoId, pairList),
		dal.Select(`commits_diffs.new_commit_sha as new_ref_commit_sha, commits_diffs.old_commit_sha as old_ref_commit_sha, 
			pull_request_issues.issue_id as issue_id, pull_request_issues.issue_key as issue_number, 
			new_refs.id as new_ref_id, old_refs.id as old_ref_id`),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(crossdomain.RefsIssuesDiffs{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: "commits_diffs",
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
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
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
}
