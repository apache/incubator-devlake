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
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Profile struct {
	DisplayName  string    `json:"displayName"`
	PublicAlias  string    `json:"publicAlias"`
	EmailAddress string    `json:"emailAddress"`
	CoreRevision int       `json:"coreRevision"`
	TimeStamp    time.Time `json:"timeStamp"`
	Id           string    `json:"id"`
	Revision     int       `json:"revision"`
}

type Account struct {
	AccountId        string      `json:"AccountId"`
	NamespaceId      string      `json:"NamespaceId"`
	AccountName      string      `json:"AccountName"`
	OrganizationName interface{} `json:"OrganizationName"`
	AccountType      int         `json:"AccountType"`
	AccountOwner     string      `json:"AccountOwner"`
	CreatedBy        string      `json:"CreatedBy"`
	CreatedDate      string      `json:"CreatedDate"`
	AccountStatus    int         `json:"AccountStatus"`
	StatusReason     interface{} `json:"StatusReason"`
	LastUpdatedBy    string      `json:"LastUpdatedBy"`
	Properties       struct {
	} `json:"Properties"`
}

type AccountResponse []Account

type vsClient struct {
	c          http.Client
	connection *models.AzuredevopsConnection
	url        string
}

func newVsClient(con *models.AzuredevopsConnection, url string) vsClient {
	return vsClient{
		c: http.Client{
			Timeout: 2 * time.Second,
		},
		connection: con,
		url:        url,
	}
}

func (vsc *vsClient) UserProfile() (Profile, errors.Error) {
	var p Profile
	endpoint, err := url.JoinPath(vsc.url, "/_apis/profile/profiles/me")
	if err != nil {
		return Profile{}, errors.Internal.Wrap(err, "failed to join user profile path")
	}

	res, err := vsc.doGet(endpoint)
	if err != nil {
		return Profile{}, errors.Internal.Wrap(err, "failed to read user accounts")
	}

	if res.StatusCode == 302 || res.StatusCode == 401 {
		return Profile{}, errors.Unauthorized.New("failed to read user profile")
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return Profile{}, errors.Internal.Wrap(err, "failed to read response body")
	}

	if err := json.Unmarshal(resBody, &p); err != nil {
		panic(err)
	}
	return p, nil
}

func (vsc *vsClient) UserAccounts(memberId string) (AccountResponse, errors.Error) {
	var a AccountResponse
	endpoint := fmt.Sprintf(vsc.url+"/_apis/accounts?memberId=%s", memberId)
	res, err := vsc.doGet(endpoint)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to read user accounts")
	}

	if res.StatusCode == 302 || res.StatusCode == 401 {
		return nil, errors.Unauthorized.New("failed to read user accounts")
	}

	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Internal.Wrap(err, "failed to read response body")
	}

	if err := json.Unmarshal(resBody, &a); err != nil {
		return nil, errors.Internal.Wrap(err, "failed to read unmarshal response body")
	}
	return a, nil
}

func (vsc *vsClient) doGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if err = vsc.connection.GetAccessTokenAuthenticator().SetupAuthentication(req); err != nil {
		return nil, errors.Internal.Wrap(err, "failed to authorize the request using the plugin connection")
	}
	return http.DefaultClient.Do(req)
}
