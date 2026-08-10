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
	"net/http"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

const jiraSearchEndpointV2 = "api/2/search"
const jiraSearchEndpointV3 = "api/3/search/jql"

func isJiraCloudDeployment(deploymentType models.DeploymentType) bool {
	return strings.EqualFold(string(deploymentType), string(models.DeploymentCloud))
}

func getJiraSearchEndpoint(deploymentType models.DeploymentType) string {
	if isJiraCloudDeployment(deploymentType) {
		return jiraSearchEndpointV3
	}
	// Jira Server and Data Center continue using v2 search for compatibility.
	return jiraSearchEndpointV2
}

func buildJiraV3SearchRequestBody(jql string, reqData *api.RequestData) map[string]interface{} {
	body := map[string]interface{}{
		"jql":        jql,
		"maxResults": reqData.Pager.Size,
		"expand":     "changelog",
		"fields":     "*all",
	}
	if nextPageToken, ok := reqData.CustomData.(string); ok && nextPageToken != "" {
		body["nextPageToken"] = nextPageToken
	}
	return body
}

func GetTotalPagesFromResponse(res *http.Response, args *api.ApiCollectorArgs) (int, errors.Error) {
	body := &JiraPagination{}
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

func getStdStatus(statusKey string) string {
	if statusKey == "done" {
		return ticket.DONE
	} else if statusKey == "new" {
		return ticket.TODO
	} else {
		return ticket.IN_PROGRESS
	}
}
