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
	"context"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/plugins/feishu/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type addInitTables struct {
	config core.ConfigGetter
}

func (u *addInitTables) SetConfigGetter(config core.ConfigGetter) {
	u.config = config
}

func (u *addInitTables) Up(ctx context.Context, db *gorm.DB) errors.Error {
	err := db.Migrator().DropTable(
		&archived.FeishuConnection{},
		&archived.FeishuMeetingTopUserItem{},
	)
	if err != nil {
		return errors.Convert(err)
	}
	err = db.Migrator().CreateTable(
		&archived.FeishuConnection{},
		&archived.FeishuMeetingTopUserItem{},
	)
	if err != nil {
		return errors.Convert(err)
	}

	encodeKey := u.config.GetString(core.EncodeKeyEnvStr)
	connection := &archived.FeishuConnection{}
	connection.Endpoint = u.config.GetString(`FEISHU_ENDPOINT`)
	connection.AppId = u.config.GetString(`FEISHU_APPID`)
	connection.SecretKey = u.config.GetString(`FEISHU_APPSCRECT`)
	connection.Name = `Feishu`
	if connection.Endpoint != `` && connection.AppId != `` && connection.SecretKey != `` && encodeKey != `` {
		err = helper.UpdateEncryptFields(connection, func(plaintext string) (string, errors.Error) {
			return core.Encrypt(encodeKey, plaintext)
		})
		if err != nil {
			return errors.Convert(err)
		}
		// update from .env and save to db
		db.Create(connection)
	}
	return nil
}

func (*addInitTables) Version() uint64 {
	return 20220714000001
}

func (*addInitTables) Name() string {
	return "Feishu init schemas"
}
