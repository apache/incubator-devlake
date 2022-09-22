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
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

type BitbucketAccountResponse struct {
	UserName      string    `json:"username"`
	DisplayName   string    `json:"display_name"`
	AccountId     string    `json:"account_id"`
	AccountStatus string    `json:"account_status"`
	CreateOn      time.Time `json:"create_on"`
	Links         struct {
		Self       struct{ Href string } `json:"self"`
		Html       struct{ Href string } `json:"html"`
		Avatar     struct{ Href string } `json:"avatar"`
		Followers  struct{ Href string } `json:"followers"`
		Following  struct{ Href string } `json:"following"`
		Repository struct{ Href string } `json:"repository"`
	}
}

func convertAccount(res *BitbucketAccountResponse, connId uint64) (*models.BitbucketAccount, errors.Error) {
	bitbucketAccount := &models.BitbucketAccount{
		ConnectionId:  connId,
		UserName:      res.UserName,
		DisplayName:   res.DisplayName,
		AccountId:     res.AccountId,
		AccountStatus: res.AccountStatus,
		AvatarUrl:     res.Links.Avatar.Href,
		HtmlUrl:       res.Links.Html.Href,
	}
	return bitbucketAccount, nil
}
