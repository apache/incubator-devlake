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

var _ plugin.SubTaskEntryPoint = ExtractHotspots

func ExtractHotspots(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_HOTSPOTS_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			var res struct {
				Key                      string              `json:"key" gorm:"primaryKey"`
				RuleKey                  string              `json:"ruleKey"`
				Component                string              `json:"component" gorm:"index"`
				ProjectKey               string              `json:"project" gorm:"index"`
				Line                     int                 `json:"line"`
				Status                   string              `json:"status"`
				Message                  string              `json:"message"`
				Author                   string              `json:"author"`
				Assignee                 string              `json:"assignee"`
				SecurityCategory         string              `json:"securityCategory"`
				VulnerabilityProbability string              `json:"vulnerabilityProbability"`
				CreationDate             *helper.Iso8601Time `json:"creationDate"`
				UpdateDate               *helper.Iso8601Time `json:"updateDate"`
			}
			err := errors.Convert(json.Unmarshal(resData.Data, &res))
			if err != nil {
				return nil, err
			}
			body := &models.SonarqubeHotspot{
				ConnectionId:             data.Options.ConnectionId,
				HotspotKey:               res.Key,
				RuleKey:                  res.RuleKey,
				Component:                res.Component,
				ProjectKey:               res.ProjectKey,
				Line:                     res.Line,
				Status:                   res.Status,
				Message:                  res.Message,
				Author:                   res.Author,
				Assignee:                 res.Assignee,
				SecurityCategory:         res.SecurityCategory,
				VulnerabilityProbability: res.VulnerabilityProbability,
				CreationDate:             res.CreationDate,
				UpdateDate:               res.UpdateDate,
			}
			return []interface{}{body}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractHotspotsMeta = plugin.SubTaskMeta{
	Name:             "ExtractHotspots",
	EntryPoint:       ExtractHotspots,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table sonarqube_hotspots",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}
