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
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/models"
	common "github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	"github.com/apache/incubator-devlake/plugins/webhook/models/migrationscripts/archived"
	"regexp"
	"strings"
	"time"
)

var _ plugin.MigrationScript = (*addApiKeys)(nil)

const (
	creator         = "migration_scripts"
	pluginName      = "webhook"
	apiKeyLen       = 128
	EncodeKeyEnvStr = "ENCRYPTION_SECRET"
)

type addApiKeys struct{}

func generateApiKey(logger log.Logger) (apiKey string, hashedApiKey string, err errors.Error) {
	apiKey, randomLetterErr := utils.RandLetterBytes(apiKeyLen)
	if randomLetterErr != nil {
		err = errors.Default.Wrap(randomLetterErr, "random letters")
		return
	}
	hashedApiKey, err = digestToken(apiKey, logger)
	return apiKey, hashedApiKey, err
}

func digestToken(token string, logger log.Logger) (string, errors.Error) {
	cfg := config.GetConfig()
	encryptionSecret := strings.TrimSpace(cfg.GetString(EncodeKeyEnvStr))
	h := hmac.New(sha256.New, []byte(encryptionSecret))
	if _, err := h.Write([]byte(token)); err != nil {
		logger.Error(err, "hmac write api key")
		return "", errors.Default.Wrap(err, "hmac write token")
	}
	hashedApiKey := fmt.Sprintf("%x", h.Sum(nil))
	return hashedApiKey, nil
}

func createForPlugin(db dal.Dal, logger log.Logger, user string, name string, pluginName string, allowedPath string, extra string) (*models.ApiKey, errors.Error) {
	if _, err := regexp.Compile(allowedPath); err != nil {
		logger.Error(err, "Compile allowed path")
		return nil, errors.Default.Wrap(err, fmt.Sprintf("compile allowed path: %s", allowedPath))
	}
	apiKey, hashedApiKey, err := generateApiKey(logger)
	if err != nil {
		logger.Error(err, "generateApiKey")
		return nil, err
	}
	now := time.Now()
	apiKeyRecord := &models.ApiKey{
		Model: common.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        name,
		ApiKey:      hashedApiKey,
		ExpiredAt:   nil,
		AllowedPath: allowedPath,
		Type:        fmt.Sprintf("plugin:%s", pluginName),
		Extra:       extra,
	}
	if user != "" {
		apiKeyRecord.Creator = common.Creator{
			Creator:      user,
			CreatorEmail: "",
		}
		apiKeyRecord.Updater = common.Updater{
			Updater:      user,
			UpdaterEmail: "",
		}
	}
	if err := db.Create(apiKeyRecord); err != nil {
		logger.Error(err, "create api key record")
		if db.IsDuplicationError(err) {
			return nil, errors.BadInput.New(fmt.Sprintf("An api key with name [%s] has already exists", name))
		}
		return nil, errors.Default.Wrap(err, "error creating DB api key")
	}
	apiKeyRecord.ApiKey = apiKey
	return apiKeyRecord, nil
}

func (u *addApiKeys) Up(baseRes context.BasicRes) errors.Error {
	db := baseRes.GetDal()
	if !db.HasTable(archived.WebhookConnection{}.TableName()) {
		return nil
	}
	var webhooks []archived.WebhookConnection
	if err := db.All(&webhooks); err != nil {
		return err
	}
	logger := baseRes.GetLogger()
	for _, webhook := range webhooks {
		name := fmt.Sprintf("%s-%d", pluginName, webhook.ID)
		apiKey := &models.ApiKey{}
		if err := db.First(apiKey, dal.Where("name = ?", name)); err != nil {
			if db.IsErrorNotFound(err) {
				logger.Info("api key with name: %s not found in db", name)
				allowedPath := fmt.Sprintf("/plugins/%s/connections/%d/.*", pluginName, webhook.ID)
				extra := fmt.Sprintf("connectionId:%d", webhook.ID)
				_, err := createForPlugin(db, logger, creator, name, pluginName, allowedPath, extra)
				if err != nil {
					logger.Error(err, "CreateForPlugin")
					return err
				}
			} else {
				logger.Error(err, "query api key from db, name: %s", name)
				return err
			}
		}
	}
	return nil
}

func (*addApiKeys) Version() uint64 {
	return 20231116103000
}

func (*addApiKeys) Name() string {
	return "associate api keys for webhook record automatically"
}
