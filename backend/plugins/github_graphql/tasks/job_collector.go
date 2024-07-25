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
			Nodes      []struct {
				Id          string
				Name        string
				DetailsUrl  string
				DatabaseId  int
				Status      string
				StartedAt   *time.Time
				Conclusion  string
				CompletedAt *time.Time
				//ExternalId   string
				//Url          string
				//Title        interface{}
				//Text         interface{}
				//Summary      interface{}

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
		} `graphql:"checkRuns(first: 50)"`
	} `graphql:"... on CheckSuite"`
}

type SimpleWorkflowRun struct {
	CheckSuiteNodeID string
}

var CollectJobsMeta = plugin.SubTaskMeta{
	Name:             "Collect Job Runs",
	EntryPoint:       CollectJobs,
	EnabledByDefault: true,
	Description:      "Collect Jobs(CheckRun) data from GithubGraphql api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

var _ plugin.SubTaskEntryPoint = CollectAccount

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
		InputStep:     20,
		GraphqlClient: data.GraphqlClient,
		BuildQuery: func(reqData *helper.GraphqlRequestData) (interface{}, map[string]interface{}, error) {
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
				"node": checkSuiteIds,
			}
			return query, variables, nil
		},
		ResponseParserWithDataErrors: func(iQuery interface{}, variables map[string]interface{}, dataErrors []graphql.DataError) ([]interface{}, error) {
			for _, dataError := range dataErrors {
				// log and ignore
				taskCtx.GetLogger().Warn(dataError, `query check run get error but ignore`)
			}
			return nil, nil
		},
	})
	if err != nil {
		return err
	}

	return apiCollector.Execute()
}
