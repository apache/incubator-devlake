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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
)

var _ plugin.SubTaskEntryPoint = ExtractProjects

func ExtractProjects(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECTS_TABLE)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			body := &SonarqubeRawProject{}
			err := errors.Convert(json.Unmarshal(resData.Data, body))
			if err != nil {
				return nil, err
			}

			sonarqubeProject := &models.SonarqubeProject{
				ConnectionId:     data.Options.ConnectionId,
				ProjectKey:       body.Key,
				Name:             body.Name,
				Qualifier:        body.Qualifier,
				Visibility:       body.Visibility,
				LastAnalysisDate: body.LastAnalysisDate,
				Revision:         body.Revision,
			}

			return []interface{}{sonarqubeProject}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractProjectsMeta = plugin.SubTaskMeta{
	Name:             "ExtractProjects",
	EntryPoint:       ExtractProjects,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table sonarqube_projects",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_SECURITY_TESTING},
}

type SonarqubeRawProject struct {
	Key              string           `json:"key"`
	Name             string           `json:"name"`
	Qualifier        string           `json:"qualifier"`
	Visibility       string           `json:"visibility"`
	LastAnalysisDate *api.Iso8601Time `json:"lastAnalysisDate"`
	Revision         string           `json:"revision"`
}
