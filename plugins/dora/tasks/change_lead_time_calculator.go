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
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"
)

func CalculateChangeLeadTime(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	log := taskCtx.GetLogger()
	clauses := []dal.Clause{
		dal.From(&code.PullRequest{}),
		dal.Where("merged_date IS NOT NULL"),
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
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
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
			deployment, err := getDeployment(devops.PRODUCTION, *pr.MergedDate, db)
			if err != nil {
				return nil, err
			}
			if deployment != nil && deployment.FinishedDate != nil {
				timespan := deployment.FinishedDate.Sub(*pr.MergedDate)
				pr.OrigDeployTimespan = int64(timespan.Minutes())
			} else {
				log.Debug("deploy time of pr %v is nil\n", pr.PullRequestKey)
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

func getDeployment(environment string, mergeDate time.Time, db dal.Dal) (*devops.CICDTask, errors.Error) {
	// ignore environment at this point because detecting it by name is obviously not engouh
	// take https://github.com/apache/incubator-devlake/actions/workflows/build.yml for example
	// one can not distingush testing/production by looking at the job name solely.
	cicdTask := &devops.CICDTask{}
	cicdTaskClauses := []dal.Clause{
		dal.From(&devops.CICDTask{}),
		dal.Where(`
			type = ?
			AND cicd_tasks.result = ?
			AND cicd_tasks.started_date > ?`,
			"DEPLOYMENT",
			"SUCCESS",
			mergeDate,
		),
		dal.Orderby("cicd_tasks.started_date ASC"),
		dal.Limit(1),
	}
	err := db.First(cicdTask, cicdTaskClauses...)
	if goerror.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return cicdTask, nil
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
