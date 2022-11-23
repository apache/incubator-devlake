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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"time"
)

type bitbucketApiDeploymentsResponse struct {
	Type   string `json:"type"`
	Number int    `json:"number"`
	UUID   string `json:"uuid"`
	Key    string `json:"key"`
	Step   struct {
		UUID string `json:"uuid"`
	} `json:"step"`
	Environment struct {
		UUID string `json:"uuid"`
	} `json:"environment"`
	Release struct {
		Type     string `json:"type"`
		UUID     string `json:"uuid"`
		Pipeline struct {
			UUID string `json:"uuid"`
			Type string `json:"type"`
		} `json:"pipeline"`
		Key    string `json:"key"`
		Name   string `json:"name"`
		URL    string `json:"url"`
		Commit struct {
			Type  string `json:"type"`
			Hash  string `json:"hash"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
			} `json:"links"`
		} `json:"commit"`
		CreatedOn *time.Time `json:"created_on"`
	} `json:"release"`
	State struct {
		Type   string `json:"type"`
		Name   string `json:"name"`
		URL    string `json:"url"`
		Status struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"status"`
		StartedOn   *time.Time `json:"started_on"`
		CompletedOn *time.Time `json:"completed_on"`
	} `json:"state"`
	LastUpdateTime *time.Time `json:"last_update_time"`
	Version        int        `json:"version"`
}

var ExtractApiDeploymentsMeta = core.SubTaskMeta{
	Name:             "extractApiDeployments",
	EntryPoint:       ExtractApiDeployments,
	EnabledByDefault: true,
	Description:      "Extract raw deployments data into tool layer table BitbucketDeployment",
	DomainTypes:      []string{core.DOMAIN_TYPE_CICD},
}

func ExtractApiDeployments(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_DEPLOYMENT_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			bitbucketApiDeployments := &bitbucketApiDeploymentsResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, bitbucketApiDeployments))
			if err != nil {
				return nil, err
			}

			bitbucketDeployment := &models.BitbucketDeployment{
				ConnectionId:   data.Options.ConnectionId,
				BitbucketId:    bitbucketApiDeployments.UUID,
				PipelineId:     bitbucketApiDeployments.Release.Pipeline.UUID,
				Type:           bitbucketApiDeployments.Type,
				Name:           bitbucketApiDeployments.Release.Name,
				Key:            bitbucketApiDeployments.Release.Key,
				WebUrl:         bitbucketApiDeployments.Release.URL,
				CommitSha:      bitbucketApiDeployments.Release.Commit.Hash,
				CommitUrl:      bitbucketApiDeployments.Release.Commit.Links.HTML.Href,
				Status:         bitbucketApiDeployments.State.Name,
				StateUrl:       bitbucketApiDeployments.State.URL,
				CreatedOn:      bitbucketApiDeployments.Release.CreatedOn,
				StartedOn:      bitbucketApiDeployments.State.StartedOn,
				CompletedOn:    bitbucketApiDeployments.State.CompletedOn,
				LastUpdateTime: bitbucketApiDeployments.LastUpdateTime,
			}
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0, 2)
			results = append(results, bitbucketDeployment)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
