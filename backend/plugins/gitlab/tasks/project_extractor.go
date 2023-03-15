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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

// Convert the API response to our DB model instance
func ConvertProject(gitlabApiProject *GitlabApiProject) *models.GitlabProject {
	gitlabProject := &models.GitlabProject{
		GitlabId:          gitlabApiProject.GitlabId,
		Name:              gitlabApiProject.Name,
		Description:       gitlabApiProject.Description,
		DefaultBranch:     gitlabApiProject.DefaultBranch,
		CreatorId:         gitlabApiProject.CreatorId,
		PathWithNamespace: gitlabApiProject.PathWithNamespace,
		WebUrl:            gitlabApiProject.WebUrl,
		HttpUrlToRepo:     gitlabApiProject.HttpUrlToRepo,
		Visibility:        gitlabApiProject.Visibility,
		OpenIssuesCount:   gitlabApiProject.OpenIssuesCount,
		StarCount:         gitlabApiProject.StarCount,
		CreatedDate:       gitlabApiProject.CreatedAt.ToNullableTime(),
		UpdatedDate:       helper.Iso8601TimeToTime(gitlabApiProject.LastActivityAt),
	}
	if gitlabApiProject.ForkedFromProject != nil {
		gitlabProject.ForkedFromProjectId = gitlabApiProject.ForkedFromProject.GitlabId
		gitlabProject.ForkedFromProjectWebUrl = gitlabApiProject.ForkedFromProject.WebUrl
	}
	return gitlabProject
}
