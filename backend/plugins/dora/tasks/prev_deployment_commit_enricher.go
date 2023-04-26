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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

var EnrichPrevSuccessDeploymentCommitMeta = plugin.SubTaskMeta{
	Name:             "enrichPrevSuccessDeploymentCommits",
	EntryPoint:       EnrichPrevSuccessDeploymentCommit,
	EnabledByDefault: false,
	Description:      "filling the prev_success_deployment_commit_id for cicd_deployment_commits table",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

// EnrichPrevSuccessDeploymentCommit
// Please note that deploying multiple environment (such as TESTING) copies
// (such as testing1 and testing2) using multiple steps with Deployment tools
// like Bitbucket or Gitlab is not supported and may result in incorrect
// outcomes. It is recommended that you deploy all copies in a single step.
// We arrived at this decision because we believe that deploying multiple
// environment copies using multiple steps is not a common or reasonable
// practice. However, if you have strong evidence to suggest otherwise, you are
// free to file an issue on our GitHub repository.
func EnrichPrevSuccessDeploymentCommit(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*DoraTaskData)
	// step 1. select all successful deployments in the project and sort them by cicd_scope_id, repo_url, env
	// and finished_date
	cursor, err := db.Cursor(
		dal.Select("dc.*"),
		dal.From("cicd_deployment_commits dc"),
		dal.Join("LEFT JOIN project_mapping pm ON (pm.table = 'cicd_scopes' AND pm.row_id = dc.cicd_scope_id)"),
		dal.Where(
			`
			dc.finished_date IS NOT NULL
			AND dc.environment IS NOT NULL AND dc.environment != ''
			AND dc.repo_url IS NOT NULL AND dc.repo_url != '' 
			AND pm.project_name = ? AND dc.result = ?
			`,
			data.Options.ProjectName, devops.SUCCESS,
		),
		dal.Orderby(`dc.cicd_scope_id, dc.repo_url, dc.environment, dc.finished_date`),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	prev_cicd_scope_id := ""
	prev_repo_url := ""
	prev_env := ""
	prev_success_deployment_id := ""

	enricher, err := api.NewDataEnricher(api.DataEnricherArgs[devops.CicdDeploymentCommit]{
		Ctx:   taskCtx,
		Name:  "prev_deployment_commit_id_enricher",
		Input: cursor,
		Enrich: func(deploymentCommit *devops.CicdDeploymentCommit) ([]interface{}, errors.Error) {
			// step 2. group them by cicd_scope_id/repo_url/env
			// whenever cicd_scope_id/repo_url/env shifted, it is a new set of consecutive deployments
			if prev_cicd_scope_id != deploymentCommit.CicdScopeId ||
				prev_repo_url != deploymentCommit.RepoUrl ||
				prev_env != deploymentCommit.Environment {
				// reset prev_success_deployment_id
				prev_success_deployment_id = ""
			}

			// now, simply connect the consecurtive deployment to its previous one
			deploymentCommit.PrevSuccessDeploymentCommitId = prev_success_deployment_id

			// preserve variables for the next record
			prev_cicd_scope_id = deploymentCommit.CicdScopeId
			prev_repo_url = deploymentCommit.RepoUrl
			prev_env = deploymentCommit.Environment
			prev_success_deployment_id = deploymentCommit.Id
			return []interface{}{deploymentCommit}, nil
		},
	})
	if err != nil {
		return err
	}

	return enricher.Execute()
}
