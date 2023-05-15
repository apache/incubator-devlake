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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

var _ plugin.SubTaskEntryPoint = ExtractPipelineSteps

var ExtractPipelineStepsMeta = plugin.SubTaskMeta{
	Name:             "ExtractPipelineSteps",
	EntryPoint:       ExtractPipelineSteps,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table bitbucket_pipeline_steps",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

type BitbucketPipelineStepsResponse struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Uuid     string `json:"uuid"`
	Pipeline struct {
		//Type string `json:"type"`
		Uuid string `json:"uuid"`
	} `json:"pipeline"`
	Trigger struct {
		Type string `json:"type"`
	} `json:"trigger"`
	State struct {
		Name string `json:"name"`
		//Type   string `json:"type"`
		Result struct {
			Name string `json:"name"`
			//Type string `json:"type"`
		} `json:"result"`
	} `json:"state"`
	MaxTime           int        `json:"maxTime"`
	StartedOn         *time.Time `json:"started_on"`
	CompletedOn       *time.Time `json:"completed_on"`
	DurationInSeconds int        `json:"duration_in_seconds"`
	BuildSecondsUsed  int        `json:"build_seconds_used"`
	RunNumber         int        `json:"run_number"`
}

func ExtractPipelineSteps(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_STEPS_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,

		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			apiPipelineStep := &BitbucketPipelineStepsResponse{}
			err := errors.Convert(json.Unmarshal(resData.Data, apiPipelineStep))
			if err != nil {
				return nil, err
			}

			bitbucketStep := &models.BitbucketPipelineStep{
				ConnectionId:      data.Options.ConnectionId,
				BitbucketId:       apiPipelineStep.Uuid,
				PipelineId:        apiPipelineStep.Pipeline.Uuid,
				Name:              apiPipelineStep.Name,
				Trigger:           apiPipelineStep.Trigger.Type,
				State:             apiPipelineStep.State.Name,
				Result:            apiPipelineStep.State.Result.Name,
				RepoId:            data.Options.FullName,
				MaxTime:           apiPipelineStep.MaxTime,
				StartedOn:         apiPipelineStep.StartedOn,
				CompletedOn:       apiPipelineStep.CompletedOn,
				DurationInSeconds: apiPipelineStep.DurationInSeconds,
				BuildSecondsUsed:  apiPipelineStep.BuildSecondsUsed,
				RunNumber:         apiPipelineStep.RunNumber,
				Type:              data.RegexEnricher.ReturnNameIfMatched(devops.DEPLOYMENT, apiPipelineStep.Name),
				Environment:       data.RegexEnricher.ReturnNameIfOmittedOrMatched(devops.PRODUCTION, apiPipelineStep.Name),
			}
			return []interface{}{
				bitbucketStep,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}
