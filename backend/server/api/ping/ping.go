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

package ping

import (
	"net/http"

	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/apache/incubator-devlake/server/services"
	"github.com/gin-gonic/gin"
)

// @Summary Ping
// @Description check http status
// @Tags framework/ping
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /ping [get]
func Get(c *gin.Context) {
	c.Status(http.StatusOK)
}

// @Summary Ready
// @Description check if service is ready
// @Tags framework/ping
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /ready [get]
func Ready(c *gin.Context) {
	status, err := services.Ready()
	if err != nil {
		shared.ApiOutputError(c, err)
		return
	}
	shared.ApiOutputSuccess(c, shared.ApiBody{Success: true, Message: status}, http.StatusOK)
}

// @Summary Health
// @Description check if service is health
// @Tags framework/ping
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /health [get]
func Health(c *gin.Context) {
	msg, err := services.Health()
	if err != nil {
		shared.ApiOutputError(c, err)
		return
	}
	shared.ApiOutputSuccess(c, shared.ApiBody{Success: true, Message: msg}, http.StatusOK)
}
