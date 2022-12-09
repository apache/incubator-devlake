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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/models/domainlayer/devops"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_REPOSITORIES_TABLE = "github_api_repositories"

type ApiRepoResponse GithubApiRepo

type GithubApiRepo struct {
	Name        string `json:"name"`
	GithubId    int    `json:"id"`
	HTMLUrl     string `json:"html_url"`
	Language    string `json:"language"`
	Description string `json:"description"`
	Owner       *GithubAccountResponse
	Parent      *GithubApiRepo      `json:"parent"`
	CreatedAt   helper.Iso8601Time  `json:"created_at"`
	UpdatedAt   *helper.Iso8601Time `json:"updated_at"`
	CloneUrl    string              `json:"clone_url"`
}

var ConvertRepoMeta = core.SubTaskMeta{
	Name:             "convertRepo",
	EntryPoint:       ConvertRepo,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_repos into  domain layer table repos and boards",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE, core.DOMAIN_TYPE_TICKET, core.DOMAIN_TYPE_CICD, core.DOMAIN_TYPE_CODE_REVIEW, core.DOMAIN_TYPE_CROSS},
}

func ConvertRepo(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	cursor, err := db.Cursor(
		dal.From(&models.GithubRepo{}),
		dal.Where("github_id = ?", repoId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	repoIdGen := didgen.NewDomainIdGenerator(&models.GithubRepo{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(models.GithubRepo{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_REPOSITORIES_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			repository := inputRow.(*models.GithubRepo)
			domainRepository := &code.Repo{
				DomainEntity: domainlayer.DomainEntity{
					Id: repoIdGen.Generate(data.Options.ConnectionId, repository.GithubId),
				},
				Name:        fmt.Sprintf("%s/%s", repository.OwnerLogin, repository.Name),
				Url:         repository.HTMLUrl,
				Description: repository.Description,
				ForkedFrom:  repository.ParentHTMLUrl,
				Language:    repository.Language,
				CreatedDate: &repository.CreatedDate,
				UpdatedDate: repository.UpdatedDate,
			}

			domainBoard := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: repoIdGen.Generate(data.Options.ConnectionId, repository.GithubId),
				},
				Name:        fmt.Sprintf("%s/%s", repository.OwnerLogin, repository.Name),
				Url:         fmt.Sprintf("%s/%s", repository.HTMLUrl, "issues"),
				Description: repository.Description,
				CreatedDate: &repository.CreatedDate,
			}

			domainBoardRepo := &crossdomain.BoardRepo{
				BoardId: repoIdGen.Generate(data.Options.ConnectionId, repository.GithubId),
				RepoId:  repoIdGen.Generate(data.Options.ConnectionId, repository.GithubId),
			}

			domainCicdScope := &devops.CicdScope{
				DomainEntity: domainlayer.DomainEntity{
					Id: repoIdGen.Generate(data.Options.ConnectionId, repository.GithubId),
				},
				Name:        repository.Name,
				Url:         repository.HTMLUrl,
				Description: repository.Description,
				CreatedDate: &repository.CreatedDate,
				UpdatedDate: repository.UpdatedDate,
			}

			return []interface{}{
				domainRepository,
				domainBoard,
				domainBoardRepo,
				domainCicdScope,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
