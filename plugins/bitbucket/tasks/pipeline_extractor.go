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
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"time"
)

type BitbucketApiCommit struct {
	Type  string `json:"type"`
	Hash  string `json:"hash"`
	Links struct {
		Self struct {
			Href string
		}
		Html struct {
			Href string
		}
	} `json:"links"`
}

type BitbucketApiPipelineTarget struct {
	Type     string `json:"type"`
	RefType  string `json:"ref_type"`
	RefName  string `json:"ref_name"`
	Selector struct {
		Type string
	} `json:"selector"`
	Commit *BitbucketApiCommit `json:"commit"`
}

type BitbucketApiPipeline struct {
	Uuid  string `json:"uuid"`
	Type  string `json:"type"`
	State struct {
		Name   string
		Type   string
		Result struct {
			Name string
			Type string
		}
	} `json:"state"`
	BuildNumber int                         `json:"build_number"`
	Creator     *BitbucketAccountResponse   `json:"creator"`
	Repo        *BitbucketApiRepo           `json:"repository"`
	Target      *BitbucketApiPipelineTarget `json:"target"`
	Trigger     struct {
		Name string
		Type string
	} `json:"trigger"`
	CreatedOn         time.Time `json:"created_on"`
	CompletedOn       time.Time `json:"completed_on"`
	RunNumber         int       `json:"run_number"`
	DurationInSeconds int       `json:"duration_in_seconds"`
	BuildSecondsUsed  int       `json:"build_seconds_used"`
	FirstSuccessful   bool      `json:"first_successful"`
	Expired           bool      `json:"expired"`
	HasVariables      bool      `json:"has_variables"`
	Links             struct {
		Self struct {
			Href string
		} `json:"self"`
		Steps struct {
			Href string
		}
	} `json:"links"`
}

var ExtractApiPipelinesMeta = core.SubTaskMeta{
	Name:             "extractApiPipelines",
	EntryPoint:       ExtractApiPipelines,
	EnabledByDefault: true,
	Description:      "Extract raw pipelines data into tool layer table BitbucketPipeline",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ExtractApiPipelines(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			// create bitbucket commit
			bitbucketApiPipeline := &BitbucketApiPipeline{}
			err := json.Unmarshal(row.Data, bitbucketApiPipeline)
			if err != nil {
				return nil, err
			}

			bitbucketPipeline := &models.BitbucketPipeline{
				ConnectionId:      data.Options.ConnectionId,
				BitbucketId:       bitbucketApiPipeline.Uuid,
				WebUrl:            bitbucketApiPipeline.Links.Self.Href,
				Status:            bitbucketApiPipeline.State.Name,
				Result:            bitbucketApiPipeline.State.Result.Name,
				RefName:           bitbucketApiPipeline.Target.RefName,
				DurationInSeconds: bitbucketApiPipeline.DurationInSeconds,
			}
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0, 2)
			results = append(results, bitbucketPipeline)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
