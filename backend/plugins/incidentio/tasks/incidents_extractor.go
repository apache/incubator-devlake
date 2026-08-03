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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/incidentio/models"
	"github.com/apache/incubator-devlake/plugins/incidentio/models/raw"
)

var _ plugin.SubTaskEntryPoint = ExtractIncidents

var ExtractIncidentsMeta = plugin.SubTaskMeta{
	Name:             "extractIncidents",
	EntryPoint:       ExtractIncidents,
	EnabledByDefault: true,
	Description:      "Extract incident.io incidents",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	ProductTables:    []string{models.Incident{}.TableName()},
}

func ExtractIncidents(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*IncidentioTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_INCIDENTS_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			return extractIncidentioIncident(row.Data, data.Options)
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}

func extractIncidentioIncident(rawData []byte, op *IncidentioOptions) ([]interface{}, errors.Error) {
	rawIncident := &raw.Incident{}
	if err := errors.Convert(json.Unmarshal(rawData, rawIncident)); err != nil {
		return nil, err
	}

	// "test" and "tutorial" incidents are practice data; keep only
	// "standard" and "retrospective" (and future real modes).
	if rawIncident.Mode == "test" || rawIncident.Mode == "tutorial" {
		return nil, nil
	}

	// The collector fetches all incidents (no server-side filtering), so
	// the scope filter lives here. When IncidentTypeId is empty we are
	// collecting all incidents globally, so skip this check.
	if op.IncidentTypeId != "" {
		if rawIncident.IncidentType == nil || rawIncident.IncidentType.Id != op.IncidentTypeId {
			return nil, nil
		}
	}

	if rawIncident.CreatedAt.IsZero() {
		return nil, errors.Default.New("incident.io incident missing created_at")
	}

	declaredDate := rawIncident.CreatedAt
	if declared := timestampValueByName(rawIncident.IncidentTimestampValues, "Declared at"); declared != nil {
		declaredDate = *declared
	}
	resolvedDate := timestampValueByName(rawIncident.IncidentTimestampValues, "Resolved at")
	if resolvedDate == nil {
		resolvedDate = timestampValueByName(rawIncident.IncidentTimestampValues, "Closed at")
	}

	incident := &models.Incident{
		ConnectionId: op.ConnectionId,
		Id:           rawIncident.Id,
		Reference:    rawIncident.Reference,
		Name:         rawIncident.Name,
		Summary:      resolve(rawIncident.Summary),
		Url:          resolve(rawIncident.Permalink),
		Mode:         rawIncident.Mode,
		CreatedDate:  rawIncident.CreatedAt,
		UpdatedDate:  rawIncident.UpdatedAt,
		DeclaredDate: declaredDate,
		ResolvedDate: resolvedDate,
	}
	if rawIncident.IncidentStatus != nil {
		incident.StatusName = rawIncident.IncidentStatus.Name
		incident.StatusCategory = rawIncident.IncidentStatus.Category
	}
	if rawIncident.Severity != nil {
		incident.SeverityName = rawIncident.Severity.Name
		incident.SeverityRank = rawIncident.Severity.Rank
	}
	if rawIncident.IncidentType != nil {
		incident.IncidentTypeId = rawIncident.IncidentType.Id
	}

	return []interface{}{incident}, nil
}

func timestampValueByName(values []raw.IncidentTimestampValue, name string) *time.Time {
	for _, v := range values {
		if v.IncidentTimestamp.Name != name {
			continue
		}
		if v.Value == nil || v.Value.Value.IsZero() {
			return nil
		}
		value := v.Value.Value
		return &value
	}
	return nil
}

func resolve[T any](t *T) T {
	if t == nil {
		return *new(T)
	}
	return *t
}
