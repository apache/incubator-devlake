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
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/codequality"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	sonarqubeModels "github.com/apache/incubator-devlake/plugins/sonarqube/models"
)

var ConvertProjectMetricsHistoryMeta = plugin.SubTaskMeta{
	Name:             "convertProjectMetricsHistory",
	EntryPoint:       ConvertProjectMetricsHistory,
	EnabledByDefault: true,
	Description:      "Convert tool layer project metrics history into domain layer table cq_project_metrics_history",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func ConvertProjectMetricsHistory(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	_, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_METRICS_HISTORY_TABLE)

	// Query all narrow metric rows ordered so we can group-pivot them
	cursor, err := db.Cursor(
		dal.From(sonarqubeModels.SonarqubeProjectMetricsHistory{}),
		dal.Where("connection_id = ? AND project_key = ?", data.Options.ConnectionId, data.Options.ProjectKey),
		dal.Orderby("analysis_date"),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	projectIdGen := didgen.NewDomainIdGenerator(&sonarqubeModels.SonarqubeProject{})
	domainProjectKey := projectIdGen.Generate(data.Options.ConnectionId, data.Options.ProjectKey)

	batchSave, err := helper.NewBatchSave(taskCtx, reflect.TypeOf(&codequality.CqProjectMetricsHistory{}), 200)
	if err != nil {
		return err
	}
	defer batchSave.Close()

	// Group narrow rows by analysis_date, pivot into wide domain rows
	var currentDate *time.Time
	var currentDomain *codequality.CqProjectMetricsHistory

	flushCurrent := func() errors.Error {
		if currentDomain != nil {
			return batchSave.Add(currentDomain)
		}
		return nil
	}

	for cursor.Next() {
		row := &sonarqubeModels.SonarqubeProjectMetricsHistory{}
		err = db.Fetch(cursor, row)
		if err != nil {
			return err
		}

		if currentDate == nil || !currentDate.Equal(row.AnalysisDate) {
			if flushErr := flushCurrent(); flushErr != nil {
				return flushErr
			}
			domainId := fmt.Sprintf("%s:%s",
				domainProjectKey,
				row.AnalysisDate.UTC().Format("2006-01-02T15:04:05Z"),
			)
			currentDomain = &codequality.CqProjectMetricsHistory{
				DomainEntity: domainlayer.DomainEntity{Id: domainId},
				ProjectKey:   domainProjectKey,
				AnalysisDate: row.AnalysisDate,
			}
			t := row.AnalysisDate
			currentDate = &t
		}

		applyMetricValue(currentDomain, row.MetricKey, row.MetricValue)
	}

	if flushErr := flushCurrent(); flushErr != nil {
		return flushErr
	}

	return batchSave.Close()
}

func applyMetricValue(d *codequality.CqProjectMetricsHistory, metricKey, value string) {
	switch metricKey {
	case "coverage":
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			d.Coverage = &v
		}
	case "ncloc":
		if v, err := strconv.Atoi(value); err == nil {
			d.Ncloc = &v
		}
	case "bugs":
		if v, err := strconv.Atoi(value); err == nil {
			d.Bugs = &v
		}
	case "reliability_rating":
		d.ReliabilityRating = alphabetMap[value]
	case "code_smells":
		if v, err := strconv.Atoi(value); err == nil {
			d.CodeSmells = &v
		}
	case "sqale_rating":
		d.SqaleRating = alphabetMap[value]
	case "complexity":
		if v, err := strconv.Atoi(value); err == nil {
			d.Complexity = &v
		}
	case "cognitive_complexity":
		if v, err := strconv.Atoi(value); err == nil {
			d.CognitiveComplexity = &v
		}
	case "vulnerabilities":
		if v, err := strconv.Atoi(value); err == nil {
			d.Vulnerabilities = &v
		}
	case "security_rating":
		d.SecurityRating = alphabetMap[value]
	case "security_hotspots":
		if v, err := strconv.Atoi(value); err == nil {
			d.SecurityHotspots = &v
		}
	case "duplicated_lines_density":
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			d.DuplicatedLinesDensity = &v
		}
	}
}
