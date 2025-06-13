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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/teambition/models"
)

var _ plugin.MigrationScript = (*addAppIdBack)(nil)

type teambitionConnection20250529 struct {
	AppId     string
	SecretKey string `gorm:"serializer:encdec"`
}

func (teambitionConnection20250529) TableName() string {
	return "_tool_teambition_connections"
}

type addAppIdBack struct{}

func (*addAppIdBack) Up(basicRes context.BasicRes) errors.Error {
	basicRes.GetLogger().Warn(nil, "*********")
	migrationhelper.AutoMigrateTables(basicRes, &teambitionConnection20250529{})
	err := migrationhelper.AutoMigrateTables(basicRes, &models.TeambitionScopeConfig{})
	basicRes.GetLogger().Warn(err, "err scope")
	return err
}

func (*addAppIdBack) Version() uint64 {
	return 20250529165745
}

func (*addAppIdBack) Name() string {
	return "add app id back to teambition_connections"
}
