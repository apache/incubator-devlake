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

package plugin

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/apache/incubator-devlake/core/dal"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
	"github.com/go-playground/validator/v10"
)

// ApiConnection represents a API Connection
type ApiConnection interface {
	GetEndpoint() string
	GetProxy() string
	GetRateLimitPerHour() int
}

type CacheableConnection interface {
	ApiConnection
	GetHash() string
}

// ApiAuthenticator is to be implemented by a Concreate Connection if Authorization is required
type ApiAuthenticator interface {
	// SetupAuthentication is a hook function for connection to set up authentication for the HTTP request
	// before sending it to the server
	SetupAuthentication(request *http.Request) errors.Error
}

// TODO: deprecated, remove
// ConnectionValidator represents the API Connection would validate its fields with customized logic
type ConnectionValidator interface {
	ValidateConnection(connection interface{}, valdator *validator.Validate) errors.Error
}

// PrepareApiClient is to be implemented by the concrete Connection which requires
// preparation for the ApiClient created by NewApiClientFromConnection, i.e. fetch token for future requests
type PrepareApiClient interface {
	PrepareApiClient(apiClient apihelperabstract.ApiClientAbstract) errors.Error
}

// MultiAuth
const (
	AUTH_METHOD_BASIC  = "BasicAuth"
	AUTH_METHOD_TOKEN  = "AccessToken"
	AUTH_METHOD_APPKEY = "AppKey"
)

var ALL_AUTH = map[string]bool{
	AUTH_METHOD_BASIC:  true,
	AUTH_METHOD_TOKEN:  true,
	AUTH_METHOD_APPKEY: true,
}

// MultiAuthenticator represents the API Connection supports multiple authorization methods
type MultiAuthenticator interface {
	GetAuthMethod() string
}

// BasicAuthenticator represents HTTP Basic Authentication
type BasicAuthenticator interface {
	GetBasicAuthenticator() ApiAuthenticator
}

// AccessTokenAuthenticator represents HTTP Bearer Authentication with Access Token
type AccessTokenAuthenticator interface {
	GetAccessTokenAuthenticator() ApiAuthenticator
}

// AppKeyAuthenticator represents the API Key and Secret authentication mechanism
type AppKeyAuthenticator interface {
	GetAppKeyAuthenticator() ApiAuthenticator
}

// Scope represents the top level entity for a data source, i.e. github repo,
// gitlab project, jira board. They turn into repo, board in Domain Layer. In
// Apache Devlake, a Project is essentially a set of these top level entities,
// for the framework to maintain these relationships dynamically and
// automatically, all Domain Layer Top Level Entities should implement this
// interface
type Scope interface {
	ScopeId() string
	ScopeName() string
	TableName() string
}

type ToolLayerConnection interface {
	dal.Tabler
	ConnectionId() uint64
}

type ToolLayerScope interface {
	dal.Tabler
	Scope
	ScopeFullName() string
	ScopeParams() interface{}
	ScopeConnectionId() uint64
	ScopeScopeConfigId() uint64
}

type ToolLayerScopeConfig interface {
	dal.Tabler
	ScopeConfigId() uint64
	ScopeConfigConnectionId() uint64
}

type ApiScope interface {
	ConvertApiScope() ToolLayerScope
}

type ApiGroup interface {
	GroupId() string
	GroupName() string
}

func MarshalScopeParams(params interface{}) string {
	bytes, err := json.Marshal(params)
	if err != nil {
		panic(fmt.Errorf("Failed to marshal %v, due to: %v", params, err))
	}
	return string(bytes)
}
