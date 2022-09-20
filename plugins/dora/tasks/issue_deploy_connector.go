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
	"reflect"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"gorm.io/gorm"
)

var ConnectIssueDeployMeta = core.SubTaskMeta{
	Name:             "ConnectIssueDeploy",
	EntryPoint:       ConnectIssueDeploy,
	EnabledByDefault: true,
	Description:      "TODO",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ConnectIssueDeploy(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*DoraTaskData)

	issue := &ticket.Issue{}
	// select all issues belongs to the board
	clauses := []dal.Clause{
		dal.From(issue),
		dal.Join(`left join board_issues 
			on issues.id = board_issues.issue_id`),
		dal.Join("left join board_repos on board_repos.board_id = board_issues.board_id"),
		dal.Where(
			"board_repos.repo_id = ? and issues.type = ?",
			data.Options.RepoId, "INCIDENT",
		),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	enricher, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:    taskCtx,
			Params: DoraApiParams{
				// TODO
			},
			Table: "issues",
		},
		InputRowType: reflect.TypeOf(ticket.Issue{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			issueToBeUpdate := inputRow.(*ticket.Issue)
			cicdTask := &devops.CICDTask{}
			cicdTakClauses := []dal.Clause{
				dal.From(cicdTask),
				dal.Join("left join cicd_pipeline_commits on cicd_tasks.pipeline_id = cicd_pipeline_commits.pipeline_id"),
				dal.Where(
					`cicd_tasks.finished_date < ? 
								and cicd_tasks.result = ? and cicd_tasks.environment = ?`,
					issueToBeUpdate.CreatedDate, "SUCCESS", devops.PRODUCTION,
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
			issueToBeUpdate.DeploymentId = cicdTask.Id

			return []interface{}{issueToBeUpdate}, nil
		},
	})
	if err != nil {
		return err
	}

	return enricher.Execute()
}
