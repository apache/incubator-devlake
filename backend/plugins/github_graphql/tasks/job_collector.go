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
	"fmt"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
	githubTasks "github.com/apache/incubator-devlake/plugins/github/tasks"
	"github.com/merico-dev/graphql"
)

const RAW_GRAPHQL_JOBS_TABLE = "github_graphql_jobs"

type GraphqlQueryCheckRunWrapper struct {
	RateLimit struct {
		Cost int
	}
	Node []GraphqlQueryCheckSuite `graphql:"node(id: $id)" graphql-extend:"true"`
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
	queryWrapper := query.(*GraphqlQueryCheckRunWrapper)
	hasNextPage := false
	endCursor := ""
	for _, node := range queryWrapper.Node {
		if node.CheckSuite.CheckRuns.PageInfo.HasNextPage {
			hasNextPage = true
			endCursor = node.CheckSuite.CheckRuns.PageInfo.EndCursor
			break
		}
	}
	return &helper.GraphqlQueryPageInfo{
		EndCursor:   endCursor,
		HasNextPage: hasNextPage,
	}, nil
}

func buildQuery(reqData *helper.GraphqlRequestData) (interface{}, map[string]interface{}, error) {
	query := &GraphqlQueryCheckRunWrapper{}
	if reqData == nil {
		return query, map[string]interface{}{}, nil
	}
	workflowRuns := reqData.Input.([]interface{})
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

	cursor, err := db.Cursor(
		clauses...,
	)
	if err != nil {
		return err
	}
	defer cursor.Close()
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleWorkflowRun{}))
	if err != nil {
		return err
	}

	err = apiCollector.InitGraphQLCollector(helper.GraphqlCollectorArgs{
		Input:         iterator,
		InputStep:     10,
		GraphqlClient: data.GraphqlClient,
		BuildQuery:    buildQuery,
		GetPageInfo:   getPageInfo,
		ResponseParser: func(queryWrapper any) (messages []json.RawMessage, err errors.Error) {
			query := queryWrapper.(*GraphqlQueryCheckRunWrapper)
			for _, node := range query.Node {
				runId := node.CheckSuite.WorkflowRun.DatabaseId
				for _, checkRun := range node.CheckSuite.CheckRuns.Nodes {
					dbCheckRun := &DbCheckRun{
						RunId:                runId,
						GraphqlQueryCheckRun: &checkRun,
					}
					// A checkRun without a startedAt time is a run that was never started (skipped)
					// akwardly, GitHub assigns a completedAt time to such runs, which is the time when the run was skipped
					// TODO: Decide if we want to skip those runs or should we assign the startedAt time to the completedAt time
					if dbCheckRun.StartedAt == nil || dbCheckRun.StartedAt.IsZero() {
						debug := fmt.Sprintf("collector: checkRun.StartedAt is nil or zero: %s", dbCheckRun.Id)
						taskCtx.GetLogger().Debug(debug, "Collector: CheckRun started at is nil or zero")
						continue
					}
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
			return
		},
		IgnoreQueryErrors: true,
		PageSize:          20,
	})
	if err != nil {
		return err
	}

	return apiCollector.Execute()
}
