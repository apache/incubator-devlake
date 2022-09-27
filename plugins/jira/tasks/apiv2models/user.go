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

import "github.com/apache/incubator-devlake/plugins/jira/models"

type Account struct {
	Self         string `json:"self"`
	Key          string `json:"key"`
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	AccountId    string `json:"accountId"`
	AccountType  string `json:"accountType"`
	AvatarUrls   struct {
		Four8X48  string `json:"48x48"`
		Two4X24   string `json:"24x24"`
		One6X16   string `json:"16x16"`
		Three2X32 string `json:"32x32"`
	} `json:"avatarUrls"`
	DisplayName string `json:"displayName"`
	Active      bool   `json:"active"`
	Deleted     bool   `json:"deleted"`
	TimeZone    string `json:"timeZone"`
	Locale      string `json:"locale"`
}

func (u *Account) getAccountId() string {
	if u == nil {
		return ""
	}
	if u.AccountId != "" {
		return u.AccountId
	}
	if u.Key != "" {
		return u.Key
	}
	return u.EmailAddress
}

func (u *Account) ToToolLayer(connectionId uint64) *models.JiraAccount {
	return &models.JiraAccount{
		ConnectionId: connectionId,
		AccountId:    u.getAccountId(),
		AccountType:  u.AccountType,
		Name:         u.DisplayName,
		Email:        u.EmailAddress,
		Timezone:     u.TimeZone,
		AvatarUrl:    u.AvatarUrls.Four8X48,
	}
}
