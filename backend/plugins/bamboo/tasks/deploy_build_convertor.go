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
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
)

var ConvertDeployBuildsMeta = plugin.SubTaskMeta{
	Name:             "convertDeployBuilds",
	EntryPoint:       ConvertDeployBuilds,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bamboo_deploy_builds into  domain layer table deployBuilds",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

type deployBuildWithVcsRevision struct {
	models.BambooDeployBuild
	RepositoryId    int
	RepositoryName  string
	VcsRevisionKey  string
	ProjectPlanName string
	ProjectName     string
}

func (deployBuildWithVcsRevision deployBuildWithVcsRevision) GenerateCICDDeploymentCommitName() string {
	if deployBuildWithVcsRevision.ProjectPlanName != "" {
		return fmt.Sprintf("%s/%s", deployBuildWithVcsRevision.ProjectPlanName, deployBuildWithVcsRevision.DeploymentVersionName)
	}
	return deployBuildWithVcsRevision.DeploymentVersionName
}

func ConvertDeployBuilds(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_JOB_BUILD_TABLE)
	cursor, err := db.Cursor(
		dal.Select("db.*, pbc.repository_id, pbc.repository_name, pbc.vcs_revision_key, p.name as project_plan_name, p.project_name"),
		dal.From("_tool_bamboo_deploy_builds AS db"),
		dal.Join("INNER JOIN _tool_bamboo_plan_build_commits AS pbc ON db.connection_id = pbc.connection_id AND db.plan_result_key = pbc.plan_result_key"),
		dal.Join("LEFT JOIN _tool_bamboo_plans as p ON db.plan_key = p.plan_key"),
		dal.Where("db.connection_id = ? and db.plan_key = ?", data.Options.ConnectionId, data.Options.PlanKey))
	if err != nil {
		return err
	}
	defer cursor.Close()

	deployBuildIdGen := didgen.NewDomainIdGenerator(&models.BambooDeployBuild{})
	planIdGen := didgen.NewDomainIdGenerator(&models.BambooPlan{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(deployBuildWithVcsRevision{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			input := inputRow.(*deployBuildWithVcsRevision)
			if input.VcsRevisionKey == "" {
				return nil, nil
			}
			deploymentCommit := &devops.CicdDeploymentCommit{
				DomainEntity: domainlayer.DomainEntity{
					Id: deployBuildIdGen.Generate(data.Options.ConnectionId, input.DeployBuildId),
				},
				CicdScopeId:      planIdGen.Generate(data.Options.ConnectionId, data.Options.PlanKey),
				CicdDeploymentId: deployBuildIdGen.Generate(data.Options.ConnectionId, input.DeployBuildId),
				Name:             input.GenerateCICDDeploymentCommitName(),
				Result: devops.GetResult(&devops.ResultRule{
					Failed:          []string{"Failed"},
					Success:         []string{"Successful"},
					CaseInsensitive: true,
					Default:         input.DeploymentState,
				}, input.DeploymentState),
				Status: devops.GetStatus(&devops.StatusRule[string]{
					Done:       []string{"Finished", "FINISHED"},
					NotStarted: []string{"not_built", "NOT_BUILT", "Not_Built", "PENDING", "QUEUED"},
					Default:    devops.STATUS_IN_PROGRESS,
				}, input.LifeCycleState),
				Environment:  input.Environment,
				StartedDate:  input.StartedDate,
				FinishedDate: input.FinishedDate,
				CommitSha:    input.VcsRevisionKey,
				RefName:      input.PlanBranchName,
				RepoId:       strconv.Itoa(input.RepositoryId),
			}
			deploymentCommit.CreatedDate = time.Now()
			if input.StartedDate != nil {
				deploymentCommit.CreatedDate = *input.StartedDate
			}
			if input.QueuedDate != nil {
				deploymentCommit.CreatedDate = *input.QueuedDate
			}
			if data.RegexEnricher.ReturnNameIfMatched(models.ENV_NAME_PATTERN, input.Environment) != "" {
				deploymentCommit.Environment = devops.PRODUCTION
			}
			if input.FinishedDate != nil && input.StartedDate != nil {
				duration := uint64(input.FinishedDate.Sub(*input.StartedDate).Seconds())
				deploymentCommit.DurationSec = &duration
			}
			fakeRepoUrl, err := generateFakeRepoUrl(data.ApiClient.GetEndpoint(), input.RepositoryId)
			if err != nil {
				logger.Warn(err, "generate fake repo url, endpoint: %s, repo id: %d", data.ApiClient.GetEndpoint(), input.RepositoryId)
			} else {
				deploymentCommit.RepoUrl = fakeRepoUrl
			}

			return []interface{}{deploymentCommit}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
