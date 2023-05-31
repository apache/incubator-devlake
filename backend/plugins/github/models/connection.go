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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
	"github.com/golang-jwt/jwt/v5"
)

// GithubAccessToken supports fetching data with multiple tokens
type GithubAccessToken struct {
	helper.AccessToken `mapstructure:",squash"`
	tokens             []string `gorm:"-" json:"-" mapstructure:"-"`
	tokenIndex         int      `gorm:"-" json:"-" mapstructure:"-"`
}

type GithubAppKey struct {
	helper.AppKey  `mapstructure:",squash"`
	InstallationID int `mapstructure:"installationId" validate:"required" json:"installationId"`
}

// GithubConn holds the essential information to connect to the Github API
type GithubConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.MultiAuth      `mapstructure:",squash"`
	GithubAccessToken     `mapstructure:",squash" authMethod:"AccessToken"`
	GithubAppKey          `mapstructure:",squash" authMethod:"AppKey"`
}

// PrepareApiClient splits Token to tokens for SetupAuthentication to utilize
func (conn *GithubConn) PrepareApiClient(apiClient apihelperabstract.ApiClientAbstract) errors.Error {

	if conn.AuthMethod == "AccessToken" {
		conn.tokens = strings.Split(conn.Token, ",")
	}

	if conn.AuthMethod == "AppKey" && conn.InstallationID != 0 {
		token, err := conn.getInstallationAccessToken(apiClient)
		if err != nil {
			return err
		}

		conn.Token = token.Token
		conn.tokens = []string{token.Token}
	}

	return nil
}

// SetupAuthentication sets up the HTTP Request Authentication
func (conn *GithubConn) SetupAuthentication(req *http.Request) errors.Error {
	// Rotates token on each request.
	if len(conn.tokens) > 0 {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", conn.tokens[conn.tokenIndex]))
		// Set next token index
		conn.tokenIndex = (conn.tokenIndex + 1) % len(conn.tokens)
	}

	return nil
}

func (gat *GithubAccessToken) GetTokensCount() int {
	return len(gat.tokens)
}

// GithubConnection holds GithubConn plus ID/Name for database storage
type GithubConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	GithubConn            `mapstructure:",squash"`
	EnableGraphql         bool `mapstructure:"enableGraphql" json:"enableGraphql"`
}

func (GithubConnection) TableName() string {
	return "_tool_github_connections"
}

// Using GithubUserOfToken because it requires authentication, and it is public information anyway.
type GithubUserOfToken struct {
	Login string `json:"login"`
}

type InstallationToken struct {
	Token string `json:"token"`
}

type GithubApp struct {
	ID   int32  `json:"id"`
	Slug string `json:"slug"`
}

type GithubAppInstallation struct {
	Id      int `json:"id"`
	Account struct {
		Login string `json:"login"`
	} `json:"account"`
}

type GithubAppInstallationWithToken struct {
	GithubAppInstallation
	Token string
}

func (gak *GithubAppKey) CreateJwt() (string, errors.Error) {
	token := jwt.New(jwt.SigningMethodRS256)
	t := time.Now().Unix()

	token.Claims = jwt.MapClaims{
		"iat": t,
		"exp": t + (10 * 60),
		"iss": gak.AppId,
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(gak.SecretKey))
	if err != nil {
		return "", errors.BadInput.Wrap(err, "invalid private key")
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.BadInput.Wrap(err, "invalid private key")
	}

	return tokenString, nil
}

func (gak *GithubAppKey) getInstallationAccessToken(
	apiClient apihelperabstract.ApiClientAbstract,
) (*InstallationToken, errors.Error) {

	jwt, err := gak.CreateJwt()
	if err != nil {
		return nil, err
	}

	resp, err := apiClient.Post(fmt.Sprintf("/app/installations/%d/access_tokens", gak.InstallationID), nil, nil, http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", jwt)},
	})
	if err != nil {
		return nil, err
	}

	body, err := errors.Convert01(io.ReadAll(resp.Body))
	if err != nil {
		return nil, err
	}

	var installationToken InstallationToken
	err = errors.Convert(json.Unmarshal(body, &installationToken))
	if err != nil {
		return nil, err
	}

	return &installationToken, nil
}
