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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertReleaseDeploymentsMeta)
}

var ConvertReleaseDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "convertApiReleaseDeployments",
	EntryPoint:       ConvertReleaseDeployments,
	EnabledByDefault: true,
	Description:      "Convert tool layer table azuredevops_release_deployments into domain layer table cicd_pipelines",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{
		models.AzuredevopsReleaseDeployment{}.TableName(),
	},
}

// Release deployment status and operation status mappings
// Reference: https://learn.microsoft.com/en-us/rest/api/azure/devops/release/deployments/list
const (
	releaseDeploymentStatusSucceeded        = "succeeded"
	releaseDeploymentStatusFailed           = "failed"
	releaseDeploymentStatusNotDeployed      = "notDeployed"
	releaseDeploymentStatusPartiallySucceeded = "partiallySucceeded"
	releaseOperationStatusApproved          = "Approved"
	releaseOperationStatusCanceled          = "Canceled"
	releaseOperationStatusCancelling        = "Cancelling"
	releaseOperationStatusDeferred          = "Deferred"
	releaseOperationStatusEvaluatingGates   = "EvaluatingGates"
	releaseOperationStatusGateFailed        = "GateFailed"
	releaseOperationStatusManualInterventionPending = "ManualInterventionPending"
	releaseOperationStatusPending           = "Pending"
	releaseOperationStatusPhaseCanceled     = "PhaseCanceled"
	releaseOperationStatusPhaseFailed       = "PhaseFailed"
	releaseOperationStatusPhaseInProgress   = "PhaseInProgress"
	releaseOperationStatusPhasePartiallySucceeded = "PhasePartiallySucceeded"
	releaseOperationStatusPhaseSucceeded    = "PhaseSucceeded"
	releaseOperationStatusQueued            = "Queued"
	releaseOperationStatusRejected          = "Rejected"
	releaseOperationStatusScheduled         = "Scheduled"
	releaseOperationStatusUndefined         = "Undefined"
)

var releaseDeploymentResultRule = devops.ResultRule{
	Success: []string{releaseDeploymentStatusSucceeded, releaseDeploymentStatusPartiallySucceeded},
	Failure: []string{releaseDeploymentStatusFailed, releaseDeploymentStatusNotDeployed},
	Default: devops.RESULT_DEFAULT,
}

var releaseDeploymentStatusRule = devops.StatusRule{
	Done:       []string{releaseOperationStatusApproved, releaseOperationStatusCanceled, releaseOperationStatusRejected, releaseOperationStatusPhaseCanceled, releaseOperationStatusPhaseFailed, releaseOperationStatusPhasePartiallySucceeded, releaseOperationStatusPhaseSucceeded, releaseOperationStatusGateFailed},
	InProgress: []string{releaseOperationStatusPending, releaseOperationStatusQueued, releaseOperationStatusScheduled, releaseOperationStatusDeferred, releaseOperationStatusCancelling, releaseOperationStatusEvaluatingGates, releaseOperationStatusManualInterventionPending, releaseOperationStatusPhaseInProgress, releaseOperationStatusUndefined},
	Default:    devops.STATUS_OTHER,
}

func ConvertReleaseDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawReleaseDeploymentTable)
	clauses := []dal.Clause{
		dal.Select("*"),
		dal.From(&models.AzuredevopsReleaseDeployment{}),
		dal.Where("project_id = ? and connection_id = ?",
			data.Options.ProjectId, data.Options.ConnectionId),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	deploymentIdGen := didgen.NewDomainIdGenerator(&models.AzuredevopsReleaseDeployment{})
	projectIdGen := didgen.NewDomainIdGenerator(&models.AzuredevopsProject{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.AzuredevopsReleaseDeployment{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			deployment := inputRow.(*models.AzuredevopsReleaseDeployment)
			duration := 0.0

			if deployment.CompletedOn != nil && deployment.StartedOn != nil {
				duration = float64(deployment.CompletedOn.Sub(*deployment.StartedOn).Milliseconds() / 1e3)
			}

			// Create a unique pipeline name combining release definition and environment
			pipelineName := deployment.DefinitionName
			if deployment.EnvironmentName != "" {
				pipelineName = deployment.DefinitionName + " - " + deployment.EnvironmentName
			}

			domainPipeline := &devops.CICDPipeline{
				DomainEntity: domainlayer.DomainEntity{
					Id: deploymentIdGen.Generate(data.Options.ConnectionId, deployment.AzuredevopsId),
				},
				Name:           pipelineName,
				Result:         devops.GetResult(&releaseDeploymentResultRule, deployment.DeploymentStatus),
				Status:         devops.GetStatus(&releaseDeploymentStatusRule, deployment.OperationStatus),
				OriginalStatus: deployment.OperationStatus,
				OriginalResult: deployment.DeploymentStatus,
				CicdScopeId:    projectIdGen.Generate(data.Options.ConnectionId, data.Options.ProjectId),
				Environment:    data.RegexEnricher.ReturnNameIfMatched(devops.PRODUCTION, pipelineName+";"+deployment.EnvironmentName),
				Type:           devops.DEPLOYMENT,
				DurationSec:    duration,
			}

			if deployment.QueuedOn != nil {
				domainPipeline.TaskDatesInfo = devops.TaskDatesInfo{
					CreatedDate:  *deployment.QueuedOn,
					QueuedDate:   deployment.QueuedOn,
					StartedDate:  deployment.StartedOn,
					FinishedDate: deployment.CompletedOn,
				}
			}

			return []interface{}{
				domainPipeline,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
