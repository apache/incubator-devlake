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

package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/codequality"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	aha "github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
)

func MakeDataSourcePipelinePlanV200(subtaskMetas []plugin.SubTaskMeta, connectionId uint64, bpScopes []*plugin.BlueprintScopeV200, syncPolicy *plugin.BlueprintSyncPolicy) (plugin.PipelinePlan, []plugin.Scope, errors.Error) {
	plan := make(plugin.PipelinePlan, len(bpScopes))
	plan, err := makeDataSourcePipelinePlanV200(subtaskMetas, plan, bpScopes, connectionId, syncPolicy)
	if err != nil {
		return nil, nil, err
	}
	scopes, err := makeScopesV200(bpScopes, connectionId)
	if err != nil {
		return nil, nil, err
	}

	return plan, scopes, nil
}

func makeDataSourcePipelinePlanV200(
	subtaskMetas []plugin.SubTaskMeta,
	plan plugin.PipelinePlan,
	bpScopes []*plugin.BlueprintScopeV200,
	connectionId uint64,
	syncPolicy *plugin.BlueprintSyncPolicy,
) (plugin.PipelinePlan, errors.Error) {
	for i, bpScope := range bpScopes {
		stage := plan[i]
		if stage == nil {
			stage = plugin.PipelineStage{}
		}
		// construct task options for Sonarqube
		options := make(map[string]interface{})
		options["connectionId"] = connectionId
		options["projectKey"] = bpScope.Id

		subtasks, err := helper.MakePipelinePlanSubtasks(subtaskMetas, bpScope.Entities)
		if err != nil {
			return nil, err
		}
		stage = append(stage, &plugin.PipelineTask{
			Plugin:   "sonarqube",
			Subtasks: subtasks,
			Options:  options,
		})
		plan[i] = stage
	}

	return plan, nil
}

func makeScopesV200(bpScopes []*plugin.BlueprintScopeV200, connectionId uint64) ([]plugin.Scope, errors.Error) {
	scopes := make([]plugin.Scope, 0)
	for _, bpScope := range bpScopes {
		sonarqubeProject := &models.SonarqubeProject{}
		// get repo from db
		err := basicRes.GetDal().First(sonarqubeProject,
			dal.Where("connection_id = ? and project_key = ?",
				connectionId, bpScope.Id))
		if err != nil {
			return nil, errors.Default.Wrap(err, fmt.Sprintf("fail to find sonarqube project %s", bpScope.Id))
		}
		// add board to scopes
		if utils.StringsContains(bpScope.Entities, plugin.DOMAIN_TYPE_CODE_QUALITY) {
			stProject := &codequality.CqProject{
				DomainEntity: domainlayer.DomainEntity{
					Id: didgen.NewDomainIdGenerator(&models.SonarqubeProject{}).Generate(sonarqubeProject.ConnectionId, sonarqubeProject.ProjectKey),
				},
				Name: sonarqubeProject.Name,
			}
			scopes = append(scopes, stProject)
		}
	}
	return scopes, nil
}

func GetApiProject(
	projectKey string,
	apiClient aha.ApiClientAbstract,
) (*models.SonarqubeApiProject, errors.Error) {
	var resData struct {
		Data []models.SonarqubeApiProject `json:"components"`
	}
	query := url.Values{}
	query.Set("q", projectKey)
	res, err := apiClient.Get("projects/search", query, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code when requesting project detail from %s", res.Request.URL.String()))
	}
	err = helper.UnmarshalResponse(res, &resData)
	if err != nil {
		return nil, err
	}
	if len(resData.Data) > 0 {
		return &resData.Data[0], nil
	}
	return nil, errors.BadInput.New(fmt.Sprintf("Cannot find project: %s", projectKey))
}
