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
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertProjectMeta)
}

const RAW_PROJECT_TABLE = "gitlab_api_project"

type GitlabApiProject struct {
	GitlabId          int    `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	DefaultBranch     string `json:"default_branch"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebUrl            string `json:"web_url"`
	CreatorId         int
	Visibility        string              `json:"visibility"`
	OpenIssuesCount   int                 `json:"open_issues_count"`
	StarCount         int                 `json:"star_count"`
	ForkedFromProject *GitlabApiProject   `json:"forked_from_project"`
	CreatedAt         common.Iso8601Time  `json:"created_at"`
	LastActivityAt    *common.Iso8601Time `json:"last_activity_at"`
	HttpUrlToRepo     string              `json:"http_url_to_repo"`
	Archived          bool                `json:"archived"`
}

var ConvertProjectMeta = plugin.SubTaskMeta{
	Name:             "Convert Projects",
	EntryPoint:       ConvertApiProjects,
	EnabledByDefault: true,
	Description:      "Add domain layer Repo according to GitlabProject",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE, plugin.DOMAIN_TYPE_TICKET},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertAccountsMeta},
}

func ConvertApiProjects(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.From(&models.GitlabProject{}),
		dal.Where("gitlab_id=? and connection_id = ?", data.Options.ProjectId, data.Options.ConnectionId),
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
			domainCicdScope := convertToCicdScopeModel(gitlabProject)
			domainBoard := convertToBoardModel(gitlabProject)
			domainBoardRepo := convertToBoardRepoModel(gitlabProject)
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

func convertToRepositoryModel(project *models.GitlabProject) *code.Repo {
	domainRepository := &code.Repo{
		DomainEntity: domainlayer.DomainEntity{
			Id: didgen.NewDomainIdGenerator(project).Generate(project.ConnectionId, project.GitlabId),
		},
		Name:        project.PathWithNamespace,
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
		Name:        project.PathWithNamespace,
		Url:         project.WebUrl,
		Description: project.Description,
		CreatedDate: project.CreatedDate,
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

func convertToCicdScopeModel(project *models.GitlabProject) *devops.CicdScope {
	domainCicdScope := &devops.CicdScope{
		DomainEntity: domainlayer.DomainEntity{
			Id: didgen.NewDomainIdGenerator(project).Generate(project.ConnectionId, project.GitlabId),
		},
		Name:        project.PathWithNamespace,
		Url:         project.WebUrl,
		Description: project.Description,
		CreatedDate: project.CreatedDate,
		UpdatedDate: project.UpdatedDate,
	}
	return domainCicdScope
}
