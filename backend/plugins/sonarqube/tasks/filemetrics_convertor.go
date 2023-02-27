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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/codequality"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	sonarqubeModels "github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"reflect"
)

var ConvertFileMetricsMeta = plugin.SubTaskMeta{
	Name:             "convertFileMetrics",
	EntryPoint:       ConvertFileMetrics,
	EnabledByDefault: true,
	Description:      "Convert tool layer table sonarqube_file_metrics into  domain layer table cq_file_metrics",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func ConvertFileMetrics(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_FILEMETRICS_TABLE)
	cursor, err := db.Cursor(dal.From(sonarqubeModels.SonarqubeWholeFileMetrics{}),
		dal.Where("connection_id = ? and project_key = ?", data.Options.ConnectionId, data.Options.ProjectKey))
	if err != nil {
		return err
	}
	defer cursor.Close()

	issueIdGen := didgen.NewDomainIdGenerator(&sonarqubeModels.SonarqubeFileMetrics{})
	projectIdGen := didgen.NewDomainIdGenerator(&sonarqubeModels.SonarqubeProject{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(sonarqubeModels.SonarqubeWholeFileMetrics{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			sonarqubeFileMetric := inputRow.(*sonarqubeModels.SonarqubeWholeFileMetrics)
			domainFileMetric := &codequality.CqFileMetrics{
				DomainEntity:                        domainlayer.DomainEntity{Id: issueIdGen.Generate(data.Options.ConnectionId, sonarqubeFileMetric.FileMetricsKey)},
				ProjectKey:                          projectIdGen.Generate(data.Options.ConnectionId, sonarqubeFileMetric.ProjectKey),
				FileName:                            sonarqubeFileMetric.FileName,
				FilePath:                            sonarqubeFileMetric.FilePath,
				FileLanguage:                        sonarqubeFileMetric.FileLanguage,
				CodeSmells:                          sonarqubeFileMetric.CodeSmells,
				SqaleIndex:                          sonarqubeFileMetric.SqaleIndex,
				SqaleRating:                         sonarqubeFileMetric.SqaleRating,
				Bugs:                                sonarqubeFileMetric.Bugs,
				ReliabilityRating:                   sonarqubeFileMetric.ReliabilityRating,
				Vulnerabilities:                     sonarqubeFileMetric.Vulnerabilities,
				SecurityRating:                      sonarqubeFileMetric.SecurityRating,
				SecurityHotspots:                    sonarqubeFileMetric.SecurityHotspots,
				SecurityHotspotsReviewed:            sonarqubeFileMetric.SecurityHotspotsReviewed,
				SecurityReviewRating:                sonarqubeFileMetric.SecurityReviewRating,
				Ncloc:                               sonarqubeFileMetric.Ncloc,
				UncoveredLines:                      sonarqubeFileMetric.UncoveredLines,
				LinesToCover:                        sonarqubeFileMetric.LinesToCover,
				DuplicatedLinesDensity:              sonarqubeFileMetric.DuplicatedLinesDensity,
				DuplicatedBlocks:                    sonarqubeFileMetric.DuplicatedBlocks,
				DuplicatedFiles:                     sonarqubeFileMetric.DuplicatedFiles,
				DuplicatedLines:                     sonarqubeFileMetric.DuplicatedLines,
				EffortToReachMaintainabilityRatingA: sonarqubeFileMetric.EffortToReachMaintainabilityRatingA,
				Complexity:                          sonarqubeFileMetric.Complexity,
				CognitiveComplexity:                 sonarqubeFileMetric.CognitiveComplexity,
				NumOfLines:                          sonarqubeFileMetric.NumOfLines,
				Coverage:                            sonarqubeFileMetric.Coverage,
			}
			return []interface{}{
				domainFileMetric,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}
