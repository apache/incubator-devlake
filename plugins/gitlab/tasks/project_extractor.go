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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ExtractProjectMeta = core.SubTaskMeta{
	Name:             "extractApiProject",
	EntryPoint:       ExtractApiProject,
	EnabledByDefault: true,
	Description:      "Extract raw project data into tool layer table GitlabProject",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE, core.DOMAIN_TYPE_TICKET},
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

func ExtractApiProject(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, errors.Error) {
			// create gitlab commit
			gitlabApiProject := &GitlabApiProject{}
			err := errors.Convert(json.Unmarshal(row.Data, gitlabApiProject))
			if err != nil {
				return nil, err
			}
			gitlabProject := convertProject(gitlabApiProject)
			gitlabProject.ConnectionId = data.Options.ConnectionId
			results := make([]interface{}, 0, 1)
			results = append(results, gitlabProject)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
