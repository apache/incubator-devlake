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
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
)

type BambooConnection struct {
	api.BaseConnection `mapstructure:",squash"`
	BambooConn         `mapstructure:",squash"`
}

// TODO Please modify the following code to fit your needs
// This object conforms to what the frontend currently sends.
type BambooConn struct {
	api.RestConnection `mapstructure:",squash"`
	//TODO you may need to use helper.BasicAuth instead of helper.AccessToken
	api.AccessToken `mapstructure:",squash"`
}

// PrepareApiClient test api and set the IsPrivateToken,version,UserId and so on.
func (conn *BambooConn) PrepareApiClient(apiClient apihelperabstract.ApiClientAbstract) errors.Error {
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %v", conn.Token))

	res, err := apiClient.Get("", nil, header)
	if err != nil {
		return errors.HttpStatus(400).New(fmt.Sprintf("Get failed %s", err.Error()))
	}
	resources := &ApiXMLResourcesResponse{}

	if res.StatusCode != http.StatusOK {
		return errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code: %d", res.StatusCode))
	}

	err = api.UnmarshalResponseXML(res, resources)

	if err != nil {
		return errors.HttpStatus(500).New(fmt.Sprintf("UnmarshalResponse failed %s", err.Error()))
	}

	if resources.Link.Href != conn.Endpoint {
		return errors.HttpStatus(400).New(fmt.Sprintf("Response Data error for connection endpoint: %s, it should like: http://{domain}/rest/api/latest/", conn.Endpoint))
	}

	res, err = apiClient.Get("repository", nil, header)
	if err != nil {
		return errors.HttpStatus(400).New(fmt.Sprintf("Get failed %s", err.Error()))
	}
	repo := &ApiRepository{}

	if res.StatusCode != http.StatusOK {
		return errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code: %d", res.StatusCode))
	}

	err = api.UnmarshalResponse(res, repo)

	if err != nil {
		return errors.HttpStatus(400).New(fmt.Sprintf("UnmarshalResponse repository failed %s", err.Error()))
	}

	return nil
}

// This object conforms to what the frontend currently expects.
type BambooResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	BambooConnection
}

type ApiRepository struct {
	Size          int         `json:"size"`
	SearchResults interface{} `json:"searchResults"`
	StartIndex    int         `json:"start-index"`
	MaxResult     int         `json:"max-result"`
}

type ApiXMLLink struct {
	XMLName xml.Name `json:"xml_name" xml:"link"`
	Href    string   `json:"href" xml:"href,attr"`
}

type ApiXMLResource struct {
	XMLName xml.Name `json:"xml_name" xml:"resource"`
	Name    string   `json:"name" xml:"name,attr"`
}

type ApiXMLResources struct {
	XMLName   xml.Name         `json:"xml_name" xml:"resources"`
	Size      string           `json:"size" xml:"size,attr"`
	Resources []ApiXMLResource `json:"resource" xml:"resource"`
}

// Using User because it requires authentication.
type ApiXMLResourcesResponse struct {
	XMLName xml.Name `json:"xml_name" xml:"resources"`
	Expand  string   `json:"expand" xml:"expand,attr"`

	Link      ApiXMLLink      `json:"link" xml:"link"`
	Resources ApiXMLResources `json:"resources" xml:"resources"`
}

func (BambooConnection) TableName() string {
	return "_tool_bamboo_connections"
}
