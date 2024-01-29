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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
)

type ApiUserResponse struct {
	BitbucketId  int    `json:"id"`
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	Active       bool   `json:"active"`
	DisplayName  string `json:"displayName"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
	Links        struct {
		Self []struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"links"`
}

func convertUser(res *ApiUserResponse, connId uint64) (*models.BitbucketServerUser, errors.Error) {
	bitbucketUser := &models.BitbucketServerUser{
		ConnectionId: connId,
		BitbucketId:  res.BitbucketId,
		Name:         res.Name,
		EmailAddress: res.EmailAddress,
		Active:       res.Active,
		DisplayName:  res.DisplayName,
		Slug:         res.Slug,
		Type:         res.Type,
	}

	if len(res.Links.Self) > 0 {
		bitbucketUser.HtmlUrl = &res.Links.Self[0].Href
	}

	return bitbucketUser, nil
}
