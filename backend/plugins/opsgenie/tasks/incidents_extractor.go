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
	"github.com/apache/incubator-devlake/plugins/opsgenie/models"
	"github.com/apache/incubator-devlake/plugins/opsgenie/models/raw"
)

var _ plugin.SubTaskEntryPoint = ExtractIncidents

var ExtractIncidentsMeta = plugin.SubTaskMeta{
	Name:             "extractIncidents",
	EntryPoint:       ExtractIncidents,
	EnabledByDefault: true,
	Description:      "Extract Opsgenie incidents",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	DependencyTables: []string{RAW_INCIDENTS_TABLE},
}

func ExtractIncidents(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*OpsgenieTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_INCIDENTS_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			incidentRaw := &raw.Incident{}
			err := errors.Convert(json.Unmarshal(row.Data, incidentRaw))
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0, 1)
			incident := models.Incident{
				ConnectionId: data.Options.ConnectionId,
				Id:           *incidentRaw.Id,
				Url:          resolve(incidentRaw.Links.Web),
				Message:      *incidentRaw.Message,
				OwnerTeam:    resolve(incidentRaw.OwnerTeam),
				Description:  resolve(incidentRaw.Description),
				ServiceId:    data.Options.ServiceId,
				ServiceName:  data.Options.ServiceName,
				Status:       models.IncidentStatus(*incidentRaw.Status),
				Priority:     models.IncidentPriority(*incidentRaw.Priority),
				CreatedDate:  *incidentRaw.CreatedAt,
				UpdatedDate:  *incidentRaw.UpdatedAt,
			}
			results = append(results, &incident)
			for _, responderRaw := range *incidentRaw.Responders {
				if err != nil {
					return nil, err
				}
				results = append(results, &models.Assignment{
					ConnectionId: data.Options.ConnectionId,
					ResponderId:  *responderRaw.Id,
					IncidentId:   *incidentRaw.Id,
				})
				results = append(results, &models.Responder{
					ConnectionId: data.Options.ConnectionId,
					Id:           *responderRaw.Id,
					Type:         *responderRaw.Type,
				})
			}
			return results, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}
