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

	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertProjectMeta = core.SubTaskMeta{
	Name:             "convertApiProject",
	EntryPoint:       ConvertApiProjects,
	EnabledByDefault: true,
	Description:      "Update domain layer Repo according to GitlabProject",
}

func ConvertApiProjects(taskCtx core.SubTaskContext) error {

	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)
	db := taskCtx.GetDb()

	//Find all piplines associated with the current projectid
	cursor, err := db.Model(&models.GitlabProject{}).Where("gitlab_id=?", data.Options.ProjectId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabProject{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			gitlabProject := inputRow.(*models.GitlabProject)

			domainRepository := convertToRepositoryModel(gitlabProject)

			return []interface{}{
				domainRepository,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

// Convert the API response to our DB model instance
func convertProject(gitlabApiProject *GitlabApiProject) *models.GitlabProject {
	gitlabProject := &models.GitlabProject{
		GitlabId:          gitlabApiProject.GitlabId,
		Name:              gitlabApiProject.Name,
		Description:       gitlabApiProject.Description,
		DefaultBranch:     gitlabApiProject.DefaultBranch,
		CreatorId:         gitlabApiProject.CreatorId,
		PathWithNamespace: gitlabApiProject.PathWithNamespace,
		WebUrl:            gitlabApiProject.WebUrl,
		Visibility:        gitlabApiProject.Visibility,
		OpenIssuesCount:   gitlabApiProject.OpenIssuesCount,
		StarCount:         gitlabApiProject.StarCount,
		CreatedDate:       gitlabApiProject.CreatedAt.ToTime(),
		UpdatedDate:       helper.Iso8601TimeToTime(gitlabApiProject.LastActivityAt),
	}
	if gitlabApiProject.ForkedFromProject != nil {
		gitlabProject.ForkedFromProjectId = gitlabApiProject.ForkedFromProject.GitlabId
		gitlabProject.ForkedFromProjectWebUrl = gitlabApiProject.ForkedFromProject.WebUrl
	}
	return gitlabProject
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
