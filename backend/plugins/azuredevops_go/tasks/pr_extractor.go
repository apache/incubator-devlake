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
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiPullRequestsMeta)
}

var ExtractApiPullRequestsMeta = plugin.SubTaskMeta{
	Name:             "extractApiPullRequests",
	EntryPoint:       ExtractApiPullRequests,
	EnabledByDefault: true,
	Description:      "Extract raw PullRequests data into tool layer table azuredevops_pull_requests",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS, plugin.DOMAIN_TYPE_CODE_REVIEW},
	DependencyTables: []string{RawPullRequestTable},
	ProductTables: []string{
		models.AzuredevopsPullRequest{}.TableName(),
		models.AzuredevopsPrLabel{}.TableName(),
	},
}

func ExtractApiPullRequests(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawPullRequestTable)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			rawL := &models.AzuredevopsApiPullRequest{}
			err := errors.Convert(json.Unmarshal(row.Data, rawL))
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0, 2)

			//If this is a pr, ignore
			adoApiPr, err := convertAzuredevopsPullRequest(rawL, data.Options.ConnectionId)
			if err != nil {
				return nil, err
			}
			for _, label := range rawL.Labels {
				results = append(results, &models.AzuredevopsPrLabel{
					ConnectionId:  data.Options.ConnectionId,
					PullRequestId: adoApiPr.AzuredevopsId,
					LabelName:     label.Name,
				})
			}

			results = append(results, adoApiPr)

			return results, nil
		},
	})

	if err != nil {
		return errors.Default.Wrap(err, "error initializing Azure DevOps PR extractor")
	}

	return extractor.Execute()
}
func convertAzuredevopsPullRequest(pull *models.AzuredevopsApiPullRequest, connId uint64) (*models.AzuredevopsPullRequest, errors.Error) {
	var state string
	if pull.Status == "abandoned" {
		state = "CLOSED"
	} else if pull.Status == "active" {
		state = "OPEN"
	} else if pull.Status == "completed" {
		state = "MERGED"
	}

	pr := &models.AzuredevopsPullRequest{
		ConnectionId:    connId,
		AzuredevopsId:   pull.PullRequestId,
		RepositoryId:    pull.Repository.Id,
		Description:     pull.Description,
		Status:          state,
		CreatedById:     pull.CreatedBy.Id,
		CreatedByName:   pull.CreatedBy.DisplayName,
		CreationDate:    common.Iso8601TimeToTime(pull.AzuredevopsCreationDate),
		ClosedDate:      common.Iso8601TimeToTime(pull.ClosedDate),
		SourceCommitSha: pull.LastMergeSourceCommit.CommitId,
		TargetCommitSha: pull.LastMergeTargetCommit.CommitId,
		MergeCommitSha:  pull.LastMergeCommit.CommitId,
		Url:             pull.Url,
		Type:            "",
		Title:           pull.Title,
		TargetRefName:   pull.TargetRefName,
		SourceRefName:   pull.SourceRefName,
		ForkRepoId:      "",
	}

	return pr, nil
}
