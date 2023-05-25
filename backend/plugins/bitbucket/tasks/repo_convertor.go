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
	"io"
	"net/http"
	"path"
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	aha "github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

const RAW_REPOSITORIES_TABLE = "bitbucket_api_repositories"

var ConvertRepoMeta = plugin.SubTaskMeta{
	Name:             "convertRepo",
	EntryPoint:       ConvertRepo,
	EnabledByDefault: true,
	Description:      "Convert tool layer table bitbucket_repos into  domain layer table repos and boards",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
}

type ApiRepoResponse models.BitbucketApiRepo

func GetApiRepo(
	op *BitbucketOptions,
	apiClient aha.ApiClientAbstract,
) (*models.BitbucketApiRepo, errors.Error) {
	res, err := apiClient.Get(path.Join("repositories", op.FullName), nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.Default.New(fmt.Sprintf(
			"unexpected status code when requesting repo detail %d %s",
			res.StatusCode, res.Request.URL.String(),
		))
	}
	body, err := errors.Convert01(io.ReadAll(res.Body))
	if err != nil {
		return nil, err
	}
	apiRepo := new(models.BitbucketApiRepo)
	err = errors.Convert(json.Unmarshal(body, apiRepo))
	if err != nil {
		return nil, err
	}
	for _, u := range apiRepo.Links.Clone {
		if u.Name == "https" {
			return apiRepo, nil
		}
	}
	return nil, errors.Default.New("no clone url")
}

func ConvertRepo(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_REPOSITORIES_TABLE)
	db := taskCtx.GetDal()
	repoId := data.Options.FullName

	cursor, err := db.Cursor(
		dal.From(&models.BitbucketRepo{}),
		dal.Where("connection_id = ? AND bitbucket_id = ?", data.Options.ConnectionId, repoId),
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

			repoId := repoIdGen.Generate(data.Options.ConnectionId, repository.BitbucketId)

			domainRepository := &code.Repo{
				DomainEntity: domainlayer.DomainEntity{
					Id: repoId,
				},
				Name:        repository.BitbucketId,
				Url:         repository.HTMLUrl,
				Description: repository.Description,
				Language:    repository.Language,
				CreatedDate: repository.CreatedDate,
				UpdatedDate: repository.UpdatedDate,
			}

			domainBoard := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: repoId,
				},
				Name:        repository.BitbucketId,
				Url:         fmt.Sprintf("%s/%s", repository.HTMLUrl, "issues"),
				Description: repository.Description,
				CreatedDate: repository.CreatedDate,
			}

			domainCicdScope := &devops.CicdScope{
				DomainEntity: domainlayer.DomainEntity{
					Id: repoId,
				},
				Name:        repository.BitbucketId,
				Url:         fmt.Sprintf("%s/%s", repository.HTMLUrl, "issues"),
				Description: repository.Description,
				CreatedDate: repository.CreatedDate,
			}

			return []interface{}{
				domainRepository,
				domainBoard,
				domainCicdScope,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}
