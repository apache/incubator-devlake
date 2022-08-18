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
	"time"

	"github.com/apache/incubator-devlake/plugins/azure/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

type AzureApiRepo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	Project struct {
		ID             string    `json:"id"`
		Name           string    `json:"name"`
		URL            string    `json:"url"`
		State          string    `json:"state"`
		Revision       int       `json:"revision"`
		Visibility     string    `json:"visibility"`
		LastUpdateTime time.Time `json:"lastUpdateTime"`
	} `json:"project"`
	DefaultBranch string `json:"defaultBranch"`
	Size          int    `json:"size"`
	RemoteURL     string `json:"remoteUrl"`
	SSHURL        string `json:"sshUrl"`
	WebURL        string `json:"webUrl"`
	IsDisabled    bool   `json:"isDisabled"`
}

var ExtractApiRepoMeta = core.SubTaskMeta{
	Name:        "extractApiRepo",
	EntryPoint:  ExtractApiRepositories,
	Required:    true,
	Description: "Extract raw Repositories data into tool layer table azure_repos",
	DomainTypes: []string{core.DOMAIN_TYPE_CODE},
}

type ApiRepoResponse AzureApiRepo

func ExtractApiRepositories(taskCtx core.SubTaskContext) error {
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
			Table: RAW_REPOSITORIES_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			body := &ApiRepoResponse{}
			err := json.Unmarshal(row.Data, body)
			if err != nil {
				return nil, err
			}
			if body.ID == "" {
				return nil, fmt.Errorf("repo %s not found", data.Options.Project)
			}
			results := make([]interface{}, 0, 1)
			azureRepository := &models.AzureRepo{
				ConnectionId:  data.Options.ConnectionId,
				AzureId:       body.ID,
				Name:          body.Name,
				Url:           body.URL,
				ProjectId:     body.Project.ID,
				DefaultBranch: body.DefaultBranch,
				Size:          body.Size,
				RemoteURL:     body.RemoteURL,
				SshUrl:        body.SSHURL,
				WebUrl:        body.WebURL,
				IsDisabled:    body.IsDisabled,
			}
			data.Repo = azureRepository

			results = append(results, azureRepository)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}
