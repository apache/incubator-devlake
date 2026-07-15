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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
)

var _ plugin.SubTaskEntryPoint = ExtractProjectAnalyses

type projectAnalysisResponse struct {
	Key            string `json:"key"`
	Date           string `json:"date"`
	ProjectVersion string `json:"projectVersion"`
	Revision       string `json:"revision"`
	BuildString    string `json:"buildString"`
	DetectedCI     string `json:"detectedCI"`
}

func ExtractProjectAnalyses(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_ANALYSES_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			body := &projectAnalysisResponse{}
			err := errors.Convert(json.Unmarshal(resData.Data, body))
			if err != nil {
				return nil, err
			}

			analysisDate, parseErr := time.Parse("2006-01-02T15:04:05-0700", body.Date)
			if parseErr != nil {
				return nil, errors.Default.Wrap(errors.Convert(parseErr), "failed to parse analysis date")
			}

			analysis := &models.SonarqubeProjectAnalysis{
				ConnectionId:   data.Options.ConnectionId,
				ProjectKey:     data.Options.ProjectKey,
				AnalysisKey:    body.Key,
				AnalysisDate:   analysisDate,
				ProjectVersion: body.ProjectVersion,
				Revision:       body.Revision,
				BuildString:    body.BuildString,
				DetectedCI:     body.DetectedCI,
			}
			return []interface{}{analysis}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractProjectAnalysesMeta = plugin.SubTaskMeta{
	Name:             "ExtractProjectAnalyses",
	EntryPoint:       ExtractProjectAnalyses,
	EnabledByDefault: true,
	Description:      "Extract raw project analyses data into tool layer table",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}
