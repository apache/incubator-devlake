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

package migrationscripts

import (
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/core"
)

type githubRepoToBeTransform struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	GithubId     int    `gorm:"primaryKey"`
	Name         string `gorm:"type:varchar(255)"`
	OwnerLogin   string `json:"ownerLogin" gorm:"type:varchar(255)"`
}

type concatOwnerAndName struct{}

func (script *concatOwnerAndName) Up(basicRes core.BasicRes) errors.Error {
	err := migrationhelper.TransformColumns(
		basicRes,
		script,
		"_tool_github_repos",
		[]string{"name", "owner_login"},
		func(src *githubRepoToBeTransform) (*githubRepoToBeTransform, errors.Error) {
			if src.OwnerLogin != "" {
				src.Name = fmt.Sprintf("%s/%s", src.OwnerLogin, src.Name)
			}
			return &githubRepoToBeTransform{
				GithubId:     src.GithubId,
				ConnectionId: src.ConnectionId,
				Name:         src.Name,
				OwnerLogin:   "",
			}, nil
		})
	if err != nil {
		return err
	}
	return basicRes.GetDal().DropColumns("_tool_github_repos", "owner_login")
}

func (*concatOwnerAndName) Version() uint64 {
	return 20221214115900
}

func (*concatOwnerAndName) Name() string {
	return "concat owner and name for old github_repos and drop owner for github_repos"
}
