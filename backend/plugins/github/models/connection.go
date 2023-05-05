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
	"github.com/apache/incubator-devlake/core/plugin"
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
	helper.AppKey `mapstructure:",squash"`

	appJsonWebToken        string                           `gorm:"-" json:"-" mapstructure:"-"`
	selectedRepository     string                           `gorm:"-" json:"-" mapstructure:"-"`
	selectedInstallationID *int                             `gorm:"-" json:"-" mapstructure:"-"`
	installations          []GithubAppInstallationWithToken `gorm:"-" json:"-" mapstructure:"-"`
}

// GithubConn holds the essential information to connect to the Github API
type GithubConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.MultiAuth      `mapstructure:",squash"`
	GithubAccessToken     `mapstructure:",squash"`
	GithubAppKey          `mapstructure:",squash"`
}

func (conn *GithubConn) GetToken() string {
	if conn.AuthMethod == "AccessToken" {
		return strings.Split(conn.Token, ",")[0]
	}

	if conn.AuthMethod == "AppKey" {
		return conn.GithubAppKey.GetToken()
	}

	return ""
}

func (conn *GithubConn) SetRepository(repo string) {
	conn.GithubAppKey.selectedRepository = repo
	conn.GithubAppKey.selectedInstallationID = nil
}

func (conn *GithubConn) SetInstallationID(id int) {
	conn.GithubAppKey.selectedInstallationID = &id
	conn.GithubAppKey.selectedRepository = ""
}

// PrepareApiClient splits Token to tokens for SetupAuthentication to utilize
func (conn *GithubConn) PrepareApiClient(apiClient apihelperabstract.ApiClientAbstract) errors.Error {

	if conn.AuthMethod == "AccessToken" {
		conn.tokens = strings.Split(conn.Token, ",")
	}

	if conn.AuthMethod == "AppKey" {
		appToken, err := conn.GithubAppKey.createJwt()
		if err != nil {
			return err
		}

		conn.appJsonWebToken = appToken

		installations, err := conn.GithubAppKey.listAppInstallations(apiClient)
		if err != nil {
			return err
		}

		var installationTokens []GithubAppInstallationWithToken
		for _, installation := range installations {
			installationToken, err := conn.GithubAppKey.getInstallationAccessToken(installation.Id, apiClient)
			if err != nil {
				return err
			}
			installationTokens = append(installationTokens, GithubAppInstallationWithToken{
				GithubAppInstallation: installation,
				Token:                 installationToken.Token,
			})
		}
		conn.installations = installationTokens
	}

	return nil
}

// SetupAuthentication sets up the HTTP Request Authentication
func (conn *GithubConn) SetupAuthentication(req *http.Request) errors.Error {
	return conn.MultiAuth.SetupAuthenticationForConnection(conn, req)
}

func (gat *GithubAccessToken) SetupAuthentication(req *http.Request) errors.Error {
	// Rotates token on each request.
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", gat.tokens[gat.tokenIndex]))
	// Set next token index
	gat.tokenIndex = (gat.tokenIndex + 1) % len(gat.tokens)

	return nil
}

func (gak *GithubAppKey) SetupAuthentication(req *http.Request) errors.Error {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", gak.GetToken()))

	return nil
}

func (gak *GithubAppKey) GetAppKeyAuthenticator() plugin.ApiAuthenticator {
	return gak
}

func (gat *GithubAppKey) GetTokensCount() int {
	return 1
}

func (gak *GithubAppKey) GetToken() string {
	for _, installation := range gak.installations {
		if gak.selectedInstallationID != nil && installation.Id == *gak.selectedInstallationID {
			return installation.Token
		}
		if gak.selectedRepository != "" {
			owner := strings.Split(gak.selectedRepository, "/")[0]
			if installation.Account.Login == owner {
				return installation.Token
			}
		}
	}

	return gak.appJsonWebToken
}

func (gat *GithubAccessToken) GetAccessTokenAuthenticator() plugin.ApiAuthenticator {
	return gat
}

func (gat *GithubAccessToken) GetTokensCount() int {
	return len(gat.tokens)
}

// GetTokensCount returns total number of tokens
func (conn *GithubConn) GetTokensCount() int {
	if conn.AuthMethod == "AccessToken" {
		return conn.GithubAccessToken.GetTokensCount()
	}

	if conn.AuthMethod == "AppKey" {
		return conn.GithubAppKey.GetTokensCount()
	}

	return 0
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

func (gak *GithubAppKey) createJwt() (string, errors.Error) {
	token := jwt.New(jwt.SigningMethodRS256)
	t := time.Now().Unix()

	token.Claims = jwt.MapClaims{
		"iat": t,
		"exp": t + (10 * 60),
		"iss": gak.AppId,
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(gak.SecretKey))
	if err != nil {
		return "", errors.AsLakeErrorType(err)
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", errors.AsLakeErrorType(err)
	}

	return tokenString, nil
}

func (gak *GithubAppKey) listAppInstallations(
	apiClient apihelperabstract.ApiClientAbstract,
) ([]GithubAppInstallation, errors.Error) {
	installationsRes := []GithubAppInstallation{}

	res, err := apiClient.Get("app/installations", nil, http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", gak.appJsonWebToken)},
	})
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code when requesting app installation from %s", res.Request.URL.String()))
	}
	body, err := errors.Convert01(io.ReadAll(res.Body))
	if err != nil {
		return nil, err
	}
	err = errors.Convert(json.Unmarshal(body, &installationsRes))
	if err != nil {
		return nil, err
	}
	return installationsRes, nil
}

func (gak *GithubAppKey) getInstallationAccessToken(
	installationID int,
	apiClient apihelperabstract.ApiClientAbstract,
) (*InstallationToken, errors.Error) {

	resp, err := apiClient.Post(fmt.Sprintf("/app/installations/%d/access_tokens", installationID), nil, nil, http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", gak.appJsonWebToken)},
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
