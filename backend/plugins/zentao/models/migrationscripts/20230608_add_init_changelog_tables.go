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
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"github.com/apache/incubator-devlake/plugins/zentao/models/archived"
)

type addInitChangelogTables struct{}

// This object conforms to what the frontend currently sends.
type ZentaoConnection20230608 struct {
	DbUrl          string `mapstructure:"dbUrl"  json:"dbUrl" gorm:"serializer:encdec"`
	DbIdleConns    int    `json:"dbIdleConns" mapstructure:"dbIdleConns"`
	DbLoggingLevel string `json:"dbLoggingLevel" mapstructure:"dbLoggingLevel"`
	DbMaxConns     int    `json:"dbMaxConns" mapstructure:"dbMaxConns"`
}

func (ZentaoConnection20230608) TableName() string {
	return "_tool_zentao_connections"
}

func (*addInitChangelogTables) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	err := db.DropTables(
		&archived.ZentaoChangelog{},
		&archived.ZentaoChangelogDetail{},
	)
	if err != nil {
		return err
	}
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.ZentaoChangelog{},
		&archived.ZentaoChangelogDetail{},
		&ZentaoConnection20230608{},
	)
}

func (*addInitChangelogTables) Version() uint64 {
	return 20230608000001
}

func (*addInitChangelogTables) Name() string {
	return "zentao init changelog schemas"
}
