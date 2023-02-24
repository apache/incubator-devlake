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
	"encoding/hex"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"hash"
	"net/http"
	"unicode"
)

func CreateRawDataSubTaskArgs(taskCtx plugin.SubTaskContext, rawTable string) (*api.RawDataSubTaskArgs, *SonarqubeTaskData) {
	data := taskCtx.GetData().(*SonarqubeTaskData)
	filteredData := *data
	filteredData.Options = &SonarqubeOptions{}
	*filteredData.Options = *data.Options
	var params = SonarqubeApiParams{
		ConnectionId: data.Options.ConnectionId,
		ProjectKey:   data.Options.ProjectKey,
	}
	rawDataSubTaskArgs := &api.RawDataSubTaskArgs{
		Ctx:    taskCtx,
		Params: params,
		Table:  rawTable,
	}
	return rawDataSubTaskArgs, &filteredData
}

func GetTotalPagesFromResponse(res *http.Response, args *api.ApiCollectorArgs) (int, errors.Error) {
	body := &SonarqubePagination{}
	err := api.UnmarshalResponse(res, body)
	if err != nil {
		return 0, err
	}
	pages := body.Paging.Total / args.PageSize
	if body.Paging.Total%args.PageSize > 0 {
		pages++
	}
	return pages, nil
}

type SonarqubePagination struct {
	Paging Paging `json:"paging"`
}
type Paging struct {
	PageIndex int `json:"pageIndex"`
	PageSize  int `json:"pageSize"`
	Total     int `json:"total"`
}

type SonarqubeApiProject struct {
	ProjectKey       string           `json:"key"`
	Name             string           `json:"name"`
	Qualifier        string           `json:"qualifier"`
	Visibility       string           `json:"visibility"`
	LastAnalysisDate *api.Iso8601Time `json:"lastAnalysisDate"`
	Revision         string           `json:"revision"`
}

// Convert the API response to our DB model instance
func ConvertProject(sonarqubeApiProject *SonarqubeApiProject) *models.SonarqubeProject {
	sonarqubeProject := &models.SonarqubeProject{
		ProjectKey:       sonarqubeApiProject.ProjectKey,
		Name:             sonarqubeApiProject.Name,
		Qualifier:        sonarqubeApiProject.Qualifier,
		Visibility:       sonarqubeApiProject.Visibility,
		LastAnalysisDate: sonarqubeApiProject.LastAnalysisDate,
		Revision:         sonarqubeApiProject.Revision,
	}
	return sonarqubeProject
}

func generateId(hashCodeBlock hash.Hash, entity *models.SonarqubeIssueCodeBlock) {
	hashCodeBlock.Write([]byte(fmt.Sprintf("%s-%s-%d-%d-%d-%d-%s", entity.IssueKey, entity.Component, entity.StartLine, entity.EndLine, entity.StartOffset, entity.EndOffset, entity.Msg)))
	entity.Id = hex.EncodeToString(hashCodeBlock.Sum(nil))
}

func convertTimeToMinutes(timeStr string) int {
	days := 0
	hours := 0
	minutes := 0

	var currentNum int
	var currentUnit string

	for i := 0; i < len(timeStr); i++ {
		c := timeStr[i]

		if unicode.IsDigit(rune(c)) {
			currentNum = currentNum*10 + int(c-'0')
		} else {
			currentUnit += string(c)
			if currentUnit == "d" {
				days = currentNum
			} else if currentUnit == "h" {
				hours = currentNum
			} else if currentUnit == "min" {
				minutes = currentNum
			} else {
				continue
			}
			currentNum = 0
			currentUnit = ""
		}
	}

	totalMinutes := days*24*60 + hours*60 + minutes
	return totalMinutes
}
