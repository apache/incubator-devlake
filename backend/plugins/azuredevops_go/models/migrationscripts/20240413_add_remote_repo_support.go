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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type extendRepoTable struct{}

type ExtendAzuredevopsRepo struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	Id           string `gorm:"primaryKey;column:id"`

	Type       string `gorm:"type:varchar(100)"`
	IsPrivate  bool
	ExternalId string `gorm:"type:varchar(255)"`
}

func (ExtendAzuredevopsRepo) TableName() string {
	return "_tool_azuredevops_go_repos"
}

func (*extendRepoTable) Up(baseRes context.BasicRes) errors.Error {
	err := migrationhelper.AutoMigrateTables(baseRes, &ExtendAzuredevopsRepo{})
	if err != nil {
		return err
	}

	err = baseRes.GetDal().UpdateColumn(
		&ExtendAzuredevopsRepo{}, "type", "TfsGit",
		dal.Where("type IS NULL"))

	if err != nil {
		return err
	}

	return baseRes.GetDal().UpdateColumn(
		&ExtendAzuredevopsRepo{}, "is_private", false,
		dal.Where("is_private IS NULL"))

}

func (*extendRepoTable) Version() uint64 {
	return 20240413100000
}

func (*extendRepoTable) Name() string {
	return "add [type, is_private, external_id] to _tool_azuredevops_go_repos in order to support remote repositories"
}
