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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertProjectMeta = core.SubTaskMeta{
	Name:             "convertApiProject",
	EntryPoint:       ConvertApiProjects,
	EnabledByDefault: true,
	Description:      "Add domain layer Repo according to GitlabProject",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE, core.DOMAIN_TYPE_TICKET},
}

func ConvertApiProjects(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.From(&models.GitlabProject{}),
		dal.Where("gitlab_id=?", data.Options.ProjectId),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabProject{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			gitlabProject := inputRow.(*models.GitlabProject)

			domainRepository := convertToRepositoryModel(gitlabProject)
			domainBoard := convertToBoardModel(gitlabProject)
			domainBoardRepo := convertToBoardRepoModel(gitlabProject)
			return []interface{}{
				domainRepository,
				domainBoard,
				domainBoardRepo,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

func convertToRepositoryModel(project *models.GitlabProject) *code.Repo {
	domainRepository := &code.Repo{
		DomainEntity: domainlayer.DomainEntity{
			Id: didgen.NewDomainIdGenerator(project).Generate(project.ConnectionId, project.GitlabId),
		},
		Name:        project.Name,
		Url:         project.WebUrl,
		Description: project.Description,
		ForkedFrom:  project.ForkedFromProjectWebUrl,
		CreatedDate: project.CreatedDate,
		UpdatedDate: project.UpdatedDate,
	}
	return domainRepository
}

func convertToBoardModel(project *models.GitlabProject) *ticket.Board {
	domainBoard := &ticket.Board{
		DomainEntity: domainlayer.DomainEntity{
			Id: didgen.NewDomainIdGenerator(project).Generate(project.ConnectionId, project.GitlabId),
		},
		Name:        project.Name,
		Url:         project.WebUrl,
		Description: project.Description,
		CreatedDate: &project.CreatedDate,
	}
	return domainBoard
}

func convertToBoardRepoModel(project *models.GitlabProject) *crossdomain.BoardRepo {
	domainBoardRepo := &crossdomain.BoardRepo{
		BoardId: didgen.NewDomainIdGenerator(project).Generate(project.ConnectionId, project.GitlabId),
		RepoId:  didgen.NewDomainIdGenerator(project).Generate(project.ConnectionId, project.GitlabId),
	}
	return domainBoardRepo
}
