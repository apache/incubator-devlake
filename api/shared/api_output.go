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

package shared

import (
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/logger"
	"github.com/gin-gonic/gin"
)

type ApiBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ApiOutputError(c *gin.Context, err error, status int) {
	if e, ok := err.(*errors.Error); ok {
		c.JSON(e.Status, &ApiBody{
			Success: false,
			Message: err.Error(),
		})
	} else {
		logger.Global.Error("Server Internal Error: %s", err.Error())
		c.JSON(status, &ApiBody{
			Success: false,
			Message: err.Error(),
		})
	}
	c.Writer.Header().Set("Content-Type", "application/json")
}

func ApiOutputSuccess(c *gin.Context, body interface{}, status int) {
	if body == nil {
		body = &ApiBody{
			Success: true,
			Message: "success",
		}
	}
	c.JSON(status, body)
}
