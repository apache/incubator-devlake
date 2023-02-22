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
	"strconv"
)

var _ plugin.SubTaskEntryPoint = ExtractFilemetrics

func ExtractFilemetrics(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_FILEMETRICS_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,

		Extract: func(resData *helper.RawData) ([]interface{}, errors.Error) {
			body := &fileMetricsResponse{}
			err := errors.Convert(json.Unmarshal(resData.Data, body))
			if err != nil {
				return nil, err
			}
			fileMetrics := &models.SonarqubeFileMetrics{
				ConnectionId:   data.Options.ConnectionId,
				FileMetricsKey: body.Key,
				FileName:       body.Name,
				FilePath:       body.Path,
				FileLanguage:   body.Language,
				ProjectKey:     data.Options.ProjectKey,
			}
			alphabetMap := map[string]string{
				"1.0": "A",
				"2.0": "B",
				"3.0": "C",
				"4.0": "D",
				"5.0": "E",
			}
			for _, v := range body.Measures {
				switch v.Metric {
				case "sqale_index":
					fileMetrics.SqaleIndex, err = errors.Convert01(strconv.Atoi(v.Value))
					if err != nil {
						return nil, err
					}
				case "sqale_rating":
					fileMetrics.SqaleRating, err = errors.Convert01(strconv.ParseFloat(v.Value, 32))
					if err != nil {
						return nil, err
					}
				case "reliability_rating":
					fileMetrics.ReliabilityRating = alphabetMap[v.Value]
				case "security_rating":
					fileMetrics.SecurityRating = alphabetMap[v.Value]
				case "security_review_rating":
					fileMetrics.SecurityReviewRating = alphabetMap[v.Value]
				case "ncloc":
					fileMetrics.Ncloc, err = errors.Convert01(strconv.Atoi(v.Value))
					if err != nil {
						return nil, err
					}
				case "duplicated_blocks":
					fileMetrics.DuplicatedBlocks, err = errors.Convert01(strconv.Atoi(v.Value))
					if err != nil {
						return nil, err
					}

				case "duplicated_lines_density":
					fileMetrics.DuplicatedLinesDensity, err = errors.Convert01(strconv.ParseFloat(v.Value, 32))
					if err != nil {
						return nil, err
					}
				case "code_smells":
					fileMetrics.CodeSmells, err = errors.Convert01(strconv.Atoi(v.Value))
					if err != nil {
						return nil, err
					}
				case "bugs":
					fileMetrics.Bugs, err = errors.Convert01(strconv.Atoi(v.Value))
					if err != nil {
						return nil, err
					}
				case "vulnerabilities":
					fileMetrics.Vulnerabilities, err = errors.Convert01(strconv.Atoi(v.Value))
					if err != nil {
						return nil, err
					}
				case "security_hotspots":
					fileMetrics.SecurityHotspots, err = errors.Convert01(strconv.Atoi(v.Value))
					if err != nil {
						return nil, err
					}
				case "security_hotspots_reviewed":
					fileMetrics.SecurityHotspotsReviewed, err = errors.Convert01(strconv.ParseFloat(v.Value, 32))
					if err != nil {
						return nil, err
					}
				case "coverage":
					fileMetrics.Coverage, err = errors.Convert01(strconv.ParseFloat(v.Value, 32))
					if err != nil {
						return nil, err
					}
				case "lines_to_cover":
					fileMetrics.LinesToCover, err = errors.Convert01(strconv.Atoi(v.Value))
					if err != nil {
						return nil, err
					}
				}
			}
			return []interface{}{fileMetrics}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractFilemetricsMeta = plugin.SubTaskMeta{
	Name:             "ExtractFilemetrics",
	EntryPoint:       ExtractFilemetrics,
	EnabledByDefault: true,
	Description:      "Extract raw data into tool layer table sonarqube_filemetrics",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

type fileMetricsResponse struct {
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	Qualifier string    `json:"qualifier"`
	Path      string    `json:"path"`
	Language  string    `json:"language"`
	Measures  []Measure `json:"measures"`
}
type Measure struct {
	Metric    string `json:"metric"`
	Value     string `json:"value"`
	BestValue bool   `json:"bestValue,omitempty"`
}
