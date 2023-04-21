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
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

// CalculateChangeLeadTimeOldMeta will be removed in v0.17
// DEPRECATED
func CalculateChangeLeadTimeOld(taskCtx plugin.SubTaskContext) errors.Error {
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

	enricher, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:    taskCtx,
			Params: DoraApiParams{},
			Table:  "pull_requests",
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
			projectPrMetric.ProjectName = "DEFAULT_PROJECT_NAME"
			if err != nil {
				return nil, err
			}
			if firstCommit != nil {
				projectPrMetric.PrCodingTime = computeTimeSpan(&firstCommit.CommitAuthoredDate, &pr.CreatedDate)
				projectPrMetric.FirstCommitSha = firstCommit.CommitSha
			}
			firstReview, err := getFirstReview(pr.Id, pr.AuthorId, db)
			if err != nil {
				return nil, err
			}
			if firstReview != nil {
				projectPrMetric.PrPickupTime = computeTimeSpan(&pr.CreatedDate, &firstReview.CreatedDate)
				projectPrMetric.PrReviewTime = computeTimeSpan(&firstReview.CreatedDate, pr.MergedDate)
				projectPrMetric.FirstReviewId = firstReview.Id
			}
			deployment, err := getDeploymentOld(devops.PRODUCTION, *pr.MergedDate, db)
			if err != nil {
				return nil, err
			}
			if deployment != nil && deployment.FinishedDate != nil {
				projectPrMetric.PrDeployTime = computeTimeSpan(pr.MergedDate, deployment.FinishedDate)
				projectPrMetric.DeploymentCommitId = deployment.Id
			} else {
				log.Debug("deploy time of pr %v is nil\n", pr.PullRequestKey)
			}
			projectPrMetric.PrCycleTime = nil
			var result int64
			if projectPrMetric.PrCodingTime != nil {
				result += *projectPrMetric.PrCodingTime
			}
			if projectPrMetric.PrPickupTime != nil {
				result += *projectPrMetric.PrPickupTime
			}
			if projectPrMetric.PrReviewTime != nil {
				result += *projectPrMetric.PrReviewTime
			}
			if projectPrMetric.PrDeployTime != nil {
				result += *projectPrMetric.PrDeployTime
			}
			if result > 0 {
				projectPrMetric.PrCycleTime = &result
			}
			return []interface{}{projectPrMetric}, nil
		},
	})
	if err != nil {
		return err
	}

	return enricher.Execute()
}

func getDeploymentOld(environment string, mergeDate time.Time, db dal.Dal) (*devops.CICDTask, errors.Error) {
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
	if db.IsErrorNotFound(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return cicdTask, nil
}

// CalculateChangeLeadTimeOldMeta will be removed in v0.17
// DEPRECATED
var CalculateChangeLeadTimeOldMeta = plugin.SubTaskMeta{
	Name:             "calculateChangeLeadTimeOld",
	EntryPoint:       CalculateChangeLeadTimeOld,
	EnabledByDefault: false,
	Description:      "Calculate change lead time",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD, plugin.DOMAIN_TYPE_CODE},
}
