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
	"net/http"
	"net/url"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

func GetTotalPagesFromResponse(res *http.Response, args *api.ApiCollectorArgs) (int, errors.Error) {
	body := &ZentaoPagination{}
	err := api.UnmarshalResponse(res, body)
	if err != nil {
		return 0, err
	}
	pages := body.Total / args.PageSize
	if body.Total%args.PageSize > 0 {
		pages++
	}
	return pages, nil
}

func getAccountId(account *models.ZentaoAccount) int64 {
	if account != nil {
		return account.ID
	}
	return 0
}

func getAccountName(account *models.ZentaoAccount) string {
	if account != nil {
		return account.Realname
	}
	return ""
}

// get the Priority string for zentao
func getPriority(pri int) string {
	switch pri {
	case 2:
		return "High"
	case 3:
		return "Middle"
	case 4:
		return "Low"
	default:
		if pri <= 1 {
			return "VeryHigh"
		}
		if pri >= 5 {
			return "VeryLow"
		}
	}
	return "Error"
}

func getOriginalProject(data *ZentaoTaskData) string {
	if data.Options.ProjectId != 0 {
		return data.ProjectName
	}
	if data.Options.ProductId != 0 {
		return data.ProductName
	}
	return ""
}

// getBugStatusMapping creates a map of original status values to bug issue standard status values
// based on the provided ZentaoTaskData. It returns the created map.
func getBugStatusMapping(data *ZentaoTaskData) map[string]string {
	stdStatusMappings := make(map[string]string)
	if data.Options.ScopeConfigs == nil {
		return stdStatusMappings
	}
	mapping := data.Options.ScopeConfigs.BugStatusMappings
	// Map original status values to standard status values
	for userStatus, stdStatus := range mapping {
		stdStatusMappings[userStatus] = strings.ToUpper(stdStatus)
	}
	return stdStatusMappings
}

// getStoryStatusMapping creates a map of original status values to story issue standard status values
// based on the provided ZentaoTaskData. It returns the created map.
func getStoryStatusMapping(data *ZentaoTaskData) map[string]string {
	stdStatusMappings := make(map[string]string)
	if data.Options.ScopeConfigs == nil {
		return stdStatusMappings
	}
	mapping := data.Options.ScopeConfigs.StoryStatusMappings
	// Map original status values to standard status values
	for userStatus, stdStatus := range mapping {
		stdStatusMappings[userStatus] = strings.ToUpper(stdStatus)
	}
	return stdStatusMappings
}

// getTaskStatusMapping creates a map of original status values to task issue standard status values
// based on the provided ZentaoTaskData. It returns the created map.
func getTaskStatusMapping(data *ZentaoTaskData) map[string]string {
	stdStatusMappings := make(map[string]string)
	if data.Options.ScopeConfigs == nil {
		return stdStatusMappings
	}
	mapping := data.Options.ScopeConfigs.TaskStatusMappings
	// Map original status values to standard status values
	for userStatus, stdStatus := range mapping {
		stdStatusMappings[userStatus] = strings.ToUpper(stdStatus)
	}
	return stdStatusMappings
}

// getStdTypeMappings creates a map of user type to standard type based on the provided ZentaoTaskData.
// It returns the created map.
func getStdTypeMappings(data *ZentaoTaskData) map[string]string {
	stdTypeMappings := make(map[string]string)
	if data.Options.ScopeConfigs == nil {
		return stdTypeMappings
	}
	mapping := data.Options.ScopeConfigs.TypeMappings
	// Map user types to standard types
	for userType, stdType := range mapping {
		stdTypeMappings[userType] = strings.ToUpper(stdType)
	}
	return stdTypeMappings
}

// parseRepoUrl parses a repository URL and returns the host, namespace, and repository name.
func parseRepoUrl(repoUrl string) (string, string, string, error) {
	parsedUrl, err := url.Parse(repoUrl)
	if err != nil {
		return "", "", "", err
	}

	host := parsedUrl.Host
	host = strings.TrimPrefix(host, "www.")
	pathParts := strings.Split(parsedUrl.Path, "/")
	if len(pathParts) < 3 {
		return "", "", "", fmt.Errorf("invaild RepoUrl: %s", repoUrl)
	}

	namespace := pathParts[1]
	repoName := pathParts[2]

	return host, namespace, repoName, nil
}

func ignoreHTTPStatus404(res *http.Response) errors.Error {
	if res.StatusCode == http.StatusUnauthorized {
		return errors.Unauthorized.New("authentication failed, please check your AccessToken")
	}
	if res.StatusCode == http.StatusNotFound {
		return api.ErrIgnoreAndContinue
	}
	return nil
}
