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

package models

import (
	"github.com/apache/incubator-devlake/core/models/common"
)

type GitlabAccount struct {
	ConnectionId    uint64 `gorm:"primaryKey"`
	GitlabId        int    `gorm:"primaryKey" json:"id"`
	Username        string `gorm:"type:varchar(255)"`
	Email           string `gorm:"type:varchar(255)"`
	Name            string `gorm:"type:varchar(255)"`
	State           string `gorm:"type:varchar(255)"`
	MembershipState string `json:"membership_state" gorm:"type:varchar(255)"`
	AvatarUrl       string `json:"avatar_url" gorm:"type:varchar(255)"`
	WebUrl          string `json:"web_url" gorm:"type:varchar(255)"`

	common.NoPKModel
}

func (GitlabAccount) TableName() string {
	return "_tool_gitlab_accounts"
}

type GitlabUser struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

type GitlabMember struct {
	AccessLevel     int        `json:"access_level"`
	CreatedAt       string     `json:"created_at"`
	CreatedBy       GitlabUser `json:"created_by"`
	ExpiresAt       string     `json:"expires_at"`
	ID              int        `json:"id"`
	Username        string     `json:"username"`
	Name            string     `json:"name"`
	State           string     `json:"state"`
	AvatarURL       string     `json:"avatar_url"`
	WebURL          string     `json:"web_url"`
	MembershipState string     `json:"membership_state"`
}
