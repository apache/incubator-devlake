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
	"github.com/apache/incubator-devlake/plugins/github/models"
)

type GithubAccountEdge struct {
	Login     string
	Id        int `graphql:"databaseId"`
	Name      string
	Company   string
	Email     string
	AvatarUrl string
	HtmlUrl   string `graphql:"url"`
	//Type      string
	Organizations struct {
		Nodes []struct {
			Email      string
			Name       string
			DatabaseId int
			Login      string
		}
	} `graphql:"organizations(first: 10)"`
}
type GithubAccountResponse struct {
	GithubAccountEdge `graphql:"... on User"`
}

func convertAccount(res GithubAccountResponse, connId uint64) ([]interface{}, error) {
	results := make([]interface{}, 0, len(res.Organizations.Nodes)+1)
	githubAccount := &models.GithubAccount{
		ConnectionId: connId,
		Id:           res.Id,
		Login:        res.Login,
		Name:         res.Name,
		Company:      res.Company,
		Email:        res.Email,
		AvatarUrl:    res.AvatarUrl,
		//Url:          res.Url,
		HtmlUrl: res.HtmlUrl,
		//Type:         res.Type,
	}
	results = append(results, githubAccount)
	for _, apiAccountOrg := range res.Organizations.Nodes {
		githubAccountOrg := &models.GithubAccountOrg{
			ConnectionId: connId,
			AccountId:    res.Id,
			OrgId:        apiAccountOrg.DatabaseId,
			OrgLogin:     apiAccountOrg.Login,
		}
		results = append(results, githubAccountOrg)
	}

	return results, nil
}
