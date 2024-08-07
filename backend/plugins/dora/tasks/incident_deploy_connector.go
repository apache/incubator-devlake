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

// ConnectIncidentToDeployment will generate data to crossdomain.ProjectIncidentDeploymentRelationship.
func ConnectIncidentToDeployment(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*DoraTaskData)
	logger := taskCtx.GetLogger()
	// Clear previous results from the project
	err := db.Exec("DELETE FROM project_incident_deployment_relationships WHERE project_name = ?", data.Options.ProjectName)
	if err != nil {
		return errors.Default.Wrap(err, "error deleting previous project_incident_deployment_relationships")
	}
	logger.Info("delete previous project_incident_deployment_relationships")
	// select all issues belongs to the board
	clauses := []dal.Clause{
		dal.From(`incidents i`),
		dal.Join(`left join project_mapping pm on pm.row_id = i.scope_id and pm.table = i.table`),
		dal.Where("pm.project_name = ?", data.Options.ProjectName),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		logger.Error(err, "db.cursor error")
		return err
	}
	defer cursor.Close()
	logger.Info("start enricher")
	enricher, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: DoraApiParams{
				ProjectName: data.Options.ProjectName,
			},
			Table: "incidents",
		},
		InputRowType: reflect.TypeOf(ticket.Incident{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			incident := inputRow.(*ticket.Incident)
			projectIssueMetric := &crossdomain.ProjectIncidentDeploymentRelationship{
				DomainEntity: domainlayer.DomainEntity{
					Id: incident.Id,
				},
				ProjectName: data.Options.ProjectName,
			}
			logger.Info("get incident: %+v", incident.Id)
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
					incident.CreatedDate, devops.RESULT_SUCCESS, devops.PRODUCTION, "cicd_scopes", data.Options.ProjectName,
				),
				dal.Orderby("finished_date DESC"),
				dal.Limit(1),
			}

			scdc := &simpleCicdDeploymentCommit{}
			err = db.All(scdc, cicdDeploymentCommitClauses...)
			if err != nil {
				if db.IsErrorNotFound(err) {
					logger.Warn(err, "deployment commit not found")
					return nil, nil
				} else {
					logger.Error(err, "get all deployment commits")
					return nil, err
				}
			}
			if scdc.Id != "" {
				projectIssueMetric.DeploymentId = scdc.Id
				return []interface{}{projectIssueMetric}, nil
			}
			logger.Info("scdc.id is empty, incident will be ignored: %+v", incident.Id)
			return nil, nil
		},
	})
	if err != nil {
		return err
	}

	return enricher.Execute()
}
