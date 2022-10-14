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
	"github.com/apache/incubator-devlake/errors"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	bitbucketModels "github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertDeploymentMeta = core.SubTaskMeta{
	Name:             "convertDeployments",
	EntryPoint:       ConvertDeployments,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bitbucket_pipeline into domain layer table pipeline",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

func ConvertDeployments(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*BitbucketTaskData)

	cursor, err := db.Cursor(dal.From(bitbucketModels.BitbucketDeployment{}))
	if err != nil {
		return err
	}
	defer cursor.Close()

	pipelineIdGen := didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketDeployment{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(bitbucketModels.BitbucketDeployment{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: BitbucketApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_DEPLOYMENT_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			bitbucketDeployment := inputRow.(*bitbucketModels.BitbucketDeployment)

			startedAt := bitbucketDeployment.CreatedOn
			if bitbucketDeployment.StartedOn != nil {
				startedAt = bitbucketDeployment.StartedOn
			}
			domainDeployment := &devops.CICDTask{
				DomainEntity: domainlayer.DomainEntity{
					Id: pipelineIdGen.Generate(data.Options.ConnectionId, bitbucketDeployment.BitbucketId),
				},
				Name: didgen.NewDomainIdGenerator(&bitbucketModels.BitbucketPipeline{}).
					Generate(data.Options.ConnectionId, bitbucketDeployment.Name),
				PipelineId: bitbucketDeployment.PipelineId,
				Result: devops.GetResult(&devops.ResultRule{
					Failed:  []string{bitbucketModels.FAILED, bitbucketModels.ERROR, bitbucketModels.UNDEPLOYED},
					Abort:   []string{bitbucketModels.STOPPED, bitbucketModels.SKIPPED},
					Success: []string{bitbucketModels.SUCCESSFUL, bitbucketModels.COMPLETED},
					Manual:  []string{bitbucketModels.PAUSED, bitbucketModels.HALTED},
					Default: devops.SUCCESS,
				}, bitbucketDeployment.Status),
				Status: devops.GetStatus(&devops.StatusRule{
					InProgress: []string{bitbucketModels.IN_PROGRESS, bitbucketModels.PENDING, bitbucketModels.BUILDING},
					Default:    devops.DONE,
				}, bitbucketDeployment.Status),
				Type:         bitbucketDeployment.Type,
				StartedDate:  *startedAt,
				FinishedDate: bitbucketDeployment.CompletedOn,
			}
			// rebuild the FinishedDate and DurationSec by Status
			finishedAt := time.Now()
			if domainDeployment.Status != devops.DONE {
				domainDeployment.FinishedDate = nil
			} else if bitbucketDeployment.CompletedOn != nil {
				finishedAt = *bitbucketDeployment.CompletedOn
			}
			durationTime := finishedAt.Sub(*startedAt)
			domainDeployment.DurationSec = uint64(durationTime.Seconds())

			return []interface{}{
				domainDeployment,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
