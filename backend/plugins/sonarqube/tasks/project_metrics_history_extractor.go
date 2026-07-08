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

var _ plugin.SubTaskEntryPoint = ExtractProjectMetricsHistory

type metricsHistoryResponse struct {
	Measures []struct {
		Metric  string `json:"metric"`
		History []struct {
			Date  string `json:"date"`
			Value string `json:"value"`
		} `json:"history"`
	} `json:"measures"`
}

func ExtractProjectMetricsHistory(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_METRICS_HISTORY_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			body := &metricsHistoryResponse{}
			err := errors.Convert(json.Unmarshal(resData.Data, body))
			if err != nil {
				return nil, err
			}

			var results []interface{}
			for _, measure := range body.Measures {
				for _, entry := range measure.History {
					if entry.Value == "" {
						continue
					}
					analysisDate, parseErr := time.Parse("2006-01-02T15:04:05-0700", entry.Date)
					if parseErr != nil {
						return nil, errors.Default.Wrap(errors.Convert(parseErr), "failed to parse analysis date")
					}
					results = append(results, &models.SonarqubeProjectMetricsHistory{
						ConnectionId: data.Options.ConnectionId,
						ProjectKey:   data.Options.ProjectKey,
						AnalysisDate: analysisDate,
						MetricKey:    measure.Metric,
						MetricValue:  entry.Value,
					})
				}
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractProjectMetricsHistoryMeta = plugin.SubTaskMeta{
	Name:             "ExtractProjectMetricsHistory",
	EntryPoint:       ExtractProjectMetricsHistory,
	EnabledByDefault: true,
	Description:      "Extract raw project metrics history into tool layer table",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}
