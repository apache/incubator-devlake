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
	"fmt"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

var ConnectIncidentToDeploymentMeta = plugin.SubTaskMeta{
	Name:             "ConnectIncidentToDeployment",
	EntryPoint:       ConnectIncidentToDeployment,
	EnabledByDefault: true,
	Description:      "Connect incident issue to deployment",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

type simpleCicdDeploymentCommit struct {
	Id           string
	FinishedDate *time.Time
}

func ConnectIncidentToDeployment(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*DoraTaskData)
	// select all issues belongs to the board
	clauses := []dal.Clause{
		dal.From(`issues i`),
		dal.Join(`left join board_issues bi on bi.issue_id = i.id`),
		dal.Join(`left join project_mapping pm on pm.row_id = bi.board_id`),
		dal.Where(
			"i.type = ? and pm.project_name = ? and pm.table = ?",
			"INCIDENT", data.Options.ProjectName, "boards",
		),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	count, err := db.Count(clauses...)
	if err != nil {
		return errors.Default.Wrap(err, "error getting count of clauses")
	}
	if count == 0 {
		// Clear previous results from the project
		deleteSql := fmt.Sprintf("DELETE FROM project_issue_metrics WHERE project_name = '%s'", data.Options.ProjectName)
		err := db.Exec(deleteSql)
		if err != nil {
			return errors.Default.Wrap(err, "error deleting previous project_issue_metrics")
		}
		return nil
	}

	enricher, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: DoraApiParams{
				ProjectName: data.Options.ProjectName,
			},
			Table: "issues",
		},
		InputRowType: reflect.TypeOf(ticket.Issue{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			issue := inputRow.(*ticket.Issue)
			projectIssueMetric := &crossdomain.ProjectIssueMetric{
				DomainEntity: domainlayer.DomainEntity{
					Id: issue.Id,
				},
				ProjectName: data.Options.ProjectName,
			}

			cicdDeploymentCommit := &devops.CicdDeploymentCommit{}
			cicdDeploymentCommitClauses := []dal.Clause{
				dal.Select("cicd_deployment_commits.cicd_deployment_id as id, cicd_deployment_commits.finished_date as finished_date"),
				dal.From(cicdDeploymentCommit),
				dal.Join("left join project_mapping pm on cicd_deployment_commits.cicd_scope_id = pm.row_id"),
				dal.Where(
					`cicd_deployment_commits.finished_date < ?
					    and cicd_deployment_commits.result = ?
						and cicd_deployment_commits.environment = ?
						and pm.table = ?
						and pm.project_name = ?`,
					issue.CreatedDate, devops.RESULT_SUCCESS, devops.PRODUCTION, "cicd_scopes", data.Options.ProjectName,
				),
				dal.Orderby("finished_date DESC"),
				dal.Limit(1),
			}

			scdc := &simpleCicdDeploymentCommit{}
			err = db.All(scdc, cicdDeploymentCommitClauses...)
			if err != nil {
				if db.IsErrorNotFound(err) {
					return nil, nil
				} else {
					return nil, err
				}
			}
			if scdc.Id != "" {
				projectIssueMetric.DeploymentId = scdc.Id
				return []interface{}{projectIssueMetric}, nil
			}
			return nil, nil
		},
	})
	if err != nil {
		return err
	}

	return enricher.Execute()
}
