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
)

var _ plugin.MigrationScript = (*addIsBotToAccounts)(nil)

type account20260720 struct {
	IsBot bool `gorm:"default:false"`
}

func (account20260720) TableName() string {
	return "accounts"
}

type addIsBotToAccounts struct{}

func (*addIsBotToAccounts) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()
	if err := db.AutoMigrate(&account20260720{}); err != nil {
		return err
	}
	return nil
}

func (*addIsBotToAccounts) Version() uint64 {
	return 20260720120000
}

func (*addIsBotToAccounts) Name() string {
	return "add is_bot to accounts so bot/automation activity can be excluded from metrics, according to #8974"
}
