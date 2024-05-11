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

package azuredevops

import (
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const apiVersion = "7.1"
const maxPageSize = 100

type Client struct {
	c http.Client

	apiClient  plugin.ApiClient
	connection *models.AzuredevopsConnection
	url        string
}

func NewClient(con *models.AzuredevopsConnection, apiClient plugin.ApiClient, url string) Client {
	return Client{
		c: http.Client{
			Timeout: 2 * time.Second,
		},
		connection: con,
		url:        url,
		apiClient:  apiClient,
	}
}

func (c *Client) GetUserProfile() (Profile, errors.Error) {
	var p Profile
	endpoint, err := url.JoinPath(c.url, "/_apis/profile/profiles/me")
	if err != nil {
		return Profile{}, errors.Internal.Wrap(err, "failed to join user profile path")
	}

	res, err := c.doGet(endpoint)
	if err != nil {
		return Profile{}, errors.Internal.Wrap(err, "failed to read user accounts")
	}

	if res.StatusCode == 203 || res.StatusCode == 401 {
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

func (c *Client) GetUserAccounts(memberId string) (AccountResponse, errors.Error) {
	var a AccountResponse
	endpoint := fmt.Sprintf(c.url+"/_apis/accounts?memberId=%s", memberId)
	res, err := c.doGet(endpoint)
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

func (c *Client) doGet(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if err = c.connection.GetAccessTokenAuthenticator().SetupAuthentication(req); err != nil {
		return nil, errors.Internal.Wrap(err, "failed to authorize the request using the plugin connection")
	}
	return http.DefaultClient.Do(req)
}

type GetProjectsArgs struct {
	// (optional) Pagination
	*OffsetPagination

	OrgId string
}

func (c *Client) GetProjects(args GetProjectsArgs) ([]Project, errors.Error) {
	query := url.Values{}
	query.Set("api-version", apiVersion)

	var top, skip int
	top = maxPageSize
	skip = 0
	if args.OffsetPagination != nil {
		top = args.Top
		skip = args.Skip
	}

	var data struct {
		Count    int       `json:"count"`
		Projects []Project `json:"value"`
	}

	var projects []Project

	for {
		query.Set("$top", strconv.Itoa(top))
		query.Set("$skip", strconv.Itoa(skip))

		path := fmt.Sprintf("%s/_apis/projects", args.OrgId)
		res, err := c.apiClient.Get(path, query, nil)
		if err != nil {
			return nil, err
		}

		if res.StatusCode == 203 || res.StatusCode == 401 {
			return nil, errors.Unauthorized.New("failed to read projects")
		}

		if res.StatusCode != 200 {
			return nil, errors.Internal.New(fmt.Sprintf("failed to read projects, upstream api call failed with (%v)", res.StatusCode))
		}

		err = api.UnmarshalResponse(res, &data)
		if err != nil {
			return nil, err
		}

		projects = append(projects, data.Projects...)

		if data.Count < top {
			return projects, nil
		}

		skip += top
	}
}

type GetRepositoriesArgs struct {
	OrgId     string
	ProjectId string
}

func (c *Client) GetRepositories(args GetRepositoriesArgs) ([]Repository, errors.Error) {
	query := url.Values{}
	query.Set("api-version", apiVersion)

	var data struct {
		Repos []Repository `json:"value"`
	}

	path := fmt.Sprintf("%s/%s/_apis/git/repositories", args.OrgId, args.ProjectId)
	res, err := c.apiClient.Get(path, query, nil)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 401:
		fallthrough
	case 403:
		return nil, errors.Unauthorized.New("failed to authorize the '.../_apis/git/repositories' request using the plugin connection")
	case 404:
		return nil, errors.NotFound.New("failed to find requested resource on '.../_apis/git/repositories'")
	default:
	}

	err = api.UnmarshalResponse(res, &data)
	if err != nil {
		return nil, err
	}

	return data.Repos, nil
}

type GetServiceEndpointsArgs struct {
	ProjectId string
	OrgId     string
}

func (c *Client) GetServiceEndpoints(args GetServiceEndpointsArgs) ([]ServiceEndpoint, errors.Error) {
	query := url.Values{}
	query.Set("api-version", apiVersion)

	path := fmt.Sprintf("%s/%s/_apis/serviceendpoint/endpoints/", args.OrgId, args.ProjectId)
	res, err := c.apiClient.Get(path, query, nil)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 401:
		fallthrough
	case 403:
		return nil, errors.Unauthorized.New("failed to authorize the '.../serviceendpoint/endpoints' request using the plugin connection")
	case 404:
		return nil, errors.NotFound.New("failed to find requested resource on '.../serviceendpoint/endpoints'")
	default:
	}

	var data struct {
		ServiceEndpoints []ServiceEndpoint `json:"value"`
	}

	err = api.UnmarshalResponse(res, &data)
	if err != nil {
		return nil, err
	}
	return data.ServiceEndpoints, nil
}

type GetRemoteRepositoriesArgs struct {
	ProjectId string
	OrgId     string
	Provider  string
	// (optional) Service Endpoint to filter for
	ServiceEndpoint string
}

func (c *Client) GetRemoteRepositories(args GetRemoteRepositoriesArgs) ([]RemoteRepository, error) {
	query := url.Values{}
	query.Set("api-version", apiVersion)
	if args.ServiceEndpoint != "" {
		query.Set("serviceEndpointId", args.ServiceEndpoint)
	}

	var repos []RemoteRepository
	var response struct {
		Repository []RemoteRepository `json:"repositories"`
	}

	for {
		path := fmt.Sprintf("%s/%s/_apis/sourceProviders/%s/repositories/", args.OrgId, args.ProjectId, args.Provider)
		res, err := c.apiClient.Get(path, query, nil)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "failed to read remote repositories")
		}
		err = api.UnmarshalResponse(res, &response)
		if err != nil {
			return nil, errors.Internal.Wrap(err, "failed to unmarshal remote repositories response")
		}

		repos = append(repos, response.Repository...)
		contToken := res.Header.Get("X-Ms-Continuationtoken")
		if contToken == "" {
			return repos, nil
		}

		query.Set("continuationToken", contToken)
	}
}
