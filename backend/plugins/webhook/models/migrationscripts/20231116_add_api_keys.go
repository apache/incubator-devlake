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
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	archivedBase "github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/webhook/models/migrationscripts/archived"
	"math/big"
	"regexp"
	"strings"
	"time"
)

var _ plugin.MigrationScript = (*addApiKeys)(nil)

const (
	creator         = "migration_scripts"
	pluginName      = "webhook"
	apiKeyLen       = 128
	encodeKeyEnvStr = "ENCRYPTION_SECRET"
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type addApiKeys struct{}

func randLetterBytes(n int) (string, errors.Error) {
	if n < 0 {
		return "", errors.Default.New("n must be greater than 0")
	}
	ret := make([]byte, n)
	bi := big.NewInt(int64(len(letterBytes)))
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, bi)
		if err != nil {
			return "", errors.Convert(err)
		}
		ret[i] = letterBytes[num.Int64()]
	}

	return string(ret), nil
}

func generateApiKey() (apiKey string, hashedApiKey string, err errors.Error) {
	apiKey, randomLetterErr := randLetterBytes(apiKeyLen)
	if randomLetterErr != nil {
		err = errors.Default.Wrap(randomLetterErr, "random letters")
		return
	}
	hashedApiKey, err = digestToken(apiKey)
	return apiKey, hashedApiKey, err
}

func digestToken(token string) (string, errors.Error) {
	cfg := config.GetConfig()
	encryptionSecret := strings.TrimSpace(cfg.GetString(encodeKeyEnvStr))
	h := hmac.New(sha256.New, []byte(encryptionSecret))
	if _, err := h.Write([]byte(token)); err != nil {
		return "", errors.Default.Wrap(err, "hmac write token")
	}
	hashedApiKey := fmt.Sprintf("%x", h.Sum(nil))
	return hashedApiKey, nil
}

func createForPlugin(db dal.Dal, user string, name string, pluginName string, allowedPath string, extra string) errors.Error {
	if _, err := regexp.Compile(allowedPath); err != nil {
		return errors.Default.Wrap(err, fmt.Sprintf("compile allowed path: %s", allowedPath))
	}
	apiKey, hashedApiKey, err := generateApiKey()
	if err != nil {
		return err
	}

	now := time.Now()
	apiKeyRecord := &archived.ApiKey{
		Model: archivedBase.Model{
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
		apiKeyRecord.Creator = archivedBase.Creator{
			Creator:      user,
			CreatorEmail: "",
		}
		apiKeyRecord.Updater = archivedBase.Updater{
			Updater:      user,
			UpdaterEmail: "",
		}
	}
	if err := db.Create(apiKeyRecord); err != nil {
		if db.IsDuplicationError(err) {
			return errors.BadInput.New(fmt.Sprintf("An api key with name [%s] has already exists", name))
		}
		return errors.Default.Wrap(err, "error creating DB api key")
	}
	apiKeyRecord.ApiKey = apiKey
	return nil
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
	for _, webhook := range webhooks {
		name := fmt.Sprintf("%s-%d", pluginName, webhook.ID)
		apiKey := &archived.ApiKey{}
		if err := db.First(apiKey, dal.Where("name = ?", name)); err != nil {
			if db.IsErrorNotFound(err) {
				allowedPath := fmt.Sprintf("/plugins/%s/connections/%d/.*", pluginName, webhook.ID)
				extra := fmt.Sprintf("connectionId:%d", webhook.ID)
				if err := createForPlugin(db, creator, name, pluginName, allowedPath, extra); err != nil {
					return err
				}
			} else {
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
