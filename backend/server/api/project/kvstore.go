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

package project

import (
	"fmt"
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/apache/incubator-devlake/server/services"

	"github.com/gin-gonic/gin"
)

// @Summary Get list of info by onboard
// @Description GET onboard info
// @Tags framework/projects
// @Param onboard path string true "onboard"
// @Success 200  {object} models.Kvstore
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /kvstore/{onboard} [get]
func GetKvstore(c *gin.Context) {
	storeKey := c.Param("onboard")
	result, err := services.GetKvstore(storeKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching %s data", storeKey)})
		return
	}

	shared.ApiOutputSuccess(c, result, http.StatusOK)
}

// @Summary Put a on board project
// @Description Put a board project
// @Tags framework/projects
// @Accept application/json
// @Param project body models.ApiKvstore false "json"
// @Success 200  {object} models.Kvstore
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /kvstore [PUT]
func PutKvstore(c *gin.Context) {
	var body models.Kvstore
	err := c.ShouldBind(&body)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}
	storeKey := body.StoreKey
	fmt.Println("storeKey: ", storeKey)

	onBoardOutput, err := services.PutKvstore(storeKey, &body)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error put on board project"))
		return
	}

	shared.ApiOutputSuccess(c, onBoardOutput, http.StatusCreated)
}
