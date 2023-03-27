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

package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type RemoteScopesChild struct {
	Type     string      `json:"type"`
	ParentId *string     `json:"parentId"`
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	Data     interface{} `json:"data"`
}

type RemoteScopesOutput struct {
	Children      []RemoteScopesChild `json:"children"`
	NextPageToken string              `json:"nextPageToken"`
}

type SearchRemoteScopesOutput struct {
	Children []RemoteScopesChild `json:"children"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
}

type PageData struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

type TeamResponse struct {
	Offset int  `json:"offset"`
	Limit  int  `json:"limit"`
	More   bool `json:"more"`
	Total  int  `json:"total"`
	Teams  []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"teams"`
}

type WorkspaceItem struct {
	//Type string `json:"type"`
	//Uuid string `json:"uuid"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

type ReposResponse struct {
	Pagelen int              `json:"pagelen"`
	Page    int              `json:"page"`
	Size    int              `json:"size"`
	Values  []models.Service `json:"values"`
}

type ServiceResponse struct {
	Offset   int  `json:"offset"`
	Limit    int  `json:"limit"`
	More     bool `json:"more"`
	Total    int  `json:"total"`
	Services []struct {
		Id                      string    `json:"id"`
		Summary                 string    `json:"summary"`
		Type                    string    `json:"type"`
		Self                    string    `json:"self"`
		HtmlUrl                 string    `json:"html_url"`
		Name                    string    `json:"name"`
		AutoResolveTimeout      int       `json:"auto_resolve_timeout"`
		AcknowledgementTimeout  int       `json:"acknowledgement_timeout"`
		CreatedAt               time.Time `json:"created_at"`
		Status                  string    `json:"status"`
		AlertCreation           string    `json:"alert_creation"`
		AlertGroupingParameters struct {
			Type string `json:"type"`
		} `json:"alert_grouping_parameters"`
		Integrations []struct {
			Id      string `json:"id"`
			Type    string `json:"type"`
			Summary string `json:"summary"`
			Self    string `json:"self"`
			HtmlUrl string `json:"html_url"`
		} `json:"integrations"`
		EscalationPolicy struct {
			Id      string `json:"id"`
			Type    string `json:"type"`
			Summary string `json:"summary"`
			Self    string `json:"self"`
			HtmlUrl string `json:"html_url"`
		} `json:"escalation_policy"`
		Teams []struct {
			Id      string `json:"id"`
			Type    string `json:"type"`
			Summary string `json:"summary"`
			Self    string `json:"self"`
			HtmlUrl string `json:"html_url"`
		} `json:"teams"`
		IncidentUrgencyRule struct {
			Type               string `json:"type"`
			DuringSupportHours struct {
				Type    string `json:"type"`
				Urgency string `json:"urgency"`
			} `json:"during_support_hours"`
			OutsideSupportHours struct {
				Type    string `json:"type"`
				Urgency string `json:"urgency"`
			} `json:"outside_support_hours"`
		} `json:"incident_urgency_rule"`
		SupportHours struct {
			Type       string `json:"type"`
			TimeZone   string `json:"time_zone"`
			StartTime  string `json:"start_time"`
			EndTime    string `json:"end_time"`
			DaysOfWeek []int  `json:"days_of_week"`
		} `json:"support_hours"`
		ScheduledActions []struct {
			Type string `json:"type"`
			At   struct {
				Type string `json:"type"`
				Name string `json:"name"`
			} `json:"at"`
			ToUrgency string `json:"to_urgency"`
		} `json:"scheduled_actions"`
		AutoPauseNotificationsParameters struct {
			Enabled bool `json:"enabled"`
			Timeout int  `json:"timeout"`
		} `json:"auto_pause_notifications_parameters"`
	}
}

const RemoteScopesPerPage int = 100
const TypeScope string = "scope"
const TypeGroup string = "group"

//const TypeGroup string = "group"

// RemoteScopes list all available scopes (services) for this connection
// @Summary list all available scopes (services) for this connection
// @Description list all available scopes (services) for this connection
// @Tags plugins/pagerduty
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} RemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}

	connection := &models.PagerDutyConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}

	pageToken, ok := input.Query["pageToken"]
	if !ok || len(pageToken) == 0 {
		pageToken = []string{""}
	}

	pageData, err := DecodeFromPageToken(pageToken[0])
	if err != nil {
		return nil, errors.BadInput.New("failed to get page token")
	}

	// create api client
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, connection)
	if err != nil {
		return nil, err
	}

	query, err := GetQueryFromPageData(pageData)
	if err != nil {
		return nil, err
	}

	var res *http.Response
	outputBody := &RemoteScopesOutput{}
	res, err = apiClient.Get("/services", query, nil)
	if err != nil {
		return nil, err
	}
	response := &ServiceResponse{}
	err = api.UnmarshalResponse(res, response)
	if err != nil {
		return nil, err
	}
	// append service to output
	for _, service := range response.Services {
		child := RemoteScopesChild{
			Type: TypeScope,
			Id:   service.Id,
			Name: service.Name,
			Data: models.Service{
				Url:                  service.HtmlUrl,
				Id:                   service.Id,
				TransformationRuleId: 0, // this is not determined here
				Name:                 service.Name,
			},
		}
		outputBody.Children = append(outputBody.Children, child)
	}

	// check service count
	if !response.More {
		pageData = nil
	}

	// get the next page token
	outputBody.NextPageToken = ""
	if pageData != nil {
		pageData.Page += 1
		outputBody.NextPageToken, err = EncodeToPageToken(pageData)
		if err != nil {
			return nil, err
		}
	}

	return &plugin.ApiResourceOutput{Body: outputBody, Status: http.StatusOK}, nil
}

// SearchRemoteScopes use the Search API and only return project
// @Summary use the Search API and only return project
// @Description use the Search API and only return project
// @Tags plugins/pagerduty
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} SearchRemoteScopesOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// Not supported
	return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusMethodNotAllowed}, nil
}

func EncodeToPageToken(pageData *PageData) (string, errors.Error) {
	// Marshal json
	pageTokenDecode, err := json.Marshal(pageData)
	if err != nil {
		return "", errors.Default.Wrap(err, fmt.Sprintf("Marshal pageToken failed %+v", pageData))
	}
	// Encode pageToken Base64
	return base64.StdEncoding.EncodeToString(pageTokenDecode), nil
}

func DecodeFromPageToken(pageToken string) (*PageData, errors.Error) {
	if pageToken == "" {
		return &PageData{
			Page:    0,
			PerPage: RemoteScopesPerPage,
		}, nil
	}
	// Decode pageToken Base64
	pageTokenDecode, err := base64.StdEncoding.DecodeString(pageToken)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("decode pageToken failed %s", pageToken))
	}
	// Unmarshal json
	pt := &PageData{}
	err = json.Unmarshal(pageTokenDecode, pt)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("json Unmarshal pageTokenDecode failed %s", pageTokenDecode))
	}
	return pt, nil
}

func GetQueryFromPageData(pageData *PageData) (url.Values, errors.Error) {
	query := url.Values{}
	query.Set("offset", fmt.Sprintf("%v", pageData.Page))
	query.Set("limit", fmt.Sprintf("%v", pageData.PerPage))
	return query, nil
}

func extractParam(params map[string]string) (uint64, uint64) {
	connectionId, _ := strconv.ParseUint(params["connectionId"], 10, 64)
	serviceId, _ := strconv.ParseUint(params["serviceId"], 10, 64)
	return connectionId, serviceId
}
