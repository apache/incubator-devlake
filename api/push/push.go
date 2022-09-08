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

package push

import (
	"fmt"
	"net/http"

	"github.com/apache/incubator-devlake/api/shared"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/services"
	"github.com/gin-gonic/gin"
)

/*
	POST /push/:tableName
	[
		{
			"id": 1,
			"sha": "osidjfoawehfwh08"
		}
	]
*/
// @Summary POST /push/:tableName
// @Description POST /push/:tableName
// @Tags framework/push
// @Accept application/json
// @Param tableName path string true "table name"
// @Param data body string true "data"
// @Success 200  {object} gin.H "{"rowsAffected": rowsAffected}"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /push/{tableName} [post]
func Post(c *gin.Context) {
	var err error
	tableName := c.Param("tableName")
	var rowsToInsert []map[string]interface{}
	err = c.ShouldBindJSON(&rowsToInsert)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody, errors.AsUserMessage()))
		return
	}
	rowsAffected, err := services.InsertRow(tableName, rowsToInsert)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, fmt.Sprintf("error inserting request body into table %s", tableName), errors.AsUserMessage()))
		return
	}
	shared.ApiOutputSuccess(c, gin.H{"rowsAffected": rowsAffected}, http.StatusOK)
}
