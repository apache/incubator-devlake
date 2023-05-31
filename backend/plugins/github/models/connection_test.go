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
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestAccessTokenSuccess(t *testing.T) {
	accessToken := &GithubConnection{}
	accessToken.AuthMethod = "AccessToken"
	accessToken.Endpoint = "https://api.github.com/"
	accessToken.Name = "test"
	accessToken.Token = "some_token"
	err := accessToken.ValidateConnection(accessToken, validator.New())
	assert.NoError(t, err)
}

func TestAccessTokenFail(t *testing.T) {
	accessToken := &GithubConnection{}
	accessToken.AuthMethod = "AccessToken"
	accessToken.Endpoint = "https://api.github.com/"
	accessToken.Name = "test"
	accessToken.Token = ""
	err := accessToken.ValidateConnection(accessToken, validator.New())
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "Token")
	}
}

func TestAppKeyFail(t *testing.T) {
	appkey := &GithubConnection{}
	appkey.AuthMethod = "AppKey"
	appkey.Endpoint = "https://api.github.com/"
	appkey.Name = "test"
	err := appkey.ValidateConnection(appkey, validator.New())
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "InstallationID")
		assert.Contains(t, err.Error(), "AppId")
		assert.Contains(t, err.Error(), "SecretKey")
		println()
	}
}
