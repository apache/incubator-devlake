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

	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/incidentio/models"
	"github.com/apache/devlake/plugins/incidentio/models/raw"
)

var _ plugin.SubTaskEntryPoint = ExtractIncidentTypes

var ExtractIncidentTypesMeta = plugin.SubTaskMeta{
	Name:             "extractIncidentTypes",
	EntryPoint:       ExtractIncidentTypes,
	EnabledByDefault: true,
	Description:      "Extract incident.io incident types",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	ProductTables:    []string{models.IncidentType{}.TableName()},
}

func ExtractIncidentTypes(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*IncidentioTaskData)
	db := taskCtx.GetDal()
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_INCIDENT_TYPES_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			rawIncidentType := &raw.IncidentType{}
			if err := errors.Convert(json.Unmarshal(row.Data, rawIncidentType)); err != nil {
				return nil, err
			}
			// The collector fetches the full incident type list; keep only
			// the type this scope is bound to.
			if data.Options.IncidentTypeId != "" && rawIncidentType.Id != data.Options.IncidentTypeId {
				return nil, nil
			}
			incidentType := &models.IncidentType{
				Id:   rawIncidentType.Id,
				Name: rawIncidentType.Name,
			}
			incidentType.ConnectionId = data.Options.ConnectionId
			// Preserve operator-set ScopeConfigId across re-collections.
			existing := &models.IncidentType{}
			if err := db.First(existing, dal.Where("connection_id = ? AND id = ?", data.Options.ConnectionId, rawIncidentType.Id)); err == nil {
				incidentType.ScopeConfigId = existing.ScopeConfigId
			}
			return []interface{}{incidentType}, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}
