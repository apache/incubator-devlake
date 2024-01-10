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
	"math"
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

// CalculateChangeLeadTimeMeta contains metadata for the CalculateChangeLeadTime subtask.
var CalculateChangeLeadTimeMeta = plugin.SubTaskMeta{
	Name:             "calculateChangeLeadTime",
	EntryPoint:       CalculateChangeLeadTime,
	EnabledByDefault: true,
	Description:      "Calculate change lead time",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD, plugin.DOMAIN_TYPE_CODE},
}

// CalculateChangeLeadTime calculates change lead time for a project.
func CalculateChangeLeadTime(taskCtx plugin.SubTaskContext) errors.Error {
	// Get instances of the DAL and logger
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	data := taskCtx.GetData().(*DoraTaskData)

	// Get pull requests by repo project_name
	cursor, err := db.Cursor(
		dal.Select("pr.id, pr.pull_request_key, pr.author_id, pr.merge_commit_sha, pr.created_date, pr.merged_date"),
		dal.From("pull_requests pr"),
		dal.Join(`LEFT JOIN project_mapping pm ON (pm.row_id = pr.base_repo_id)`),
		dal.Where("pr.merged_date IS NOT NULL AND pm.project_name = ? AND pm.table = 'repos'", data.Options.ProjectName),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			// table and params are essential for deleting data from the target table
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
			// Initialize a new ProjectPrMetric
			projectPrMetric := &crossdomain.ProjectPrMetric{}
			projectPrMetric.Id = pr.Id
			projectPrMetric.ProjectName = data.Options.ProjectName

			// Get the first commit for the PR
			firstCommit, err := getFirstCommit(pr.Id, db)
			if err != nil {
				return nil, err
			}
			// Calculate PR coding time
			if firstCommit != nil {
				projectPrMetric.PrCodingTime = computeTimeSpan(&firstCommit.CommitAuthoredDate, &pr.CreatedDate)
				projectPrMetric.FirstCommitSha = firstCommit.CommitSha
			}

			// Get the first review for the PR
			firstReview, err := getFirstReview(pr.Id, pr.AuthorId, db)
			if err != nil {
				return nil, err
			}
			// Calculate PR pickup time and PR review time
			prDuring := computeTimeSpan(&pr.CreatedDate, pr.MergedDate)
			if firstReview != nil {
				projectPrMetric.PrPickupTime = computeTimeSpan(&pr.CreatedDate, &firstReview.CreatedDate)
				projectPrMetric.PrReviewTime = computeTimeSpan(&firstReview.CreatedDate, pr.MergedDate)
				projectPrMetric.FirstReviewId = firstReview.Id
			}

			// Get the deployment for the PR
			deployment, err := getDeploymentCommit(pr.MergeCommitSha, data.Options.ProjectName, db)
			if err != nil {
				return nil, err
			}

			// Calculate PR deploy time
			if deployment != nil && deployment.FinishedDate != nil {
				projectPrMetric.PrDeployTime = computeTimeSpan(pr.MergedDate, deployment.FinishedDate)
				projectPrMetric.DeploymentCommitId = deployment.Id
			} else {
				logger.Debug("deploy time of pr %v is nil\n", pr.PullRequestKey)
			}

			// Calculate PR cycle time
			if projectPrMetric.PrDeployTime != nil {
				var cycleTime int64
				if projectPrMetric.PrCodingTime != nil {
					cycleTime += *projectPrMetric.PrCodingTime
				}
				if prDuring != nil {
					cycleTime += *prDuring
				}
				cycleTime += *projectPrMetric.PrDeployTime
				projectPrMetric.PrCycleTime = &cycleTime
			}
			// Return the projectPrMetric
			return []interface{}{projectPrMetric}, nil
		},
	})
	if err != nil {
		return err
	}
	// Execute the data converter
	return converter.Execute()
}

// getFirstCommit takes a PR ID and a database connection as input, and returns the first commit of the PR.
func getFirstCommit(prId string, db dal.Dal) (*code.PullRequestCommit, errors.Error) {
	// Initialize a pull_request_commit object
	commit := &code.PullRequestCommit{}
	// Define the SQL clauses for the database query
	commitClauses := []dal.Clause{
		dal.From(&code.PullRequestCommit{}),                          // Select from the "pull_request_commits" table
		dal.Where("pull_request_commits.pull_request_id = ?", prId),  // Filter by the PR ID
		dal.Orderby("pull_request_commits.commit_authored_date ASC"), // Order by the authored date of the commits (ascending)
	}

	// Execute the query and retrieve the first commit
	err := db.First(commit, commitClauses...)

	// If any other error occurred, return nil and the error
	if err != nil {
		// If the error indicates that no commit was found, return nil and no error
		if db.IsErrorNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	// If there were no errors, return the first commit and no error
	return commit, nil
}

// getFirstReview takes a PR ID, PR creator ID, and a database connection as input, and returns the first review comment of the PR.
func getFirstReview(prId string, prCreator string, db dal.Dal) (*code.PullRequestComment, errors.Error) {
	// Initialize a review comment object
	review := &code.PullRequestComment{}
	// Define the SQL clauses for the database query
	commentClauses := []dal.Clause{
		dal.From(&code.PullRequestComment{}),                                  // Select from the "pull_request_comments" table
		dal.Where("pull_request_id = ? and account_id != ?", prId, prCreator), // Filter by the PR ID and exclude comments from the PR creator
		dal.Orderby("created_date ASC"),                                       // Order by the created date of the review comments (ascending)
	}

	// Execute the query and retrieve the first review comment
	err := db.First(review, commentClauses...)

	// If any other error occurred, return nil and the error
	if err != nil {
		// If the error indicates that no review comment was found, return nil and no error
		if db.IsErrorNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	// If there were no errors, return the first review comment and no error
	return review, nil
}

// getDeploymentCommit takes a merge commit SHA, a repository ID, a list of deployment pairs, and a database connection as input.
// It returns the deployment pair related to the merge commit, or nil if not found.
func getDeploymentCommit(mergeSha string, projectName string, db dal.Dal) (*devops.CicdDeploymentCommit, errors.Error) {
	deploymentCommits := make([]*devops.CicdDeploymentCommit, 0, 1)
	// do not use `.First` method since gorm would append ORDER BY ID to the query which leads to a error
	err := db.All(
		&deploymentCommits,
		dal.Select("dc.*"),
		dal.From("cicd_deployment_commits dc"),
		dal.Join("LEFT JOIN cicd_deployment_commits p ON (dc.prev_success_deployment_commit_id = p.id)"),
		dal.Join("LEFT JOIN project_mapping pm ON (pm.table = 'cicd_scopes' AND pm.row_id = dc.cicd_scope_id)"),
		dal.Join("INNER JOIN commits_diffs cd ON (cd.new_commit_sha = dc.commit_sha AND cd.old_commit_sha = COALESCE (p.commit_sha, ''))"),
		dal.Where("dc.environment = 'PRODUCTION'"), // TODO: remove this when multi-environment is supported
		dal.Where("pm.project_name = ? AND cd.commit_sha = ? AND dc.RESULT = ?", projectName, mergeSha, devops.RESULT_SUCCESS),
		dal.Orderby("dc.started_date, dc.id ASC"),
		dal.Limit(1),
	)
	if err != nil {
		return nil, err
	}
	if len(deploymentCommits) == 0 {
		return nil, nil
	}
	return deploymentCommits[0], nil
}

func computeTimeSpan(start, end *time.Time) *int64 {
	if start == nil || end == nil {
		return nil
	}
	span := end.Sub(*start)
	minutes := int64(math.Ceil(span.Minutes()))
	if minutes < 0 {
		return nil
	}
	return &minutes
}
