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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	common "github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/apikeyhelper"
	"github.com/apache/incubator-devlake/plugins/webhook/models/migrationscripts/archived"
)

var _ plugin.MigrationScript = (*addApiKeys)(nil)

const (
	creator    = "migration_scripts"
	pluginName = "webhook"
)

type addApiKeys struct{}

func (u *addApiKeys) Up(baseRes context.BasicRes) errors.Error {
	db := baseRes.GetDal()
	if !db.HasTable(archived.WebhookConnection{}.TableName()) {
		return nil
	}
	var webhooks []archived.WebhookConnection
	if err := db.All(&webhooks); err != nil {
		return err
	}
	user := &common.User{
		Name:  creator,
		Email: "",
	}
	logger := baseRes.GetLogger()
	apiKeyHelper := apikeyhelper.NewApiKeyHelper(baseRes, logger)
	tx := db.Begin()
	for _, webhook := range webhooks {
		name := apiKeyHelper.GenApiKeyNameForPlugin(pluginName, webhook.ID)
		_, err := apiKeyHelper.GetApiKey(db, dal.Where("name = ?", name))
		if err != nil {
			if db.IsErrorNotFound(err) {
				logger.Info("api key with name: %s not found in db", name)
				allowedPath := fmt.Sprintf("/plugins/%s/connections/%d/.*", pluginName, webhook.ID)
				extra := fmt.Sprintf("connectionId:%d", webhook.ID)
				_, err := apiKeyHelper.CreateForPlugin(tx, user, name, pluginName, allowedPath, extra)
				if err != nil {
					logger.Error(err, "CreateForPlugin")
					if err := tx.Rollback(); err != nil {
						logger.Error(err, "rollback transaction")
						return err
					}
					return err
				}
			} else {
				logger.Error(err, "query api key from db, name: %s", name)
				if err := tx.Rollback(); err != nil {
					logger.Error(err, "rollback transaction")
					return err
				}
				return err
			}
		}
	}
	if err := tx.Commit(); err != nil {
		logger.Error(err, "commit transaction")
		return err
	}
	return nil
}

func (*addApiKeys) Version() uint64 {
	return 20231116103000
}

func (*addApiKeys) Name() string {
	return "associate api keys for webhook record automatically"
}
