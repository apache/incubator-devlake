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
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
)

var _ plugin.SubTaskEntryPoint = ExtractMetrics

func ExtractMetrics(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_METRICS_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,

		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			body := &models.SonarqubeApiMetrics{}
			err := errors.Convert(json.Unmarshal(resData.Data, body))
			if err != nil {
				return nil, err
			}
			metric := &models.SonarqubeMetrics{
				ConnectionId: data.Options.ConnectionId,
				Id:           body.Id,
				MetricsKey:   body.Key,
				Type:         body.Type,
				Name:         body.Name,
				Description:  body.Description,
				Domain:       body.Domain,
				Direction:    body.Direction,
				Qualitative:  body.Qualitative,
				Hidden:       body.Hidden,
				Custom:       body.Custom,
			}
			return []interface{}{metric}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractMetricsMeta = plugin.SubTaskMeta{
	Name:             "ExtractMetrics",
	EntryPoint:       ExtractMetrics,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table sonarqube_metrics",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}
