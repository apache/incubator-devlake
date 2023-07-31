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

package api

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/helpers/apikeyhelper"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func Authentication(router *gin.Engine, basicRes context.BasicRes) gin.HandlerFunc {
	type ApiBody struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	db := basicRes.GetDal()
	logger := basicRes.GetLogger()
	if db == nil {
		panic(fmt.Errorf("db is not initialised"))
	}
	apiKeyHelper := apikeyhelper.NewApiKeyHelper(basicRes, logger)
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") {
			path = strings.TrimPrefix(path, "/api")
		}

		// Only open api needs to check api key
		if !strings.HasPrefix(path, "/rest") {
			logger.Info("path %s will continue", path)
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Abort()
			c.JSON(http.StatusUnauthorized, &ApiBody{
				Success: false,
				Message: "token is missing",
			})
			return
		}
		apiKeyStr := strings.TrimPrefix(authHeader, "Bearer ")
		if apiKeyStr == authHeader || apiKeyStr == "" {
			c.Abort()
			c.JSON(http.StatusUnauthorized, &ApiBody{
				Success: false,
				Message: "token is not present or malformed",
			})
			return
		}

		hashedApiKey, err := apiKeyHelper.GenerateApiKeyWithToken(apiKeyStr)
		if err != nil {
			logger.Error(err, "GenerateApiKeyWithToken")
			c.Abort()
			c.JSON(http.StatusInternalServerError, &ApiBody{
				Success: false,
				Message: err.Error(),
			})
			return
		}

		var apiKey models.ApiKey
		err = db.First(&apiKey, dal.Where("api_key = ?", hashedApiKey))
		if err != nil {
			c.Abort()
			if db.IsErrorNotFound(err) {
				c.JSON(http.StatusForbidden, &ApiBody{
					Success: false,
					Message: "api key is invalid",
				})
			} else {
				logger.Error(err, "query api key from db")
				c.JSON(http.StatusInternalServerError, &ApiBody{
					Success: false,
					Message: err.Error(),
				})
			}
			return
		}

		if apiKey.ExpiredAt != nil && apiKey.ExpiredAt.Sub(time.Now()) < 0 {
			c.Abort()
			c.JSON(http.StatusForbidden, &ApiBody{
				Success: false,
				Message: "api key has expired",
			})
			return
		}
		matched, matchErr := regexp.MatchString(apiKey.AllowedPath, path)
		if matchErr != nil {
			logger.Error(err, "regexp match path error")
			c.Abort()
			c.JSON(http.StatusInternalServerError, &ApiBody{
				Success: false,
				Message: matchErr.Error(),
			})
			return
		}
		if !matched {
			c.JSON(http.StatusForbidden, &ApiBody{
				Success: false,
				Message: "path doesn't match api key's scope",
			})
			return
		}

		if strings.HasPrefix(path, "/rest") {
			logger.Info("redirect path: %s to: %s", path, strings.TrimPrefix(path, "/rest"))
			c.Request.URL.Path = strings.TrimPrefix(path, "/rest")
		}
		router.HandleContext(c)
		c.Abort()
	}
}
