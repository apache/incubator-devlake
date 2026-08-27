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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/monorepo/models"
)

var AttributeDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "attributeDeployments",
	EntryPoint:       AttributeDeployments,
	EnabledByDefault: true,
	Description:      "Attribute each deployment to a monorepo sub-project by the name of the CI job that deployed it",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

// deploymentJobRow is one (deployment, deploy job) pair as returned by the query below.
//
// RawDataOrigin is embedded because DataConverter copies that field from the input row
// onto every result; without it the conversion panics.
type deploymentJobRow struct {
	common.RawDataOrigin
	CicdDeploymentId string
	CommitSha        string
	Result           string
	Environment      string
	FinishedDate     *time.Time
	JobName          string
}

// AttributeDeployments populates the core cicd_deployment_subprojects mapping table from
// cicd_deployment_commits + cicd_tasks regex matching, and - for one release, for backward
// compatibility - the monorepo plugin's own monorepo_subproject_deployments table with the
// exact same matches. There is no heuristic involved in either write (both are a direct
// regex match against the deploying job's name), so dual-writing is cheap and safe.
func AttributeDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*MonorepoTaskData)

	// Rebuild from scratch: attribution depends on configuration that may have changed
	// since the last run, so stale rows cannot be reconciled incrementally.
	if err := db.Exec(
		"DELETE FROM cicd_deployment_subprojects WHERE project_name = ?",
		data.Options.ProjectName,
	); err != nil {
		return errors.Default.Wrap(err, "error deleting previous cicd_deployment_subprojects")
	}
	if err := db.Exec(
		"DELETE FROM monorepo_subproject_deployments WHERE project_name = ?",
		data.Options.ProjectName,
	); err != nil {
		return errors.Default.Wrap(err, "error deleting previous monorepo_subproject_deployments")
	}

	// Only deployments generated from pipelines can be attributed: cicd_deployment_id is
	// the pipeline id, which is what cicd_tasks rows hang off. Deployments imported
	// straight from a provider's deployment API carry no job and are skipped.
	clauses := []dal.Clause{
		dal.Select(`dc.cicd_deployment_id, dc.commit_sha, dc.result, dc.environment,
			dc.finished_date, t.name AS job_name`),
		dal.From("cicd_deployment_commits dc"),
		dal.Join("JOIN project_mapping pm ON (pm.table = 'cicd_scopes' AND pm.row_id = dc.cicd_scope_id)"),
		dal.Join("JOIN cicd_tasks t ON (t.pipeline_id = dc.cicd_deployment_id)"),
		dal.Where("pm.project_name = ? AND t.type = ?", data.Options.ProjectName, devops.DEPLOYMENT),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	includeUnattributed := data.Options.ShouldIncludeUnattributed()
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: MonorepoApiParams{
				ProjectName: data.Options.ProjectName,
			},
			Table: "cicd_deployment_commits",
		},
		InputRowType: reflect.TypeOf(deploymentJobRow{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			row := inputRow.(*deploymentJobRow)
			matched := data.Matcher.MatchDeployJob(row.JobName)
			if len(matched) == 0 {
				if !includeUnattributed {
					// No sub-project matches and the caller opted out of the
					// 'unattributed' bucket: behave as before and skip the row.
					return nil, nil
				}
				matched = []string{UnattributedSubProject}
			}

			results := make([]interface{}, 0, len(matched)*2)
			for _, subProject := range matched {
				results = append(results, &devops.CicdDeploymentSubproject{
					ProjectName:      data.Options.ProjectName,
					CicdDeploymentId: row.CicdDeploymentId,
					SubProject:       subProject,
				})
				results = append(results, &models.SubProjectDeployment{
					ProjectName:      data.Options.ProjectName,
					SubProject:       subProject,
					CicdDeploymentId: row.CicdDeploymentId,
					CommitSha:        row.CommitSha,
					JobName:          row.JobName,
					Result:           row.Result,
					Environment:      row.Environment,
					FinishedDate:     row.FinishedDate,
				})
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
