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
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
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
	deploymentClause := []dal.Clause{
		dal.Select(`ct.id as task_id, cpc.commit_sha as new_deploy_commit_sha, 
			ct.finished_date as task_finished_date, cpc.repo_id as repo_id`),
		dal.From(`cicd_tasks ct`),
		dal.Join(`left join cicd_pipeline_commits cpc on ct.pipeline_id = cpc.pipeline_id`),
		dal.Join(`left join project_mapping pm on pm.row_id = ct.cicd_scope_id`),
		dal.Where(`ct.environment = ? and ct.type = ? and ct.result = ? and pm.project_name = ? and pm.table = ?`,
			devops.PRODUCTION, devops.DEPLOYMENT, devops.SUCCESS, data.Options.ProjectName, "cicd_scopes"),
		dal.Orderby(`cpc.repo_id, ct.started_date `),
	}
	deploymentDiffPairs := make([]deploymentPair, 0)
	err := db.All(&deploymentDiffPairs, deploymentClause...)
	if err != nil {
		return err
	}
	// deploymentDiffPairs[i-1].NewDeployCommitSha is deploymentDiffPairs[i].OldDeployCommitSha
	oldDeployCommitSha := ""
	lastRepoId := ""
	for i := 0; i < len(deploymentDiffPairs); i++ {
		// if two deployments belong to different repo, let's skip
		if lastRepoId == deploymentDiffPairs[i].RepoId {
			deploymentDiffPairs[i].OldDeployCommitSha = oldDeployCommitSha
		} else {
			lastRepoId = deploymentDiffPairs[i].RepoId
		}
		oldDeployCommitSha = deploymentDiffPairs[i].NewDeployCommitSha
	}

	// get prs by repo project_name
	clauses := []dal.Clause{
		dal.From(&code.PullRequest{}),
		dal.Join(`left join project_mapping pm on pm.row_id = pull_requests.base_repo_id`),
		dal.Where("pull_requests.merged_date IS NOT NULL and pm.project_name = ? and pm.table = ?", data.Options.ProjectName, "repos"),
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
			firstCommit, err := getFirstCommit(pr.Id, db)
			if err != nil {
				return nil, err
			}
			projectPrMetric := &crossdomain.ProjectPrMetric{}
			projectPrMetric.Id = pr.Id
			projectPrMetric.ProjectName = data.Options.ProjectName
			if err != nil {
				return nil, err
			}
			if firstCommit != nil {
				codingTime := int64(pr.CreatedDate.Sub(firstCommit.AuthoredDate).Seconds())
				if codingTime/60 == 0 && codingTime%60 > 0 {
					codingTime = 1
				} else {
					codingTime = codingTime / 60
				}
				projectPrMetric.CodingTimespan = processNegativeValue(codingTime)
				projectPrMetric.FirstCommitSha = firstCommit.Sha
			}
			firstReview, err := getFirstReview(pr.Id, pr.AuthorId, db)
			if err != nil {
				return nil, err
			}
			if firstReview != nil {
				projectPrMetric.ReviewLag = processNegativeValue(int64(firstReview.CreatedDate.Sub(pr.CreatedDate).Minutes()))
				projectPrMetric.ReviewTimespan = processNegativeValue(int64(pr.MergedDate.Sub(firstReview.CreatedDate).Minutes()))
				projectPrMetric.FirstReviewId = firstReview.ReviewId
			}
			deployment, err := getDeployment(pr.MergeCommitSha, pr.BaseRepoId, deploymentDiffPairs, db)
			if err != nil {
				return nil, err
			}
			if deployment != nil && deployment.TaskFinishedDate != nil {
				timespan := deployment.TaskFinishedDate.Sub(*pr.MergedDate)
				projectPrMetric.DeployTimespan = processNegativeValue(int64(timespan.Minutes()))
				projectPrMetric.DeploymentId = deployment.TaskId
			} else {
				log.Debug("deploy time of pr %v is nil\n", pr.PullRequestKey)
			}
			projectPrMetric.ChangeTimespan = nil
			var result int64
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

func getFirstCommit(prId string, db dal.Dal) (*code.Commit, errors.Error) {
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
	return commit, nil
}

func getFirstReview(prId string, prCreator string, db dal.Dal) (*code.PullRequestComment, errors.Error) {
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
	return review, nil
}

func getDeployment(mergeSha string, repoId string, deploymentPairList []deploymentPair, db dal.Dal) (*deploymentPair, errors.Error) {
	// ignore environment at this point because detecting it by name is obviously not engouh
	// take https://github.com/apache/incubator-devlake/actions/workflows/build.yml for example
	// one can not distingush testing/production by looking at the job name solely.
	commitDiff := &code.CommitsDiff{}
	// find if tuple[merge_sha, new_commit_sha, old_commit_sha] exist in commits_diffs, if yes, return pair.FinishedDate
	for _, pair := range deploymentPairList {
		if repoId != pair.RepoId {
			continue
		}
		err := db.First(commitDiff, dal.Where(`commit_sha = ? and new_commit_sha = ? and old_commit_sha = ?`,
			mergeSha, pair.NewDeployCommitSha, pair.OldDeployCommitSha))
		if err == nil {
			return &pair, nil
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
	RepoId             string
	NewDeployCommitSha string
	OldDeployCommitSha string
	TaskFinishedDate   *time.Time
}
