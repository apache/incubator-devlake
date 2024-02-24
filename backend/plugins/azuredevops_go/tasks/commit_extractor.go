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
	"time"
)

func init() {
	RegisterSubtaskMeta(&ExtractApiCommitsMeta)
}

var ExtractApiCommitsMeta = plugin.SubTaskMeta{
	Name:             "extractApiCommits",
	EntryPoint:       ExtractApiCommits,
	EnabledByDefault: false,
	Description:      "Extract raw commit data into tool layer table",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
	DependencyTables: []string{RawCommitTable},
	ProductTables: []string{
		models.AzuredevopsCommit{}.TableName(),
		models.AzuredevopsRepoCommit{}.TableName()},
}

type CommitApiResponse struct {
	CommitId string `json:"commitId"`
	Author   struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	} `json:"author"`
	Committer struct {
		Name  string    `json:"name"`
		Email string    `json:"email"`
		Date  time.Time `json:"date"`
	} `json:"committer"`
	Comment      string `json:"comment"`
	ChangeCounts struct {
		Add    int `json:"Add"`
		Edit   int `json:"Edit"`
		Delete int `json:"Delete"`
	} `json:"changeCounts"`
	Url       string `json:"url"`
	RemoteUrl string `json:"remoteUrl"`
}

func ExtractApiCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RawCommitTable)

	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			results := make([]interface{}, 0, 2)

			apiCommit := &CommitApiResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, apiCommit))
			if err != nil {
				return nil, err
			}
			azuredevopsCommit, err := ConvertCommit(apiCommit)
			if err != nil {
				return nil, err
			}

			// create project/commits relationship
			repoCommit := &models.AzuredevopsRepoCommit{RepositoryId: data.Options.RepositoryId}
			repoCommit.CommitSha = azuredevopsCommit.Sha
			repoCommit.ConnectionId = data.Options.ConnectionId

			results = append(results, azuredevopsCommit)
			results = append(results, repoCommit)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func ConvertCommit(commit *CommitApiResponse) (*models.AzuredevopsCommit, errors.Error) {
	azureDevOpsCommit := &models.AzuredevopsCommit{
		NoPKModel:      common.NoPKModel{},
		Sha:            commit.CommitId,
		Message:        commit.Comment,
		AuthorName:     commit.Author.Name,
		AuthorEmail:    commit.Author.Email,
		AuthoredDate:   &commit.Author.Date,
		CommitterName:  commit.Committer.Name,
		CommitterEmail: commit.Committer.Email,
		CommittedDate:  &commit.Committer.Date,
		WebUrl:         commit.Url,
		Additions:      commit.ChangeCounts.Add,
		Deletions:      commit.ChangeCounts.Delete,
		Edit:           commit.ChangeCounts.Edit,
	}
	return azureDevOpsCommit, nil
}
