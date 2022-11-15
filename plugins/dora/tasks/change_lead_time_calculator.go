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
	goerror "errors"
	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"
)

func CalculateChangeLeadTime(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	log := taskCtx.GetLogger()
	data := taskCtx.GetData().(*DoraTaskData)
	// construct a list of tuple[task, oldPipelineCommitSha, newPipelineCommitSha, taskFinishedDate]
	pipelineIdClauses := []dal.Clause{
		dal.Select(`ct.id as task_id, cpc.commit_sha as new_deploy_commit_sha, 
			ct.finished_date as task_finished_date`),
		dal.From(`cicd_tasks ct`),
		dal.Join(`left join cicd_pipeline_commits cpc on ct.pipeline_id = cpc.pipeline_id`),
		dal.Join(`left join project_mapping pm on pm.row_id = ct.cicd_scope_id and pm.table = 'cicd_scopes'`),
		dal.Where(`ct.environment = ? and ct.type = ? and pm.project_name = ?`,
			"PRODUCTION", "DEPLOYMENT", data.Options.ProjectName),
		dal.Orderby(`ct.started_date `),
	}
	deploymentPairList := make([]deploymentPair, 0)
	err := db.All(&deploymentPairList, pipelineIdClauses...)
	if err != nil {
		return err
	}
	// deploymentPairList[i-1].NewDeployCommitSha is deploymentPairList[i].OldDeployCommitSha
	oldDeployCommitSha := ""
	for i := 0; i < len(deploymentPairList); i++ {
		deploymentPairList[i].OldDeployCommitSha = oldDeployCommitSha
		oldDeployCommitSha = deploymentPairList[i].NewDeployCommitSha
	}

	// get repo list by projectName
	repoClauses := []dal.Clause{
		dal.From(`project_mapping pm`),
		dal.Where("pm.project_name = ? and pm.table = ?", data.Options.ProjectName, "repos"),
	}
	repoList := make([]string, 0)
	err = db.Pluck(`row_id`, &repoList, repoClauses...)
	if err != nil {
		return err
	}

	// get prs by repo list
	clauses := []dal.Clause{
		dal.From(&code.PullRequest{}),
		dal.Where("merged_date IS NOT NULL and base_repo_id in ?", repoList),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: DoraApiParams{
				ProjectName: data.Options.ProjectName,
			},
			Table: "pull_requests",
		},
		BatchSize:    100,
		InputRowType: reflect.TypeOf(code.PullRequest{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			pr := inputRow.(*code.PullRequest)
			firstCommitDate, err := getFirstCommitTime(pr.Id, db)
			projectPrMetric := &crossdomain.ProjectPrMetrics{}
			projectPrMetric.Id = pr.Id
			projectPrMetric.ProjectName = data.Options.ProjectName
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
				projectPrMetric.CodingTimespan = processNegativeValue(codingTime)
			}
			firstReviewTime, err := getFirstReviewTime(pr.Id, pr.AuthorId, db)
			if err != nil {
				return nil, err
			}
			if firstReviewTime != nil {
				projectPrMetric.ReviewLag = processNegativeValue(int64(firstReviewTime.Sub(pr.CreatedDate).Minutes()))
				projectPrMetric.ReviewTimespan = processNegativeValue(int64(pr.MergedDate.Sub(*firstReviewTime).Minutes()))
			}
			deploymentFinishedDate, err := getDeploymentFinishTime(pr.MergeCommitSha, deploymentPairList, db)
			if err != nil {
				return nil, err
			}
			if deploymentFinishedDate != nil {
				timespan := deploymentFinishedDate.Sub(*pr.MergedDate)
				projectPrMetric.DeployTimespan = processNegativeValue(int64(timespan.Minutes()))
			} else {
				log.Debug("deploy time of pr %v is nil\n", pr.PullRequestKey)
			}
			projectPrMetric.ChangeTimespan = nil
			result := int64(0)
			if projectPrMetric.CodingTimespan != nil {
				result += *projectPrMetric.CodingTimespan
			}
			if projectPrMetric.ReviewLag != nil {
				result += *projectPrMetric.ReviewLag
			}
			if projectPrMetric.ReviewTimespan != nil {
				result += *projectPrMetric.ReviewTimespan
			}
			if projectPrMetric.DeployTimespan != nil {
				result += *projectPrMetric.DeployTimespan
			}
			if result > 0 {
				projectPrMetric.ChangeTimespan = &result
			}
			return []interface{}{projectPrMetric}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

func getFirstCommitTime(prId string, db dal.Dal) (*time.Time, errors.Error) {
	commit := &code.Commit{}
	commitClauses := []dal.Clause{
		dal.From(&code.Commit{}),
		dal.Join("left join pull_request_commits on commits.sha = pull_request_commits.commit_sha"),
		dal.Where("pull_request_commits.pull_request_id = ?", prId),
		dal.Orderby("commits.authored_date ASC"),
	}
	err := db.First(commit, commitClauses...)
	if goerror.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &commit.AuthoredDate, nil
}

func getFirstReviewTime(prId string, prCreator string, db dal.Dal) (*time.Time, errors.Error) {
	review := &code.PullRequestComment{}
	commentClauses := []dal.Clause{
		dal.From(&code.PullRequestComment{}),
		dal.Where("pull_request_id = ? and account_id != ?", prId, prCreator),
		dal.Orderby("created_date ASC"),
	}
	err := db.First(review, commentClauses...)
	if goerror.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &review.CreatedDate, nil
}

func getDeploymentFinishTime(mergeSha string, deploymentPairList []deploymentPair, db dal.Dal) (*time.Time, errors.Error) {
	// ignore environment at this point because detecting it by name is obviously not engouh
	// take https://github.com/apache/incubator-devlake/actions/workflows/build.yml for example
	// one can not distingush testing/production by looking at the job name solely.
	commitDiff := &code.CommitsDiff{}
	// find if tuple[merge_sha, new_commit_sha, old_commit_sha] exist in commits_diffs, if yes, return pair.FinishedDate
	for _, pair := range deploymentPairList {
		err := db.First(commitDiff, dal.Where(`commit_sha = ? and new_commit_sha = ? and old_commit_sha = ?`,
			mergeSha, pair.NewDeployCommitSha, pair.OldDeployCommitSha))
		if err == nil {
			return pair.TaskFinishedDate, nil
		}
		if goerror.Is(err, gorm.ErrRecordNotFound) {
			continue
		}
		if err != nil {
			return nil, err
		}

	}
	return nil, nil
}

func processNegativeValue(v int64) *int64 {
	if v > 0 {
		return &v
	} else {
		return nil
	}
}

var CalculateChangeLeadTimeMeta = core.SubTaskMeta{
	Name:             "calculateChangeLeadTime",
	EntryPoint:       CalculateChangeLeadTime,
	EnabledByDefault: true,
	Description:      "Calculate change lead time",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD, core.DOMAIN_TYPE_CODE},
}

type deploymentPair struct {
	TaskId             string
	NewDeployCommitSha string
	OldDeployCommitSha string
	TaskFinishedDate   *time.Time
}
