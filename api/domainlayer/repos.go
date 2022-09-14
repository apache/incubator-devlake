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

package domainlayer

import (
	"github.com/apache/incubator-devlake/errors"
	"net/http"

	"github.com/apache/incubator-devlake/api/shared"
	"github.com/apache/incubator-devlake/services"
	"github.com/gin-gonic/gin"
)

/*
Get all repos from database
GET /repos
{
	"repos": [
		{"id": "github:GithubRepo:384111310", "name": "merico-dev/lake", ...}
	],
	"count": 5
}
*/
// @Summary Get all repos from database
// @Description Get all repos from database
// @Tags framework/domainlayer
// @Accept application/json
// @Success 200  {object} gin.H "{"repos": repos, "count": count}"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /domainlayer/repos [get]
func ReposIndex(c *gin.Context) {
	repos, count, err := services.GetRepos()
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting repositories"))
		return
	}
	shared.ApiOutputSuccess(c, gin.H{"repos": repos, "count": count}, http.StatusOK)
}
