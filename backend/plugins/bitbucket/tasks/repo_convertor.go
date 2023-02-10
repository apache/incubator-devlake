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
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
	"reflect"
	"time"
)

const RAW_REPOSITORIES_TABLE = "bitbucket_api_repositories"

var ConvertRepoMeta = plugin.SubTaskMeta{
	Name:             "convertRepo",
	EntryPoint:       ConvertRepo,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bitbucket_repos into  domain layer table repos and boards",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

type ApiRepoResponse BitbucketApiRepo

type BitbucketApiRepo struct {
	Scm         string                  `json:"scm"`
	HasWiki     bool                    `json:"has_wiki"`
	Uuid        string                  `json:"uuid"`
	Name        string                  `json:"name"`
	FullName    string                  `json:"full_name"`
	Language    string                  `json:"language"`
	Description string                  `json:"description"`
	Type        string                  `json:"type"`
	HasIssue    bool                    `json:"has_issue"`
	ForkPolicy  string                  `json:"fork_policy"`
	Owner       models.BitbucketAccount `json:"owner"`
	CreatedAt   *time.Time              `json:"created_on"`
	UpdatedAt   *time.Time              `json:"updated_on"`
	Links       struct {
		Clone []struct {
			Href string `json:"href"`
			Name string `json:"name"`
		} `json:"clone"`
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		Html struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
}

func ConvertRepo(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_REPOSITORIES_TABLE)
	db := taskCtx.GetDal()
	repoId := data.Options.FullName

	cursor, err := db.Cursor(
		dal.From(&models.BitbucketRepo{}),
		dal.Where("bitbucket_id = ?", repoId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	repoIdGen := didgen.NewDomainIdGenerator(&models.BitbucketRepo{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.BitbucketRepo{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			repository := inputRow.(*models.BitbucketRepo)
			domainRepository := &code.Repo{
				DomainEntity: domainlayer.DomainEntity{
					Id: repoIdGen.Generate(data.Options.ConnectionId, repository.BitbucketId),
				},
				Name:        repository.Name,
				Url:         repository.HTMLUrl,
				Description: repository.Description,
				Language:    repository.Language,
				CreatedDate: repository.CreatedDate,
				UpdatedDate: repository.UpdatedDate,
			}

			domainBoard := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: repoIdGen.Generate(data.Options.ConnectionId, repository.BitbucketId),
				},
				Name:        repository.Name,
				Url:         fmt.Sprintf("%s/%s", repository.HTMLUrl, "issues"),
				Description: repository.Description,
				CreatedDate: repository.CreatedDate,
			}

			return []interface{}{
				domainRepository,
				domainBoard,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
