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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/github/models"
	githubTasks "github.com/apache/incubator-devlake/plugins/github/tasks"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/merico-dev/graphql"
	"reflect"
	"time"
)

const RAW_CHECK_RUNS_TABLE = "github_graphql_check_runs"

type GraphqlQueryCheckRunWrapper struct {
	RateLimit struct {
		Cost int
	}
	Node []GraphqlQueryCheckSuite `graphql:"node(id: $id)" graphql-extend:"true"`
}

type GraphqlQueryCheckSuite struct {
	Id         string
	Typename   string `graphql:"__typename"`
	CheckSuite struct {
		WorkflowRun struct {
			DatabaseId int
		}
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

var CollectCheckRunMeta = core.SubTaskMeta{
	Name:             "CollectCheckRun",
	EntryPoint:       CollectCheckRun,
	EnabledByDefault: true,
	Description:      "Collect CheckRun data from GithubGraphql api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

var _ core.SubTaskEntryPoint = CollectAccount

func CollectCheckRun(taskCtx core.SubTaskContext) errors.Error {
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*githubTasks.GithubTaskData)

	cursor, err := db.Cursor(
		dal.Select("check_suite_node_id"),
		dal.From(models.GithubRun{}.TableName()),
		dal.Where("repo_id = ? and connection_id=?", data.Options.GithubId, data.Options.ConnectionId),
	)
	if err != nil {
		return err
	}
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleWorkflowRun{}))
	if err != nil {
		return err
	}

	collector, err := helper.NewGraphqlCollector(helper.GraphqlCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: githubTasks.GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_CHECK_RUNS_TABLE,
		},
		Input:         iterator,
		InputStep:     60,
		GraphqlClient: data.GraphqlClient,
		BuildQuery: func(reqData *helper.GraphqlRequestData) (interface{}, map[string]interface{}, error) {
			workflowRuns := reqData.Input.([]interface{})
			query := &GraphqlQueryCheckRunWrapper{}
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
			query := iQuery.(*GraphqlQueryCheckRunWrapper)
			nodes := query.Node

			results := make([]interface{}, 0, 1)
			for _, node := range nodes {
				for _, checkRun := range node.CheckSuite.CheckRuns.Nodes {

					paramsBytes, err := json.Marshal(checkRun.Steps.Nodes)
					if err != nil {
						logger.Error(err, `Marshal checkRun.Steps.Nodes fail and ignore`)
					}
					githubJob := &models.GithubJob{
						ConnectionId: data.Options.ConnectionId,
						RunID:        node.CheckSuite.WorkflowRun.DatabaseId,
						RepoId:       data.Options.GithubId,
						ID:           checkRun.DatabaseId,
						NodeID:       checkRun.Id,
						HTMLURL:      checkRun.DetailsUrl,
						Status:       checkRun.Status,
						Conclusion:   checkRun.Conclusion,
						StartedAt:    checkRun.StartedAt,
						CompletedAt:  checkRun.CompletedAt,
						Name:         checkRun.Name,
						Steps:        paramsBytes,
						// these columns can not fill by graphql
						//HeadSha:       ``,  // use _tool_github_runs
						//RunURL:        ``,
						//CheckRunURL:   ``,
						//Labels:        ``, // not on use
						//RunnerID:      ``, // not on use
						//RunnerName:    ``, // not on use
						//RunnerGroupID: ``, // not on use
					}
					results = append(results, githubJob)
				}
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}
