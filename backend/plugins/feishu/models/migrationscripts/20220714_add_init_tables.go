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
	"github.com/apache/incubator-devlake/plugins/feishu/models/migrationscripts/archived"
)

type addInitTables struct {
}

func (u *addInitTables) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	err := db.DropTables(
		&archived.FeishuConnection{},
		&archived.FeishuMeetingTopUserItem{},
	)
	if err != nil {
		return err
	}

	err = migrationhelper.AutoMigrateTables(
		basicRes,
		&archived.FeishuConnection{},
		&archived.FeishuMeetingTopUserItem{},
	)
	if err != nil {
		return err
	}

	encodeKey := basicRes.GetConfig(plugin.EncodeKeyEnvStr)
	connection := &archived.FeishuConnection{}
	connection.Endpoint = basicRes.GetConfig(`FEISHU_ENDPOINT`)
	connection.AppId = basicRes.GetConfig(`FEISHU_APPID`)
	connection.SecretKey = basicRes.GetConfig(`FEISHU_APPSCRECT`)
	connection.Name = `Feishu`
	if connection.Endpoint != `` && connection.AppId != `` && connection.SecretKey != `` && encodeKey != `` {
		// update from .env and save to db
		err = db.CreateIfNotExist(connection)
		if err != nil {
			return err
		}
	}
	return nil
}

func (*addInitTables) Version() uint64 {
	return 20220714000001
}

func (*addInitTables) Name() string {
	return "Feishu init schemas"
}
