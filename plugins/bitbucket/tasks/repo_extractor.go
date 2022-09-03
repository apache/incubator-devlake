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
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"time"

	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractApiRepoMeta = core.SubTaskMeta{
	Name:        "extractApiRepo",
	EntryPoint:  ExtractApiRepositories,
	Required:    true,
	Description: "Extract raw Repositories data into tool layer table bitbucket_repos",
	DomainTypes: []string{core.DOMAIN_TYPE_CODE},
}

type ApiRepoResponse BitbucketApiRepo

type BitbucketApiRepo struct {
	BitbucketId string
	Scm         string `json:"scm"`
	HasWiki     bool   `json:"has_wiki"`
	Links       struct {
		Clone []struct {
			Href string
			Name string
		} `json:"clone"`
		Self struct {
			Href string
		} `json:"self"`

		Html struct {
			Href string
		} `json:"html"`
	} `json:"links"`
	Uuid        string `json:"uuid"`
	FullName    string `json:"full_name"`
	Language    string `json:"language"`
	Description string `json:"description"`
	Type        string `json:"type"`
	HasIssue    bool   `json:"has_issue"`
	ForkPolicy  string `json:"fork_policy"`
	Owner       models.BitbucketAccount
	CreatedAt   time.Time  `json:"created_on"`
	UpdatedAt   *time.Time `json:"updated_on"`
}

func ExtractApiRepositories(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*BitbucketTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: BitbucketApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			/*
				Table store raw data
			*/
			Table: RAW_REPOSITORIES_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			body := &ApiRepoResponse{}
			err := errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {
				return nil, err
			}
			if body.FullName == "" {
				return nil, errors.NotFound.New(fmt.Sprintf("repo %s/%s not found", data.Options.Owner, data.Options.Repo))
			}
			results := make([]interface{}, 0, 1)
			bitbucketRepository := &models.BitbucketRepo{
				ConnectionId: data.Options.ConnectionId,
				BitbucketId:  "repositories/" + body.FullName,
				Name:         body.FullName,
				HTMLUrl:      body.Links.Html.Href,
				Description:  body.Description,
				OwnerId:      body.Owner.AccountId,
				Language:     body.Language,
				CreatedDate:  body.CreatedAt,
				UpdatedDate:  body.UpdatedAt,
			}
			data.Repo = bitbucketRepository

			results = append(results, bitbucketRepository)

			parentTaskContext := taskCtx.TaskContext()
			if parentTaskContext != nil {
				parentTaskContext.GetData().(*BitbucketTaskData).Repo = bitbucketRepository
			}
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
