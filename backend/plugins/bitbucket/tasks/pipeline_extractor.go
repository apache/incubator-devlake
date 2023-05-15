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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

type bitbucketApiCommit struct {
	//Type  string `json:"type"`
	Hash string `json:"hash"`
}

type bitbucketApiPipelineTarget struct {
	//Type     string `json:"type"`
	//RefType  string `json:"ref_type"`
	RefName string `json:"ref_name"`
	//Selector struct {
	//	Type string `json:"type"`
	//} `json:"selector"`
	Commit *bitbucketApiCommit `json:"commit"`
}

type BitbucketApiPipeline struct {
	Uuid string `json:"uuid"`
	//Type  string `json:"type"`
	State struct {
		Name string `json:"name"`
		//	Type   string `json:"type"`
		Result *struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"result"`
		Stage *struct {
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"stage"`
	} `json:"state"`
	Target            *bitbucketApiPipelineTarget `json:"target"`
	CreatedOn         *api.Iso8601Time            `json:"created_on"`
	CompletedOn       *api.Iso8601Time            `json:"completed_on"`
	DurationInSeconds uint64                      `json:"duration_in_seconds"`
	Links             struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

var ExtractApiPipelinesMeta = plugin.SubTaskMeta{
	Name:             "extractApiPipelines",
	EntryPoint:       ExtractApiPipelines,
	EnabledByDefault: true,
	Description:      "Extract raw pipelines data into tool layer table BitbucketPipeline",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ExtractApiPipelines(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_TABLE)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			// create bitbucket commit
			bitbucketApiPipeline := &BitbucketApiPipeline{}
			err := errors.Convert(json.Unmarshal(row.Data, bitbucketApiPipeline))
			if err != nil {
				return nil, err
			}

			bitbucketPipeline := &models.BitbucketPipeline{
				ConnectionId:        data.Options.ConnectionId,
				BitbucketId:         bitbucketApiPipeline.Uuid,
				WebUrl:              bitbucketApiPipeline.Links.Self.Href,
				Status:              bitbucketApiPipeline.State.Name,
				RefName:             bitbucketApiPipeline.Target.RefName,
				CommitSha:           bitbucketApiPipeline.Target.Commit.Hash,
				RepoId:              data.Options.FullName,
				DurationInSeconds:   bitbucketApiPipeline.DurationInSeconds,
				BitbucketCreatedOn:  api.Iso8601TimeToTime(bitbucketApiPipeline.CreatedOn),
				BitbucketCompleteOn: api.Iso8601TimeToTime(bitbucketApiPipeline.CompletedOn),
				Type:                data.RegexEnricher.ReturnNameIfMatched(devops.DEPLOYMENT, bitbucketApiPipeline.Target.RefName),
				Environment:         data.RegexEnricher.ReturnNameIfOmittedOrMatched(devops.PRODUCTION, bitbucketApiPipeline.Target.RefName),
			}
			if err != nil {
				return nil, err
			}
			if bitbucketApiPipeline.State.Result != nil {
				bitbucketPipeline.Result = bitbucketApiPipeline.State.Result.Name
			} else if bitbucketApiPipeline.State.Stage != nil {
				bitbucketPipeline.Result = bitbucketApiPipeline.State.Stage.Name
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
