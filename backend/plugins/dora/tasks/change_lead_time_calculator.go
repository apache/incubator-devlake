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
	// Clear previous results from the project
	err := db.Exec("DELETE FROM project_pr_metrics WHERE project_name = ? ", data.Options.ProjectName)
	if err != nil {
		return errors.Default.Wrap(err, "error deleting previous project_pr_metrics")
	}

	// Batch fetch all required data upfront for better performance
	startTime := time.Now()
	logger.Info("Batch fetching data for project: %s", data.Options.ProjectName)

	firstCommitsMap, err := batchFetchFirstCommits(data.Options.ProjectName, db)
	if err != nil {
		return errors.Default.Wrap(err, "failed to batch fetch first commits")
	}
	logger.Info("Fetched %d first commits in %v", len(firstCommitsMap), time.Since(startTime))

	reviewStartTime := time.Now()
	firstReviewsMap, err := batchFetchFirstReviews(data.Options.ProjectName, db)
	if err != nil {
		return errors.Default.Wrap(err, "failed to batch fetch first reviews")
	}
	logger.Info("Fetched %d first reviews in %v", len(firstReviewsMap), time.Since(reviewStartTime))

	deploymentStartTime := time.Now()
	deploymentsMap, err := batchFetchDeployments(data.Options.ProjectName, db)
	if err != nil {
		return errors.Default.Wrap(err, "failed to batch fetch deployments")
	}
	logger.Info("Fetched %d deployments in %v", len(deploymentsMap), time.Since(deploymentStartTime))
	logger.Info("Total batch fetch time: %v", time.Since(startTime))

	// Get pull requests by repo project_name
	var clauses = []dal.Clause{
		dal.Select("pr.id, pr.pull_request_key, pr.author_id, pr.merge_commit_sha, pr.created_date, pr.merged_date"),
		dal.From("pull_requests pr"),
		dal.Join(`LEFT JOIN project_mapping pm ON (pm.row_id = pr.base_repo_id)`),
		dal.Where("pr.merged_date IS NOT NULL AND pm.project_name = ? AND pm.table = 'repos'", data.Options.ProjectName),
	}
	cursor, err := db.Cursor(clauses...)
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

			// Get the first commit for the PR from batch-fetched map
			firstCommit := firstCommitsMap[pr.Id]
			// Calculate PR coding time
			if firstCommit != nil {
				projectPrMetric.PrCodingTime = computeTimeSpan(&firstCommit.CommitAuthoredDate, &pr.CreatedDate)
				projectPrMetric.FirstCommitSha = firstCommit.CommitSha
				projectPrMetric.FirstCommitAuthoredDate = &firstCommit.CommitAuthoredDate
			}

			// Get the first review for the PR from batch-fetched map
			firstReview := firstReviewsMap[pr.Id]
			// Calculate PR pickup time and PR review time
			prDuring := computeTimeSpan(&pr.CreatedDate, pr.MergedDate)
			if firstReview != nil {
				projectPrMetric.PrPickupTime = computeTimeSpan(&pr.CreatedDate, &firstReview.CreatedDate)
				projectPrMetric.PrReviewTime = computeTimeSpan(&firstReview.CreatedDate, pr.MergedDate)
				projectPrMetric.FirstReviewId = firstReview.Id
				projectPrMetric.FirstCommentDate = &firstReview.CreatedDate
			}

			projectPrMetric.PrCreatedDate = &pr.CreatedDate
			projectPrMetric.PrMergedDate = pr.MergedDate

			// Get the deployment for the PR from batch-fetched map
			deployment := deploymentsMap[pr.MergeCommitSha]

			// Calculate PR deploy time
			if deployment != nil && deployment.FinishedDate != nil {
				projectPrMetric.PrDeployTime = computeTimeSpan(pr.MergedDate, deployment.FinishedDate)
				projectPrMetric.DeploymentCommitId = deployment.Id
				projectPrMetric.PrDeployedDate = deployment.FinishedDate
			} else {
				logger.Debug("deploy time of pr %v is nil\n", pr.PullRequestKey)
			}

			// Calculate PR cycle time
			var cycleTime int64
			if projectPrMetric.PrCodingTime != nil {
				cycleTime += *projectPrMetric.PrCodingTime
			}
			if prDuring != nil {
				cycleTime += *prDuring
			}
			if projectPrMetric.PrDeployTime != nil {
				cycleTime += *projectPrMetric.PrDeployTime
			}
			projectPrMetric.PrCycleTime = &cycleTime

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
		dal.Where("dc.prev_success_deployment_commit_id <> ''"),
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

// deploymentCommitWithMergeSha is a helper struct to capture both the deployment commit
// and the associated merge_sha from the commits_diffs join query.
type deploymentCommitWithMergeSha struct {
	devops.CicdDeploymentCommit
	MergeSha string `gorm:"column:merge_sha"`
}

// batchFetchFirstCommits retrieves the first commit for all pull requests in the given project.
// Returns a map indexed by PR ID for O(1) lookup performance.
//
// The query uses a subquery to find the minimum commit_authored_date for each PR,
// then joins back to get the full commit record. This is more efficient than
// fetching all commits and filtering in memory.
func batchFetchFirstCommits(projectName string, db dal.Dal) (map[string]*code.PullRequestCommit, errors.Error) {
	var results []*code.PullRequestCommit

	// Use a subquery to find the earliest commit for each PR, then join to get full commit details.
	// This avoids scanning all commits and is optimized by the database engine.
	err := db.All(
		&results,
		dal.Select("prc.*"),
		dal.From("pull_request_commits prc"),
		dal.Join(`INNER JOIN (
			SELECT pull_request_id, MIN(commit_authored_date) as min_date
			FROM pull_request_commits
			GROUP BY pull_request_id
		) first_commits ON prc.pull_request_id = first_commits.pull_request_id
		AND prc.commit_authored_date = first_commits.min_date`),
		dal.Join("INNER JOIN pull_requests pr ON pr.id = prc.pull_request_id"),
		dal.Join("LEFT JOIN project_mapping pm ON pm.row_id = pr.base_repo_id AND pm.table = 'repos'"),
		dal.Where("pm.project_name = ?", projectName),
		dal.Orderby("prc.pull_request_id, prc.commit_authored_date ASC"),
	)

	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to batch fetch first commits")
	}

	// Build the map for O(1) lookup by PR ID
	commitMap := make(map[string]*code.PullRequestCommit, len(results))
	for _, commit := range results {
		// Only keep the first commit if multiple commits have the same timestamp
		if _, exists := commitMap[commit.PullRequestId]; !exists {
			commitMap[commit.PullRequestId] = commit
		}
	}

	return commitMap, nil
}

// batchFetchFirstReviews retrieves the first review comment for all pull requests in the given project.
// Returns a map indexed by PR ID for O(1) lookup performance.
//
// The query uses a subquery to find the minimum created_date for each PR (excluding the PR author),
// then joins back to get the full comment record.
func batchFetchFirstReviews(projectName string, db dal.Dal) (map[string]*code.PullRequestComment, errors.Error) {
	var results []*code.PullRequestComment

	// Use a subquery to find the earliest review comment for each PR (excluding author's comments),
	// then join to get full comment details.
	err := db.All(
		&results,
		dal.Select("prc.*"),
		dal.From("pull_request_comments prc"),
		dal.Join(`INNER JOIN (
			SELECT prc2.pull_request_id, MIN(prc2.created_date) as min_date
			FROM pull_request_comments prc2
			INNER JOIN pull_requests pr2 ON pr2.id = prc2.pull_request_id
			WHERE (pr2.author_id IS NULL OR pr2.author_id = '' OR prc2.account_id != pr2.author_id)
			GROUP BY prc2.pull_request_id
		) first_reviews ON prc.pull_request_id = first_reviews.pull_request_id
		AND prc.created_date = first_reviews.min_date`),
		dal.Join("INNER JOIN pull_requests pr ON pr.id = prc.pull_request_id"),
		dal.Join("LEFT JOIN project_mapping pm ON pm.row_id = pr.base_repo_id AND pm.table = 'repos'"),
		dal.Where("pm.project_name = ? AND (pr.author_id IS NULL OR pr.author_id = '' OR prc.account_id != pr.author_id)", projectName),
		dal.Orderby("prc.pull_request_id, prc.created_date ASC"),
	)

	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to batch fetch first reviews")
	}

	// Build the map for O(1) lookup by PR ID
	reviewMap := make(map[string]*code.PullRequestComment, len(results))
	for _, review := range results {
		// Only keep the first review if multiple reviews have the same timestamp
		if _, exists := reviewMap[review.PullRequestId]; !exists {
			reviewMap[review.PullRequestId] = review
		}
	}

	return reviewMap, nil
}

// batchFetchDeployments retrieves deployment commits for all merge commits in the given project.
// Returns a map indexed by merge commit SHA for O(1) lookup performance.
//
// The query finds the first successful production deployment for each merge commit by:
// 1. Finding deployment commits that have a previous successful deployment
// 2. Joining with commits_diffs to find which deployment included each merge commit
// 3. Filtering for successful production deployments
// 4. Ordering by started_date to get the earliest deployment
//
// The map is indexed by merge_sha (from commits_diffs), not by deployment commit_sha,
// because the caller needs to look up deployments by PR merge_commit_sha.
func batchFetchDeployments(projectName string, db dal.Dal) (map[string]*devops.CicdDeploymentCommit, errors.Error) {
	var results []*deploymentCommitWithMergeSha

	// Query finds the first deployment for each merge commit by using a window function
	// to rank deployments by started_date, then filtering to keep only rank 1.
	err := db.All(
		&results,
		dal.Select("dc.*, cd.commit_sha as merge_sha"),
		dal.From("cicd_deployment_commits dc"),
		dal.Join("LEFT JOIN cicd_deployment_commits p ON dc.prev_success_deployment_commit_id = p.id"),
		dal.Join("INNER JOIN commits_diffs cd ON cd.new_commit_sha = dc.commit_sha AND cd.old_commit_sha = COALESCE(p.commit_sha, '')"),
		dal.Join("LEFT JOIN project_mapping pm ON pm.table = 'cicd_scopes' AND pm.row_id = dc.cicd_scope_id"),
		dal.Where("dc.prev_success_deployment_commit_id <> ''"),
		dal.Where("dc.environment = 'PRODUCTION'"), // TODO: remove this when multi-environment is supported
		dal.Where("dc.result = ? AND pm.project_name = ?", devops.RESULT_SUCCESS, projectName),
		dal.Orderby("cd.commit_sha, dc.started_date ASC, dc.id ASC"),
	)

	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to batch fetch deployments")
	}

	// Build the map indexed by merge_sha for O(1) lookup.
	// Keep only the first deployment for each merge commit (earliest by started_date).
	deploymentMap := make(map[string]*devops.CicdDeploymentCommit, len(results))
	for _, result := range results {
		// Only keep the first deployment for each merge_sha
		if _, exists := deploymentMap[result.MergeSha]; !exists {
			// Copy the CicdDeploymentCommit without the MergeSha field
			deploymentCopy := result.CicdDeploymentCommit
			deploymentMap[result.MergeSha] = &deploymentCopy
		}
	}

	return deploymentMap, nil
}
