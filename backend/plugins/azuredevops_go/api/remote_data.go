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
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"time"
)

type AzuredevopsRemotePagination struct {
	Skip int
	Top  int
}

type AzuredevopsApiRepo struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Url     string `json:"url"`
	Project struct {
		Id             string    `json:"id"`
		Name           string    `json:"name"`
		Description    string    `json:"description"`
		Url            string    `json:"url"`
		State          string    `json:"state"`
		Revision       int       `json:"revision"`
		Visibility     string    `json:"visibility"`
		LastUpdateTime time.Time `json:"lastUpdateTime"`
	} `json:"project"`
	DefaultBranch   string `json:"defaultBranch"`
	Size            int    `json:"size"`
	RemoteUrl       string `json:"remoteUrl"`
	SshUrl          string `json:"sshUrl"`
	WebUrl          string `json:"webUrl"`
	IsDisabled      bool   `json:"isDisabled"`
	IsInMaintenance bool   `json:"isInMaintenance"`
	IsFork          bool   `json:"isFork"`
}

func (r AzuredevopsApiRepo) toRepoModel() models.AzuredevopsRepo {
	return models.AzuredevopsRepo{
		Id:   r.Id,
		Name: r.Name,
		Url:  r.WebUrl,
		AzureDevOpsPK: models.AzureDevOpsPK{
			ProjectId: r.Project.Id,
		},
		RemoteUrl: r.RemoteUrl,
		IsFork:    r.IsFork,
	}
}
