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
	"encoding/base64"
	"fmt"
	"github.com/apache/incubator-devlake/core/plugin"
	"net/http"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/go-playground/validator/v10"
)

// BasicAuth implements HTTP Basic Authentication
type BasicAuth struct {
	Username string `mapstructure:"username" validate:"required" json:"username"`
	Password string `mapstructure:"password" validate:"required" json:"password" gorm:"serializer:encdec"`
}

// GetEncodedToken returns encoded bearer token for HTTP Basic Authentication
func (ba *BasicAuth) GetEncodedToken() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", ba.Username, ba.Password)))
}

// SetupAuthentication sets up the request headers for authentication
func (ba *BasicAuth) SetupAuthentication(request *http.Request) errors.Error {
	request.Header.Set("Authorization", fmt.Sprintf("Basic %v", ba.GetEncodedToken()))
	return nil
}

// GetBasicAuthenticator returns the ApiAuthenticator for setting up the HTTP request
// it looks odd to return itself with a different type, this is necessary because Callers
// might call the method from the Outer-Struct(`connection.SetupAuthentication(...)`)
// which would lead to a Stack Overflow  error
func (ba *BasicAuth) GetBasicAuthenticator() plugin.ApiAuthenticator {
	return ba
}

// AccessToken implements HTTP Bearer Authentication with Access Token
type AccessToken struct {
	Token string `mapstructure:"token" validate:"required" json:"token" gorm:"serializer:encdec"`
}

// SetupAuthentication sets up the request headers for authentication
func (at *AccessToken) SetupAuthentication(request *http.Request) errors.Error {
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", at.Token))
	return nil
}

// GetAccessTokenAuthenticator returns SetupAuthentication
func (at *AccessToken) GetAccessTokenAuthenticator() plugin.ApiAuthenticator {
	return at
}

// AppKey implements the API Key and Secret authentication mechanism
type AppKey struct {
	AppId     string `mapstructure:"appId" validate:"required" json:"appId"`
	SecretKey string `mapstructure:"secretKey" validate:"required" json:"secretKey" gorm:"serializer:encdec"`
}

// SetupAuthentication sets up the request headers for authentication
func (ak *AppKey) SetupAuthentication(request *http.Request) errors.Error {
	// no universal way to implement AppKey authentication, plugin should alias AppKey and
	// define its own implementation if API requires signature for each request,
	// or you should implement PrepareApiClient if API requires a Token for all requests
	return nil
}

// GetAppKeyAuthenticator returns SetupAuthentication
func (ak *AppKey) GetAppKeyAuthenticator() plugin.ApiAuthenticator {
	// no universal way to implement AppKey authentication, plugin should alias AppKey and
	// define its own implementation
	return ak
}

// MultiAuth implements the MultiAuthenticator interface
type MultiAuth struct {
	AuthMethod       string `mapstructure:"authMethod" json:"authMethod" validate:"required,oneof=BasicAuth AccessToken AppKey"`
	apiAuthenticator plugin.ApiAuthenticator
}

func (ma *MultiAuth) GetApiAuthenticator(connection plugin.ApiConnection) (plugin.ApiAuthenticator, errors.Error) {
	// cache the ApiAuthenticator for performance
	if ma.apiAuthenticator != nil {
		return ma.apiAuthenticator, nil
	}
	// cache missed
	switch ma.AuthMethod {
	case plugin.AUTH_METHOD_BASIC:
		basicAuth, ok := connection.(plugin.BasicAuthenticator)
		if !ok {
			return nil, errors.Default.New("connection doesn't support Basic Authentication")
		}
		ma.apiAuthenticator = basicAuth.GetBasicAuthenticator()
	case plugin.AUTH_METHOD_TOKEN:
		accessToken, ok := connection.(plugin.AccessTokenAuthenticator)
		if !ok {
			return nil, errors.Default.New("connection doesn't support AccessToken Authentication")
		}
		ma.apiAuthenticator = accessToken.GetAccessTokenAuthenticator()
	case plugin.AUTH_METHOD_APPKEY:
		// Note that AppKey Authentication requires complex logic like signing the request with timestamp
		// so, there is no way to solve them once and for all, each Specific Connection should implement
		// on its own.
		appKey, ok := connection.(plugin.AppKeyAuthenticator)
		if !ok {
			return nil, errors.Default.New("connection doesn't support AppKey Authentication")
		}
		// check ae/models/connection.go:AeAppKey if you needed an example
		ma.apiAuthenticator = appKey.GetAppKeyAuthenticator()
	default:
		return nil, errors.Default.New("no Authentication Method was specified")
	}
	return ma.apiAuthenticator, nil
}

// SetupAuthenticationForConnection sets up authentication for the specified `req` based on connection
// Specific Connection should implement IAuthentication and then call this method for MultiAuth to work properly,
// check jira/models/connection.go:JiraConn if you needed an example
// Note: this method would be called for each request, so it is performance-sensitive, do NOT use reflection here
func (ma *MultiAuth) SetupAuthenticationForConnection(connection plugin.ApiConnection, req *http.Request) errors.Error {
	apiAuthenticator, err := ma.GetApiAuthenticator(connection)
	if err != nil {
		return err
	}
	return apiAuthenticator.SetupAuthentication(req)
}

func (ma *MultiAuth) ValidateConnection(connection interface{}, v *validator.Validate) errors.Error {
	// the idea is to filtered out errors from unselected Authentication struct
	validationErrors := v.Struct(connection).(validator.ValidationErrors)
	if validationErrors != nil {
		filteredValidationErrors := make(validator.ValidationErrors, 0)
		for _, e := range validationErrors {
			// JiraConnection.JiraConn.BasicAuth.Username
			ns := strings.Split(e.Namespace(), ".")
			if len(ns) > 1 {
				// BasicAuth
				authName := ns[len(ns)-2]
				if plugin.ALL_AUTH[authName] && authName != ma.AuthMethod {
					continue
				}
				filteredValidationErrors = append(filteredValidationErrors, e)
			}
		}
		if len(filteredValidationErrors) > 0 {
			return errors.BadInput.Wrap(filteredValidationErrors, "validation failed")
		}
	}
	return nil
}
