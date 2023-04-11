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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"reflect"
	"time"
)

// CalculateChangeLeadTime calculates change lead time for a project.
func CalculateChangeLeadTime(taskCtx plugin.SubTaskContext) errors.Error {
	// Get instances of the DAL and logger
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	data := taskCtx.GetData().(*DoraTaskData)

	// Build deployment pairs
	deploymentDiffPairs, err := buildDeploymentPairs(db, data)
	if err != nil {
		return err
	}

	// Get pull requests by repo project_name
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

	// Initialize a new data converter
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
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
			// Process each pull request
			pr := inputRow.(*code.PullRequest)

			// Get the first commit for the PR
			firstCommit, err := getFirstCommit(pr.Id, db)
			if err != nil {
				return nil, err
			}

			// Initialize a new ProjectPrMetric
			projectPrMetric := &crossdomain.ProjectPrMetric{}
			projectPrMetric.Id = pr.Id
			projectPrMetric.ProjectName = data.Options.ProjectName

			// Calculate PR coding time
			if firstCommit != nil {
				codingTime := int64(pr.CreatedDate.Sub(firstCommit.AuthoredDate).Seconds())
				if codingTime/60 == 0 && codingTime%60 > 0 {
					codingTime = 1
				} else {
					codingTime = codingTime / 60
				}
				projectPrMetric.PrCodingTime = processNegativeValue(codingTime)
				projectPrMetric.FirstCommitSha = firstCommit.Sha
			}

			// Get the first review for the PR
			firstReview, err := getFirstReview(pr.Id, pr.AuthorId, db)
			if err != nil {
				return nil, err
			}

			// Calculate PR pickup time and PR review time
			prDuring := processNegativeValue(int64(pr.MergedDate.Sub(pr.CreatedDate).Minutes()))
			if firstReview != nil {
				projectPrMetric.PrPickupTime = processNegativeValue(int64(firstReview.CreatedDate.Sub(pr.CreatedDate).Minutes()))
				projectPrMetric.PrReviewTime = processNegativeValue(int64(pr.MergedDate.Sub(firstReview.CreatedDate).Minutes()))
				projectPrMetric.FirstReviewId = firstReview.Id
			}

			// Get the deployment for the PR
			deployment, err := getDeployment(pr.MergeCommitSha, pr.BaseRepoId, deploymentDiffPairs, db)
			if err != nil {
				return nil, err
			}

			// Calculate PR deploy time
			if deployment != nil && deployment.TaskFinishedDate != nil {
				timespan := deployment.TaskFinishedDate.Sub(*pr.MergedDate)
				projectPrMetric.PrDeployTime = processNegativeValue(int64(timespan.Minutes()))
				projectPrMetric.DeploymentId = deployment.TaskId
			} else {
				logger.Debug("deploy time of pr %v is nil\n", pr.PullRequestKey)
			}

			// Calculate PR cycle time
			if deployment == nil || projectPrMetric.PrDeployTime == nil {
				// Return the projectPrMetric with nill cycle time
				return []interface{}{projectPrMetric}, nil
			}
			var result int64
			if projectPrMetric.PrCodingTime != nil {
				result += *projectPrMetric.PrCodingTime
			}
			if prDuring != nil {
				result += *prDuring
			}
			result += *projectPrMetric.PrDeployTime
			if result > 0 {
				projectPrMetric.PrCycleTime = &result
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

// buildDeploymentPairs populates the OldDeployCommitSha field of each deploymentPair in the given slice.
func buildDeploymentPairs(db dal.Dal, data *DoraTaskData) ([]deploymentPair, errors.Error) {
	// Construct a list of tuple[task, oldPipelineCommitSha, newPipelineCommitSha, taskFinishedDate]
	deploymentClause := []dal.Clause{
		dal.Select(`ct.id as task_id, cpc.commit_sha as new_deploy_commit_sha,
			ct.finished_date as task_finished_date, cpc.repo_id as repo_id`),
		dal.From(`cicd_tasks ct`),
		dal.Join(`left join cicd_pipeline_commits cpc on ct.pipeline_id = cpc.pipeline_id`),
		dal.Join(`left join project_mapping pm on pm.row_id = ct.cicd_scope_id`),
		dal.Where(`ct.environment = ? and ct.type = ? and ct.result = ? and pm.project_name = ? and pm.table = ?`,
			devops.PRODUCTION, devops.DEPLOYMENT, devops.SUCCESS, data.Options.ProjectName, "cicd_scopes"),
		dal.Orderby(`cpc.repo_id, ct.started_date`),
	}

	// Initialize deploymentDiffPairs without oldPipelineCommitSha
	deploymentDiffPairs := make([]deploymentPair, 0)
	err := db.All(&deploymentDiffPairs, deploymentClause...)
	if err != nil {
		return nil, err
	}

	oldDeployCommitSha := ""
	lastRepoId := ""
	for i := 0; i < len(deploymentDiffPairs); i++ {
		// If two deployments belong to different repos, skip
		if lastRepoId == deploymentDiffPairs[i].RepoId {
			deploymentDiffPairs[i].OldDeployCommitSha = oldDeployCommitSha
		} else {
			lastRepoId = deploymentDiffPairs[i].RepoId
		}
		oldDeployCommitSha = deploymentDiffPairs[i].NewDeployCommitSha
	}
	return deploymentDiffPairs, nil
}

// getFirstCommit takes a PR ID and a database connection as input, and returns the first commit of the PR.
func getFirstCommit(prId string, db dal.Dal) (*code.Commit, errors.Error) {
	// Initialize a commit object
	commit := &code.Commit{}
	// Define the SQL clauses for the database query
	commitClauses := []dal.Clause{
		dal.From(&code.Commit{}), // Select from the "commits" table
		dal.Join("left join pull_request_commits on commits.sha = pull_request_commits.commit_sha"), // Join with the "pull_request_commits" table
		dal.Where("pull_request_commits.pull_request_id = ?", prId),                                 // Filter by the PR ID
		dal.Orderby("commits.authored_date ASC"),                                                    // Order by the authored date of the commits (ascending)
	}

	// Execute the query and retrieve the first commit
	err := db.First(commit, commitClauses...)

	// If the error indicates that no commit was found, return nil and no error
	if db.IsErrorNotFound(err) {
		return nil, nil
	}

	// If any other error occurred, return nil and the error
	if err != nil {
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

	// If the error indicates that no review comment was found, return nil and no error
	if db.IsErrorNotFound(err) {
		return nil, nil
	}

	// If any other error occurred, return nil and the error
	if err != nil {
		return nil, err
	}

	// If there were no errors, return the first review comment and no error
	return review, nil
}

// getDeployment takes a merge commit SHA, a repository ID, a list of deployment pairs, and a database connection as input.
// It returns the deployment pair related to the merge commit, or nil if not found.
func getDeployment(mergeSha string, repoId string, deploymentPairList []deploymentPair, db dal.Dal) (*deploymentPair, errors.Error) {
	// Ignore environment detection based on job name since it's not enough to distinguish between testing/production environments.
	commitDiff := &code.CommitsDiff{}
	// Iterate through the deploymentPairList to find if a tuple [merge_sha, new_commit_sha, old_commit_sha] exists in commits_diffs.
	// If found, return the corresponding deployment pair.
	for _, pair := range deploymentPairList {
		// Continue to the next iteration if the repoId does not match the current deployment pair's RepoId.
		if repoId != pair.RepoId {
			continue
		}

		// Query the database to find the commit diff with the given merge SHA, new commit SHA, and old commit SHA.
		err := db.First(commitDiff, dal.Where(`commit_sha = ? and new_commit_sha = ? and old_commit_sha = ?`,
			mergeSha, pair.NewDeployCommitSha, pair.OldDeployCommitSha))

		// If no error occurred, return the current deployment pair.
		if err == nil {
			return &pair, nil
		}

		// If the error indicates that no commit diff was found, continue to the next iteration.
		if db.IsErrorNotFound(err) {
			continue
		}

		// If any other error occurred, return nil and the error.
		if err != nil {
			return nil, err
		}
	}

	// If no matching deployment pair was found, return nil and no error.
	return nil, nil
}

func processNegativeValue(v int64) *int64 {
	if v > 0 {
		return &v
	} else {
		return nil
	}
}

// CalculateChangeLeadTimeMeta contains metadata for the CalculateChangeLeadTime subtask.
var CalculateChangeLeadTimeMeta = plugin.SubTaskMeta{
	Name:             "calculateChangeLeadTime",
	EntryPoint:       CalculateChangeLeadTime,
	EnabledByDefault: true,
	Description:      "Calculate change lead time",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD, plugin.DOMAIN_TYPE_CODE},
}

// deploymentPair is a struct representing a deployment pair with fields for task ID, repo ID, new and old commit SHAs, and task finished date.
type deploymentPair struct {
	TaskId             string
	RepoId             string
	NewDeployCommitSha string
	OldDeployCommitSha string
	TaskFinishedDate   *time.Time
}
