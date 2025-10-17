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
	"encoding/json"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
	githubTasks "github.com/apache/incubator-devlake/plugins/github/tasks"
	"github.com/merico-dev/graphql"
)

const RAW_GRAPHQL_JOBS_TABLE = "github_graphql_jobs"

// Collection mode configuration
const (
	JOB_COLLECTION_MODE_BATCHING   = "BATCHING"
	JOB_COLLECTION_MODE_PAGINATING = "PAGINATING"
)

// Set the collection mode here
// BATCHING: Query multiple runs at once, no pagination (may miss jobs if >20 per run)
// PAGINATING: Query one run at a time with full pagination (complete data, more API calls)
const JOB_COLLECTION_MODE = JOB_COLLECTION_MODE_PAGINATING

// Mode-specific configuration
const (
	BATCHING_INPUT_STEP   = 10 // Number of runs per request in BATCHING mode
	BATCHING_PAGE_SIZE    = 20 // Jobs per run in BATCHING mode (no pagination)
	PAGINATING_INPUT_STEP = 1  // Number of runs per request in PAGINATING mode
	PAGINATING_PAGE_SIZE  = 50 // Jobs per page in PAGINATING mode (with pagination)
)

// Batch mode: query multiple runs at once (array of nodes)
type GraphqlQueryCheckRunWrapperBatch struct {
	RateLimit struct {
		Cost int
	}
	Node []GraphqlQueryCheckSuite `graphql:"node(id: $id)" graphql-extend:"true"`
}

// Paginating mode: query single run (single node)
type GraphqlQueryCheckRunWrapperSingle struct {
	RateLimit struct {
		Cost int
	}
	Node GraphqlQueryCheckSuite `graphql:"node(id: $id)"`
}

type GraphqlQueryCheckSuite struct {
	Id       string
	Typename string `graphql:"__typename"`
	// equal to Run in rest
	CheckSuite struct {
		WorkflowRun struct {
			DatabaseId int
		}
		// equal to Job in rest
		CheckRuns struct {
			TotalCount int
			PageInfo   struct {
				EndCursor   string `graphql:"endCursor"`
				HasNextPage bool   `graphql:"hasNextPage"`
			}
			Nodes []GraphqlQueryCheckRun
		} `graphql:"checkRuns(first: $pageSize, after: $skipCursor)"`
	} `graphql:"... on CheckSuite"`
}

type GraphqlQueryCheckRun struct {
	Id          string
	Name        string
	DetailsUrl  string
	DatabaseId  int
	Status      string
	StartedAt   *time.Time
	Conclusion  string
	CompletedAt *time.Time
	// ExternalId   string
	// Url          string
	// Title        interface{}
	// Text         interface{}
	// Summary      interface{}

	Steps struct {
		TotalCount int
		Nodes      []struct {
			CompletedAt         *time.Time `json:"completed_at"`
			Conclusion          string     `json:"conclusion"`
			Name                string     `json:"name"`
			Number              int        `json:"number"`
			SecondsToCompletion int        `json:"seconds_to_completion"`
			StartedAt           *time.Time `json:"started_at"`
			Status              string     `json:"status"`
		}
	} `graphql:"steps(first: 50)"`
}

type SimpleWorkflowRun struct {
	CheckSuiteNodeID string
}

// DbCheckRun is used to store additional fields (like RunId) required for database storage
// and application logic, while embedding the GraphqlQueryCheckRun struct for API data.
type DbCheckRun struct {
	RunId int // WorkflowRunId, required for DORA calculation
	*GraphqlQueryCheckRun
}

var CollectJobsMeta = plugin.SubTaskMeta{
	Name:             "Collect Job Runs",
	EntryPoint:       CollectJobs,
	EnabledByDefault: true,
	Description:      "Collect Jobs(CheckRun) data from GithubGraphql api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

var _ plugin.SubTaskEntryPoint = CollectJobs

func getPageInfo(query interface{}, args *helper.GraphqlCollectorArgs) (*helper.GraphqlQueryPageInfo, error) {
	// Only PAGINATING mode supports pagination
	if JOB_COLLECTION_MODE == JOB_COLLECTION_MODE_PAGINATING {
		queryWrapper := query.(*GraphqlQueryCheckRunWrapperSingle)
		return &helper.GraphqlQueryPageInfo{
			EndCursor:   queryWrapper.Node.CheckSuite.CheckRuns.PageInfo.EndCursor,
			HasNextPage: queryWrapper.Node.CheckSuite.CheckRuns.PageInfo.HasNextPage,
		}, nil
	}

	// BATCHING mode: no pagination support
	// Always return false for HasNextPage to collect only first page of jobs
	return &helper.GraphqlQueryPageInfo{
		EndCursor:   "",
		HasNextPage: false,
	}, nil
}

func buildQuery(reqData *helper.GraphqlRequestData) (interface{}, map[string]interface{}, error) {
	if reqData == nil {
		// Return appropriate empty query based on mode
		if JOB_COLLECTION_MODE == JOB_COLLECTION_MODE_PAGINATING {
			return &GraphqlQueryCheckRunWrapperSingle{}, map[string]interface{}{}, nil
		}
		return &GraphqlQueryCheckRunWrapperBatch{}, map[string]interface{}{}, nil
	}

	if JOB_COLLECTION_MODE == JOB_COLLECTION_MODE_PAGINATING {
		// Single run mode
		workflowRun := reqData.Input.(*SimpleWorkflowRun)
		query := &GraphqlQueryCheckRunWrapperSingle{}
		variables := map[string]interface{}{
			"id":         graphql.ID(workflowRun.CheckSuiteNodeID),
			"pageSize":   graphql.Int(reqData.Pager.Size),
			"skipCursor": (*graphql.String)(reqData.Pager.SkipCursor),
		}
		return query, variables, nil
	}

	// Batch mode (default)
	workflowRuns := reqData.Input.([]interface{})
	query := &GraphqlQueryCheckRunWrapperBatch{}
	checkSuiteIds := []map[string]interface{}{}
	for _, iWorkflowRuns := range workflowRuns {
		workflowRun := iWorkflowRuns.(*SimpleWorkflowRun)
		checkSuiteIds = append(checkSuiteIds, map[string]interface{}{
			`id`: graphql.ID(workflowRun.CheckSuiteNodeID),
		})
	}
	variables := map[string]interface{}{
		"node":       checkSuiteIds,
		"pageSize":   graphql.Int(reqData.Pager.Size),
		"skipCursor": (*graphql.String)(reqData.Pager.SkipCursor),
	}
	return query, variables, nil
}

func CollectJobs(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*githubTasks.GithubTaskData)
	logger := taskCtx.GetLogger()

	// Log the collection mode
	logger.Info("GitHub Job Collector Mode: %s", JOB_COLLECTION_MODE)

	apiCollector, err := helper.NewStatefulApiCollector(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: githubTasks.GithubApiParams{
			ConnectionId: data.Options.ConnectionId,
			Name:         data.Options.Name,
		},
		Table: RAW_GRAPHQL_JOBS_TABLE,
	})
	if err != nil {
		return err
	}

	clauses := []dal.Clause{
		dal.Select("check_suite_node_id"),
		dal.From(models.GithubRun{}.TableName()),
		dal.Where("repo_id = ? and connection_id=?", data.Options.GithubId, data.Options.ConnectionId),
		dal.Orderby("github_updated_at DESC"),
	}
	if apiCollector.IsIncremental() && apiCollector.GetSince() != nil {
		clauses = append(clauses, dal.Where("github_updated_at > ?", *apiCollector.GetSince()))
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleWorkflowRun{}))
	if err != nil {
		return err
	}

	// Set configuration based on mode
	var inputStep, pageSize int
	var getPageInfoFunc func(interface{}, *helper.GraphqlCollectorArgs) (*helper.GraphqlQueryPageInfo, error)

	if JOB_COLLECTION_MODE == JOB_COLLECTION_MODE_PAGINATING {
		inputStep = PAGINATING_INPUT_STEP
		pageSize = PAGINATING_PAGE_SIZE
		getPageInfoFunc = getPageInfo // Enable pagination
	} else {
		inputStep = BATCHING_INPUT_STEP
		pageSize = BATCHING_PAGE_SIZE
		getPageInfoFunc = nil // Disable pagination
	}

	err = apiCollector.InitGraphQLCollector(helper.GraphqlCollectorArgs{
		Input:         iterator,
		InputStep:     inputStep,
		GraphqlClient: data.GraphqlClient,
		BuildQuery:    buildQuery,
		GetPageInfo:   getPageInfoFunc, // nil for BATCHING, function for PAGINATING
		ResponseParser: func(queryWrapper any) (messages []json.RawMessage, err errors.Error) {
			if JOB_COLLECTION_MODE == JOB_COLLECTION_MODE_PAGINATING {
				// Single node processing
				query := queryWrapper.(*GraphqlQueryCheckRunWrapperSingle)
				node := query.Node
				runId := node.CheckSuite.WorkflowRun.DatabaseId

				for _, checkRun := range node.CheckSuite.CheckRuns.Nodes {
					dbCheckRun := &DbCheckRun{
						RunId:                runId,
						GraphqlQueryCheckRun: &checkRun,
					}
					// A checkRun without a startedAt time is a run that was never started (skipped), GitHub returns
					// a ZeroTime (Due to the GO implementation) for startedAt, so we need to check for that here.
					dbCheckRun.StartedAt = utils.NilIfZeroTime(dbCheckRun.StartedAt)
					dbCheckRun.CompletedAt = utils.NilIfZeroTime(dbCheckRun.CompletedAt)
					updatedAt := dbCheckRun.StartedAt
					if dbCheckRun.CompletedAt != nil {
						updatedAt = dbCheckRun.CompletedAt
					}
					if apiCollector.GetSince() != nil && !apiCollector.GetSince().Before(*updatedAt) {
						return messages, helper.ErrFinishCollect
					}
					messages = append(messages, errors.Must1(json.Marshal(dbCheckRun)))
				}
			} else {
				// Batch processing (multiple nodes)
				query := queryWrapper.(*GraphqlQueryCheckRunWrapperBatch)
				for _, node := range query.Node {
					runId := node.CheckSuite.WorkflowRun.DatabaseId
					for _, checkRun := range node.CheckSuite.CheckRuns.Nodes {
						dbCheckRun := &DbCheckRun{
							RunId:                runId,
							GraphqlQueryCheckRun: &checkRun,
						}
						// A checkRun without a startedAt time is a run that was never started (skipped), GitHub returns
						// a ZeroTime (Due to the GO implementation) for startedAt, so we need to check for that here.
						dbCheckRun.StartedAt = utils.NilIfZeroTime(dbCheckRun.StartedAt)
						dbCheckRun.CompletedAt = utils.NilIfZeroTime(dbCheckRun.CompletedAt)
						updatedAt := dbCheckRun.StartedAt
						if dbCheckRun.CompletedAt != nil {
							updatedAt = dbCheckRun.CompletedAt
						}
						if apiCollector.GetSince() != nil && !apiCollector.GetSince().Before(*updatedAt) {
							return messages, helper.ErrFinishCollect
						}
						messages = append(messages, errors.Must1(json.Marshal(dbCheckRun)))
					}
				}
			}
			return
		},
		IgnoreQueryErrors: true,
		PageSize:          pageSize,
	})
	if err != nil {
		return err
	}

	return apiCollector.Execute()
}
