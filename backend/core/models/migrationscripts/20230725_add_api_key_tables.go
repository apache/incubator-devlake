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
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
	"time"
)

var _ plugin.MigrationScript = (*addApiKeyTables)(nil)

type addApiKeyTables struct{}

type apiKey20230728 struct {
	archived.Model
	archived.Creator
	archived.Updater
	Name        string    `json:"name" gorm:"type:varchar(255);uniqueIndex"`
	ApiKey      string    `json:"apiKey" gorm:"type:varchar(255);column:api_key;uniqueIndex"`
	ExpiredAt   time.Time `json:"expiredAt" gorm:"column:expired_at"`
	AllowedPath string    `json:"allowedPath" gorm:"type:varchar(255);column:allowed_path"`
	Type        string    `json:"type" gorm:"type:varchar(40);column:type;index"`
	Extra       string    `json:"extra" gorm:"type:varchar(255);column:extra;index"`
}

func (apiKey20230728) TableName() string {
	return "_devlake_api_keys"
}

func (script *addApiKeyTables) Up(basicRes context.BasicRes) errors.Error {
	// To create multiple tables with migration helper
	return migrationhelper.AutoMigrateTables(
		basicRes,
		&apiKey20230728{},
	)
}

func (*addApiKeyTables) Version() uint64 {
	return 20230725142900
}

func (*addApiKeyTables) Name() string {
	return "add api key tables"
}
