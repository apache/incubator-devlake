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

	"github.com/apache/incubator-devlake/plugins/azure/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type AzureApiBuildDefinition struct {
	Quality    string `json:"quality"`
	AuthoredBy struct {
		DisplayName string `json:"displayName"`
		URL         string `json:"url"`
		Links       struct {
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"_links"`
		ID         string `json:"id"`
		UniqueName string `json:"uniqueName"`
		ImageURL   string `json:"imageUrl"`
		Descriptor string `json:"descriptor"`
	} `json:"authoredBy"`
	Queue struct {
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
		ID   int    `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
		Pool struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			IsHosted bool   `json:"isHosted"`
		} `json:"pool"`
	} `json:"queue"`
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	URI         string    `json:"uri"`
	Path        string    `json:"path"`
	Type        string    `json:"type"`
	QueueStatus string    `json:"queueStatus"`
	Revision    int       `json:"revision"`
	CreatedDate time.Time `json:"createdDate"`
	Project     struct {
		ID             string    `json:"id"`
		Name           string    `json:"name"`
		URL            string    `json:"url"`
		State          string    `json:"state"`
		Revision       int       `json:"revision"`
		Visibility     string    `json:"visibility"`
		LastUpdateTime time.Time `json:"lastUpdateTime"`
	} `json:"project"`
}

var ExtractApiBuildDefinitionMeta = core.SubTaskMeta{
	Name:        "extractApiBuild",
	EntryPoint:  ExtractApiBuildDefinition,
	Required:    true,
	Description: "Extract raw BuildDefinition data into tool layer table azure_repos",
	DomainTypes: []string{core.DOMAIN_TYPE_CICD},
}

func ExtractApiBuildDefinition(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*AzureTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: AzureApiParams{
				ConnectionId: data.Options.ConnectionId,
				Project:      data.Options.Project,
			},
			/*
				Table store raw data
			*/
			Table: RAW_BUILD_DEFINITION_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &AzureApiBuildDefinition{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}

			results := make([]interface{}, 0, 1)
			azureBuildDefinition := &models.AzureBuildDefinition{
				ConnectionId:     data.Options.ConnectionId,
				ProjectId:        body.Project.ID,
				AzureId:          body.ID,
				AuthorId:         body.AuthoredBy.ID,
				QueueId:          body.Queue.ID,
				Url:              body.URL,
				Name:             body.Name,
				Path:             body.Path,
				Type:             body.Type,
				QueueStatus:      body.QueueStatus,
				Revision:         body.Revision,
				AzureCreatedDate: body.CreatedDate,
			}
			results = append(results, azureBuildDefinition)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
