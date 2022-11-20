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
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"reflect"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"
)

var ConnectIncidentToDeploymentMeta = core.SubTaskMeta{
	Name:             "ConnectIncidentToDeployment",
	EntryPoint:       ConnectIncidentToDeployment,
	EnabledByDefault: true,
	Description:      "Connect incident issue to deployment",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ConnectIncidentToDeployment(taskCtx core.SubTaskContext) errors.Error {
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

	enricher, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
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
			cicdTask := &devops.CICDTask{}
			cicdTakClauses := []dal.Clause{
				dal.From(cicdTask),
				dal.Join("left join project_mapping pm on cicd_tasks.cicd_scope_id = pm.row_id"),
				dal.Where(
					`cicd_tasks.finished_date < ? 
								and cicd_tasks.result = ? 
								and cicd_tasks.environment = ?
								and cicd_tasks.type = ?
								and pm.table = ?
								and pm.project_name = ?`,
					issue.CreatedDate, "SUCCESS", devops.PRODUCTION, "DEPLOYMENT", "cicd_scopes", data.Options.ProjectName,
				),
				dal.Orderby("cicd_tasks.finished_date DESC"),
			}
			err = db.First(cicdTask, cicdTakClauses...)
			if err != nil {
				if goerror.Is(err, gorm.ErrRecordNotFound) {
					return nil, nil
				} else {
					return nil, err
				}
			}
			projectIssueMetric.DeploymentId = cicdTask.Id

			return []interface{}{projectIssueMetric}, nil
		},
	})
	if err != nil {
		return err
	}

	return enricher.Execute()
}
