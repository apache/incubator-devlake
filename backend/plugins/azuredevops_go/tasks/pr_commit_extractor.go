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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiPullRequestCommitsMeta)
}

var ExtractApiPullRequestCommitsMeta = plugin.SubTaskMeta{
	Name:             "extractApiPullRequestCommits",
	EntryPoint:       ExtractApiPullRequestCommits,
	EnabledByDefault: true,
	Description:      "Extract raw pull requests commit data into tool layer table AzuredevopsPullRequestCommit and AzuredevopsRepoCommit",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS, plugin.DOMAIN_TYPE_CODE_REVIEW},
	DependencyTables: []string{RawPrCommitTable},
	ProductTables: []string{
		models.AzuredevopsPrCommit{}.TableName(),
		models.AzuredevopsRepoCommit{}.TableName(),
	},
}

func ExtractApiPullRequestCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawPrCommitTable)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			apiResponse := &models.AzuredevopsApiPrCommit{}
			err := errors.Convert(json.Unmarshal(row.Data, apiResponse))
			if err != nil {
				return nil, err
			}

			input := &SimplePr{}
			err = errors.Convert(json.Unmarshal(row.Input, &input))
			if err != nil {
				return nil, err
			}

			azuredevopsCommit := &models.AzuredevopsCommit{
				Sha:            apiResponse.CommitId,
				Message:        apiResponse.Comment,
				AuthorName:     apiResponse.Author.Name,
				AuthorEmail:    apiResponse.Author.Email,
				AuthoredDate:   &apiResponse.Author.Date,
				CommitterName:  apiResponse.Committer.Name,
				CommitterEmail: apiResponse.Committer.Email,
				CommittedDate:  &apiResponse.Committer.Date,
				WebUrl:         apiResponse.Url,
			}

			azuredevopsPRCommit := &models.AzuredevopsPrCommit{
				CommitSha:     apiResponse.CommitId,
				ConnectionId:  data.Options.ConnectionId,
				PullRequestId: input.AzuredevopsId,
				AuthorDate:    &apiResponse.Author.Date,
				AuthorName:    apiResponse.Author.Name,
				AuthorEmail:   apiResponse.Author.Email,
			}
			azuredevopsRepoCommit := &models.AzuredevopsRepoCommit{
				ConnectionId: data.Options.ConnectionId,
				RepositoryId: data.Options.RepositoryId,
				CommitSha:    apiResponse.CommitId,
			}

			results := make([]interface{}, 0, 3)
			results = append(results, azuredevopsCommit, azuredevopsPRCommit, azuredevopsRepoCommit)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
