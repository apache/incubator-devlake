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
	"net/http"
	"net/url"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_INCIDENTS_TABLE = "incidentio_incidents"

var _ plugin.SubTaskEntryPoint = CollectIncidents

type collectedIncidents struct {
	Incidents      []json.RawMessage        `json:"incidents"`
	PaginationMeta *collectedPaginationMeta `json:"pagination_meta"`
}

type collectedPaginationMeta struct {
	After    *string `json:"after"`
	PageSize int     `json:"page_size"`
}

var CollectIncidentsMeta = plugin.SubTaskMeta{
	Name:             "collectIncidents",
	EntryPoint:       CollectIncidents,
	EnabledByDefault: true,
	Description:      "Collect incident.io incidents",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
	ProductTables:    []string{RAW_INCIDENTS_TABLE},
}

func CollectIncidents(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*IncidentioTaskData)
	args := api.RawDataSubTaskArgs{
		Ctx:     taskCtx,
		Options: data.Options,
		Table:   RAW_INCIDENTS_TABLE,
	}
	// Pagination state captured during ResponseParser and consulted in
	// GetNextPageCustomData. Required because prevPageResponse.Body is
	// a single-read stream and is already drained by the time the
	// next-page hook fires.
	var lastAfter *string

	collector, err := api.NewStatefulApiCollectorForFinalizableEntity(api.FinalizableApiCollectorArgs{
		RawDataSubTaskArgs: args,
		ApiClient:          data.Client,
		CollectNewRecordsByList: api.FinalizableApiCollectorListArgs{
			PageSize: 250,
			GetNextPageCustomData: func(prevReqData *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
				// Safety cap against an upstream that returns full pages forever
				// while echoing a non-empty `after` cursor on every page.
				const maxPages = 10000
				if prevReqData.Pager.Page >= maxPages {
					return nil, api.ErrFinishCollect
				}
				if lastAfter == nil || *lastAfter == "" {
					return nil, api.ErrFinishCollect
				}
				return *lastAfter, nil
			},
			FinalizableApiCollectorCommonArgs: api.FinalizableApiCollectorCommonArgs{
				UrlTemplate: "v2/incidents",
				// incident.io does not support server-side filtering by incident
				// type or update time on this endpoint, so every scope collects
				// all incidents and the extractor filters; createdAfter is
				// intentionally ignored.
				Query: func(reqData *api.RequestData, createdAfter *time.Time) (url.Values, errors.Error) {
					after := ""
					if cursor, ok := reqData.CustomData.(string); ok {
						after = cursor
					}
					return buildIncidentsQuery(reqData.Pager.Size, after), nil
				},
				ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
					rawResult := collectedIncidents{}
					if err := api.UnmarshalResponse(res, &rawResult); err != nil {
						return nil, err
					}
					if rawResult.PaginationMeta != nil {
						lastAfter = rawResult.PaginationMeta.After
					} else {
						lastAfter = nil
					}
					return rawResult.Incidents, nil
				},
			},
		},
	})
	if err != nil {
		return err
	}
	return collector.Execute()
}

// buildIncidentsQuery is the pure-function core of the Query closure
// above. incident.io paginates with an opaque `after` cursor taken from
// the previous response's pagination_meta; the first page sends no cursor.
func buildIncidentsQuery(pageSize int, after string) url.Values {
	query := url.Values{}
	query.Set("page_size", fmt.Sprintf("%d", pageSize))
	if after != "" {
		query.Set("after", after)
	}
	return query
}
