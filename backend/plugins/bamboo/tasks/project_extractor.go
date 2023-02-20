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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
)

var _ plugin.SubTaskEntryPoint = ExtractProject

func ExtractProject(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,

		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			res := &models.ApiBambooProject{}
			err := errors.Convert(json.Unmarshal(resData.Data, res))
			if err != nil {
				return nil, err
			}
			body := ConvertProject(res)
			body.ConnectionId = data.Options.ConnectionId
			return []interface{}{body}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractProjectMeta = plugin.SubTaskMeta{
	Name:             "ExtractProject",
	EntryPoint:       ExtractProject,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table bamboo_project",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

// Convert the API response to our DB model instance
func ConvertProject(bambooApiProject *models.ApiBambooProject) *models.BambooProject {
	bambooProject := &models.BambooProject{
		ProjectKey:  bambooApiProject.Key,
		Name:        bambooApiProject.Name,
		Description: bambooApiProject.Description,
		Href:        bambooApiProject.Link.Href,
		Rel:         bambooApiProject.Link.Rel,
	}
	return bambooProject
}

//type ApiProject struct {
//	ProjectKey         string `json:"key"`
//	Name        string `json:"name"`
//	Description string `json:"description"`
//	Link        Link   `json:"link"`
//}
//type ApiActions struct {
//	Size       int `json:"size"`
//	StartIndex int `json:"start-index"`
//	MaxResult  int `json:"max-result"`
//}
//type ApiStages struct {
//	Size       int `json:"size"`
//	StartIndex int `json:"start-index"`
//	MaxResult  int `json:"max-result"`
//}
//type ApiBranches struct {
//	Size       int `json:"size"`
//	StartIndex int `json:"start-index"`
//	MaxResult  int `json:"max-result"`
//}
//type ApiPlanKey struct {
//	ProjectKey string `json:"key"`
//}
//type ApiPlan struct {
//	Expand                    string `json:"expand"`
//	ProjectKey                string `json:"projectKey"`
//	ProjectName               string `json:"projectName"`
//	Description               string `json:"description"`
//	ShortName                 string `json:"shortName"`
//	BuildName                 string `json:"buildName"`
//	ShortKey                  string `json:"shortKey"`
//	Type                      string `json:"type"`
//	Enabled                   bool   `json:"enabled"`
//	ApiLink                   `json:"link"`
//	IsFavourite               bool     `json:"isFavourite"`
//	IsActive                  bool     `json:"isActive"`
//	IsBuilding                bool     `json:"isBuilding"`
//	AverageBuildTimeInSeconds int      `json:"averageBuildTimeInSeconds"`
//	Actions                   Actions  `json:"actions"`
//	Stages                    Stages   `json:"stages"`
//	Branches                  Branches `json:"branches"`
//	ProjectKey                       string   `json:"key"`
//	Name                      string   `json:"name"`
//	PlanKey                   PlanKey  `json:"planKey"`
//}
//type ApiPlans struct {
//	Size       int       `json:"size"`
//	Expand     string    `json:"expand"`
//	StartIndex int       `json:"start-index"`
//	MaxResult  int       `json:"max-result"`
//	Plan       []ApiPlan `json:"plan"`
//}
