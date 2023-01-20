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
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

const RAW_INCIDENTS_TABLE = "pagerduty_incidents"

var _ plugin.SubTaskEntryPoint = CollectIncidents

type (
	pagingInfo struct {
		Limit  *int  `json:"limit"`
		Offset *int  `json:"offset"`
		Total  *int  `json:"total"`
		More   *bool `json:"more"`
	}
	collectedIncidents struct {
		pagingInfo
		Incidents []json.RawMessage `json:"incidents"`
	}

	collectedIncident struct {
		pagingInfo
		Incident json.RawMessage `json:"incident"`
	}
	simplifiedRawIncident struct {
		IncidentNumber int       `json:"incident_number"`
		CreatedAt      time.Time `json:"created_at"`
	}
)

func CollectIncidents(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*PagerDutyTaskData)
	db := taskCtx.GetDal()
	args := api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: PagerDutyParams{
			ConnectionId: data.Options.ConnectionId,
		},
		Table: RAW_INCIDENTS_TABLE,
	}
	collector, err := api.NewStatefulApiCollectorForFinalizableEntity(api.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: args,
		ApiClient:          data.Client,
		TimeAfter:          data.TimeAfter,
		CollectNewRecordsByList: api.FinalizableApiCollectorListArgs{
			PageSize: 1,
			GetCreated: func(item json.RawMessage) (time.Time, errors.Error) {
				incident := &simplifiedRawIncident{}
				err := json.Unmarshal(item, incident)
				if err != nil {
					return time.Time{}, errors.BadInput.Wrap(err, "failed to unmarshal incident")
				}
				return incident.CreatedAt, nil
			},
			FinalizableApiCollectorCommonArgs: api.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "incidents",
				Query: func(reqData *api.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					query := url.Values{}
					if createdAfter != nil {
						now := time.Now()
						if now.Sub(*createdAfter).Seconds() > 180*24*time.Hour.Seconds() {
							// beyond 6 months Pagerduty API will just return nothing, so need to query for 'all' instead
							query.Set("date_range", "all")
						} else {
							query.Set("since", data.TimeAfter.String())
						}
					} else {
						query.Set("date_range", "all")
					}
					query.Set("sort_by", "created_at:asc")
					query.Set("limit", fmt.Sprintf("%d", reqData.Pager.Size))
					query.Set("offset", fmt.Sprintf("%d", reqData.Pager.Page))
					query.Set("total", "true")
					return query, nil
				},
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					rawResult := collectedIncidents{}
					err := api.UnmarshalResponse(res, &rawResult)
					return rawResult.Incidents, err
				},
			},
			GetNextPageCustomData: func(prevReqData *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
				// not sure this is even necessary because the framework seems to auto-detect when to stop querying for the next page
				pageInfo := pagingInfo{}
				err := api.UnmarshalResponse(prevPageResponse, &pageInfo)
				if err != nil {
					return nil, err
				}
				if *pageInfo.Offset > *pageInfo.Total {
					return nil, api.ErrFinishCollect
				}
				return pageInfo, nil
			},
		},
		CollectUnfinishedDetails: api.FinalizableApiCollectorDetailArgs{
			FinalizableApiCollectorCommonArgs: api.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "incidents/{{ .Input.Number }}",
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					rawResult := collectedIncident{}
					err := api.UnmarshalResponse(res, &rawResult)
					return []json.RawMessage{rawResult.Incident}, err
				},
			},
			BuildInputIterator: func() (api.Iterator, errors.Error) {
				// select pull id from database
				cursor, err := db.Cursor(
					dal.Select("number, created_date"),
					dal.From(&models.Incident{}),
					dal.Where(
						"service_id = ? AND connection_id = ?",
						data.Options.ServiceId, data.Options.ConnectionId,
					),
				)
				if err != nil {
					return nil, err
				}
				return api.NewDalCursorIterator(db, cursor, reflect.TypeOf(simplifiedRawIncident{}))
			},
		},
	})
	if err != nil {
		return nil
	}
	return collector.Execute()
}

var CollectIncidentsMeta = plugin.SubTaskMeta{
	Name:             "collectIncidents",
	EntryPoint:       CollectIncidents,
	EnabledByDefault: true,
	Description:      "Collect PagerDuty incidents",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}
