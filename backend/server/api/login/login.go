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

package login

import (
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/apache/incubator-devlake/server/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AuthenticationResult AuthenticationResult `json:"AuthenticationResult"`
	ChallengeName        interface{}          `json:"ChallengeName"`
	ChallengeParameters  ChallengeParameters  `json:"ChallengeParameters"`
	Session              interface{}          `json:"Session"`
}
type AuthenticationResult struct {
	AccessToken       string      `json:"AccessToken"`
	ExpiresIn         int         `json:"ExpiresIn"`
	IDToken           string      `json:"IdToken"`
	NewDeviceMetadata interface{} `json:"NewDeviceMetadata"`
	RefreshToken      string      `json:"RefreshToken"`
	TokenType         string      `json:"TokenType"`
}
type ChallengeParameters struct {
}

// @Summary post login
// @Description post login
// @Tags framework/login
// @Accept application/json
// @Param blueprint body LoginRequest true "json"
// @Success 200  {object} LoginResponse
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /login [post]
func Login(ctx *gin.Context) {
	loginReq := &LoginRequest{}
	err := ctx.ShouldBind(loginReq)
	if err != nil {
		shared.ApiOutputError(ctx, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}
	res, err := auth.SignIn(auth.CreateCognitoClient(), loginReq.Username, loginReq.Password)
	if err != nil {
		shared.ApiOutputError(ctx, errors.Default.Wrap(err, "error signing in"))
		return
	}
	shared.ApiOutputSuccess(ctx, res, http.StatusOK)
}
