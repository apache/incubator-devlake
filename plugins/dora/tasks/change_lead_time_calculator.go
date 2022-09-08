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
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"
	"reflect"
	"time"
)

func CalculateChangeLeadTime(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*DoraTaskData)
	db := taskCtx.GetDal()
	repoId := data.Options.RepoId
	clauses := []dal.Clause{
		dal.From(&code.PullRequest{}),
		dal.Where("merged_date IS NOT NULL and head_repo_id = ?", repoId),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	enricher, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:    taskCtx,
			Params: DoraApiParams{
				// TODO
			},
			Table: "pull_requests",
		},
		BatchSize:    100,
		InputRowType: reflect.TypeOf(code.PullRequest{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			pr := inputRow.(*code.PullRequest)
			firstCommitDate, err := getFirstCommitTime(pr.Id, db)
			if err != nil {
				return nil, err
			}
			if firstCommitDate != nil {
				codingTime := int64(pr.CreatedDate.Sub(*firstCommitDate).Seconds())
				if codingTime/60 == 0 && codingTime%60 > 0 {
					codingTime = 1
				} else {
					codingTime = codingTime / 60
				}
				pr.OrigCodingTimespan = codingTime
			}
			firstReviewTime, err := getFirstReviewTime(pr.Id, pr.AuthorId, db)
			if err != nil {
				return nil, err
			}
			if firstReviewTime != nil {
				pr.OrigReviewLag = int64(firstReviewTime.Sub(pr.CreatedDate).Minutes())
				pr.OrigReviewTimespan = int64(pr.MergedDate.Sub(*firstReviewTime).Minutes())
			}
			deployTime, err := getDeployTime(repoId, pr.MergeCommitSha, *pr.MergedDate, db)
			if err != nil {
				return nil, err
			}
			if deployTime != nil {
				pr.OrigDeployTimespan = int64(deployTime.Sub(*pr.MergedDate).Minutes())
			}
			processNegativeValue(pr)
			pr.ChangeTimespan = nil
			result := int64(0)
			if pr.CodingTimespan != nil {
				result += *pr.CodingTimespan
			}
			if pr.ReviewLag != nil {
				result += *pr.ReviewLag
			}
			if pr.ReviewTimespan != nil {
				result += *pr.ReviewTimespan
			}
			if pr.DeployTimespan != nil {
				result += *pr.DeployTimespan
			}
			if result > 0 {
				pr.ChangeTimespan = &result
			}
			return []interface{}{pr}, nil
		},
	})
	if err != nil {
		return err
	}

	return enricher.Execute()
}

func getFirstCommitTime(prId string, db dal.Dal) (*time.Time, error) {
	commit := &code.Commit{}
	commitClauses := []dal.Clause{
		dal.From(&code.Commit{}),
		dal.Join("left join pull_request_commits on commits.sha = pull_request_commits.commit_sha"),
		dal.Where("pull_request_commits.pull_request_id = ?", prId),
		dal.Orderby("commits.authored_date ASC"),
	}
	err := db.First(commit, commitClauses...)
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &commit.AuthoredDate, nil
}

func getFirstReviewTime(prId string, prCreator string, db dal.Dal) (*time.Time, error) {
	review := &code.PullRequestComment{}
	commentClauses := []dal.Clause{
		dal.From(&code.PullRequestComment{}),
		dal.Where("pull_request_id = ? and account_id != ?", prId, prCreator),
		dal.Orderby("created_date ASC"),
	}
	err := db.First(review, commentClauses...)
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &review.CreatedDate, nil
}

func getDeployTime(repoId string, mergeSha string, mergeDate time.Time, db dal.Dal) (*time.Time, error) {
	cicdTask := &devops.CICDTask{}
	cicdTaskClauses := []dal.Clause{
		dal.From(&devops.CICDTask{}),
		dal.Join("left join cicd_pipelines on cicd_pipelines.id = cicd_tasks.pipeline_id"),
		dal.Join("left join cicd_pipeline_repos on cicd_pipelines.id = cicd_pipeline_repos.id"),
		dal.Where(`cicd_pipeline_repos.commit_sha = ? 
			and cicd_pipeline_repos.repo = ? 
			and cicd_tasks.type = ? 
			and cicd_tasks.result = ?
			and cicd_tasks.started_date > ?`,
			mergeSha, repoId, "DEPLOY", "SUCCESS", mergeDate),
		dal.Orderby("cicd_tasks.started_date ASC"),
	}
	err := db.First(cicdTask, cicdTaskClauses...)
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return cicdTask.FinishedDate, nil
}

func processNegativeValue(pr *code.PullRequest) {
	if pr.OrigCodingTimespan > 0 {
		pr.CodingTimespan = &pr.OrigCodingTimespan
	} else {
		pr.CodingTimespan = nil
	}
	if pr.OrigReviewLag > 0 {
		pr.ReviewLag = &pr.OrigReviewLag
	} else {
		pr.ReviewLag = nil
	}
	if pr.OrigReviewTimespan > 0 {
		pr.ReviewTimespan = &pr.OrigReviewTimespan
	} else {
		pr.ReviewTimespan = nil
	}
	if pr.OrigDeployTimespan > 0 {
		pr.DeployTimespan = &pr.OrigDeployTimespan
	} else {
		pr.DeployTimespan = nil
	}
}

var CalculateChangeLeadTimeMeta = core.SubTaskMeta{
	Name:             "calculateChangeLeadTime",
	EntryPoint:       CalculateChangeLeadTime,
	EnabledByDefault: true,
	Description:      "Calculate change lead time",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD, core.DOMAIN_TYPE_CODE},
}
