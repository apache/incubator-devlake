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

package apikeys

import (
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/apache/incubator-devlake/server/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PaginatedApiKeys struct {
	ApiKeys []*models.ApiKey `json:"apikeys"`
	Count   int64            `json:"count"`
}

// @Summary Get list of api keys
// @Description GET /api-keys?page=1&pageSize=10
// @Tags framework/api-keys
// @Param page query int true "query"
// @Param pageSize query int true "query"
// @Success 200  {object} PaginatedApiKeys
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /api-keys [get]
func GetApiKeys(c *gin.Context) {
	var query services.ApiKeysQuery
	err := c.ShouldBindQuery(&query)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}
	apiKeys, count, err := services.GetApiKeys(&query)
	if err != nil {
		shared.ApiOutputAbort(c, errors.Default.Wrap(err, "error getting api keys"))
		return
	}

	shared.ApiOutputSuccess(c, PaginatedApiKeys{
		ApiKeys: apiKeys,
		Count:   count,
	}, http.StatusOK)
}

// @Summary Delete an api key
// @Description Delete an api key
// @Tags framework/api-keys
// @Accept application/json
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /api-keys/:apiKeyId [delete]
func DeleteApiKey(c *gin.Context) {
	apiKeyId := c.Param("apiKeyId")
	id, err := strconv.ParseUint(apiKeyId, 10, 64)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, "bad apiKeyId format supplied"))
		return
	}
	err = services.DeleteApiKey(id)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error deleting api key"))
		return
	}
	shared.ApiOutputSuccess(c, nil, http.StatusOK)
}

// @Summary Refresh an api key
// @Description Refresh an api key
// @Tags framework/api-keys
// @Accept application/json
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /api-keys/:apiKeyId [put]
func PutApiKey(c *gin.Context) {
	apiKeyId := c.Param("apiKeyId")
	id, err := strconv.ParseUint(apiKeyId, 10, 64)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, "bad apiKeyId format supplied"))
		return
	}
	user, email, err := shared.GetUserInfo(c.Request)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting user info"))
		return
	}
	apiOutputApiKey, err := services.PutApiKey(&common.Updater{
		Updater:      user,
		UpdaterEmail: email,
	}, id)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error regenerate api key"))
		return
	}
	shared.ApiOutputSuccess(c, apiOutputApiKey, http.StatusOK)
}

// @Summary Create a new api key
// @Description Create a new api key
// @Tags framework/api-keys
// @Accept application/json
// @Param apikey body models.ApiInputApiKey true "json"
// @Success 200  {object} models.ApiOutputApiKey
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /api-keys [post]
func PostApiKey(c *gin.Context) {
	apiKeyInput := &models.ApiInputApiKey{}
	err := c.ShouldBind(apiKeyInput)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}
	user, email, err := shared.GetUserInfo(c.Request)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting user info"))
		return
	}
	apiKeyOutput, err := services.CreateApiKey(&common.Creator{
		Creator:      user,
		CreatorEmail: email,
	}, apiKeyInput)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error creating api key"))
		return
	}

	shared.ApiOutputSuccess(c, apiKeyOutput, http.StatusCreated)
}
