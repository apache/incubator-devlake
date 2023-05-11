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
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/apache/incubator-devlake/server/services/auth"

	"github.com/gin-gonic/gin"
)

// @Summary post login
// @Description post login
// @Tags framework/login
// @Accept application/json
// @Param login body auth.LoginRequest true "json"
// @Success 200  {object} LoginResponse
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /login [post]
func Login(ctx *gin.Context) {
	loginReq := &auth.LoginRequest{}
	err := ctx.ShouldBind(loginReq)
	if err != nil {
		shared.ApiOutputError(ctx, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}
	res, err := auth.Provider.SignIn(loginReq)
	if err != nil {
		shared.ApiOutputError(ctx, errors.Default.Wrap(err, "error signing in"))
		return
	}
	if res.AuthenticationResult != nil && res.AuthenticationResult.AccessToken != nil {
		token, err := auth.Provider.CheckAuth(*res.AuthenticationResult.AccessToken)
		if err != nil {
			shared.ApiOutputAbort(ctx, err)
		}
		ctx.Set("token", token)
	}
	shared.ApiOutputSuccess(ctx, res, http.StatusOK)
}

// @Summary post NewPassword
// @Description post NewPassword
// @Tags framework/NewPassword
// @Accept application/json
// @Param newpassword body auth.NewPasswordRequest true "json"
// @Success 200  {object} shared.ApiBody
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /password [post]
func NewPassword(ctx *gin.Context) {
	newPasswordReq := &auth.NewPasswordRequest{}
	err := ctx.ShouldBind(newPasswordReq)
	if err != nil {
		shared.ApiOutputError(ctx, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}
	res, err := auth.Provider.NewPassword(newPasswordReq)
	if err != nil {
		shared.ApiOutputError(ctx, errors.BadInput.Wrap(err, "failed to set new password"))
		return
	}
	shared.ApiOutputSuccess(ctx, res, http.StatusOK)
}
