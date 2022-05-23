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

package apiv2models

import (
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

type Board struct {
	ID       uint64 `json:"id"`
	Self     string `json:"self"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Location *struct {
		ProjectId      uint   `json:"projectId"`
		DisplayName    string `json:"displayName"`
		ProjectName    string `json:"projectName"`
		ProjectKey     string `json:"projectKey"`
		ProjectTypeKey string `json:"projectTypeKey"`
		AvatarURI      string `json:"avatarURI"`
		Name           string `json:"name"`
	} `json:"location"`
}

func (b Board) ToToolLayer(connectionId uint64) *models.JiraBoard {
	result := &models.JiraBoard{
		ConnectionId: connectionId,
		BoardId:      b.ID,
		Name:         b.Name,
		Self:         b.Self,
		Type:         b.Type,
	}
	if b.Location != nil {
		result.ProjectId = b.Location.ProjectId
	}
	return result
}
