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

package apikeyhelper

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
	"github.com/apache/incubator-devlake/core/utils"
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (
	EncodeKeyEnvStr = "ENCRYPTION_SECRET"
	apiKeyLen       = 128
)

type ApiKeyHelper struct {
	basicRes         context.BasicRes
	cfg              *viper.Viper
	logger           log.Logger
	encryptionSecret string
}

func NewApiKeyHelper(basicRes context.BasicRes, logger log.Logger) *ApiKeyHelper {
	cfg := config.GetConfig()
	encryptionSecret := strings.TrimSpace(cfg.GetString(EncodeKeyEnvStr))
	if encryptionSecret == "" {
		panic("ENCRYPTION_SECRET must be set in environment variable or .env file")
	}
	return &ApiKeyHelper{
		basicRes:         basicRes,
		cfg:              cfg,
		logger:           logger,
		encryptionSecret: encryptionSecret,
	}
}

func (c *ApiKeyHelper) Create(user *common.User, name string, expiredAt *time.Time, allowedPath string, apiKeyType string) (*models.ApiKey, errors.Error) {
	randomApiKey, hashedApiKey, generateApiKeyErr := c.digestApiKey()
	if generateApiKeyErr != nil {
		c.logger.Error(generateApiKeyErr, "digestApiKey")
		return nil, errors.Default.Wrap(generateApiKeyErr, "random letters")
	}

	// create transaction to update multiple tables
	var err errors.Error
	tx := c.basicRes.GetDal().Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			c.logger.Error(err, "panic occurs or something went wrong")
			if err = tx.Rollback(); err != nil {
				c.logger.Error(err, "CreateApiKey: failed to rollback")
			}
		}
	}()

	// create project first
	now := time.Now()
	apiKey := &models.ApiKey{
		Model: common.Model{
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:        name,
		ApiKey:      hashedApiKey,
		ExpiredAt:   expiredAt,
		AllowedPath: allowedPath,
		Type:        apiKeyType,
	}
	if user != nil {
		apiKey.Creator = common.Creator{
			Creator:      user.Name,
			CreatorEmail: user.Email,
		}
		apiKey.Updater = common.Updater{
			Updater:      user.Name,
			UpdaterEmail: user.Email,
		}
	}
	err = tx.Create(apiKey)
	if err != nil {
		if tx.IsDuplicationError(err) {
			return nil, errors.BadInput.New(fmt.Sprintf("An api key with name [%s] has already exists", name))
		}
		return nil, errors.Default.Wrap(err, "error creating DB api key")
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	apiKey.ApiKey = randomApiKey
	return apiKey, nil
}

func (c *ApiKeyHelper) Put(user *common.User, id uint64) (*models.ApiKey, errors.Error) {
	// verify exists
	db := c.basicRes.GetDal()
	apiKey, err := c.getApiKeyById(db, id)
	if err != nil {
		c.logger.Error(err, "get api key by id: %d", id)
		return nil, err
	}

	apiKeyStr, hashApiKey, err := c.digestApiKey()
	if err != nil {
		c.logger.Error(err, "digestApiKey")
		return nil, err
	}
	apiKey.ApiKey = hashApiKey
	apiKey.UpdatedAt = time.Now()
	if user != nil {
		apiKey.Updater = common.Updater{
			Updater:      user.Name,
			UpdaterEmail: user.Email,
		}
	}
	if err = db.Update(apiKey); err != nil {
		c.logger.Error(err, "delete api key, id: %d", id)
		return nil, errors.Default.Wrap(err, "error deleting project")
	}
	apiKey.ApiKey = apiKeyStr
	return apiKey, nil
}

func (c *ApiKeyHelper) getApiKeyById(tx dal.Dal, id uint64, additionalClauses ...dal.Clause) (*models.ApiKey, errors.Error) {
	apiKey := &models.ApiKey{}
	err := tx.First(apiKey, append([]dal.Clause{dal.Where("id = ?", id)}, additionalClauses...)...)
	if err != nil {
		if tx.IsErrorNotFound(err) {
			return nil, errors.NotFound.Wrap(err, fmt.Sprintf("could not find api key id[%d] in DB", id))
		}
		return nil, errors.Default.Wrap(err, "error getting api key from DB")
	}
	return apiKey, nil
}

func (c *ApiKeyHelper) digestApiKey() (string, string, errors.Error) {
	randomApiKey, randomLetterErr := utils.RandLetterBytes(apiKeyLen)
	if randomLetterErr != nil {
		return "", "", errors.Default.Wrap(randomLetterErr, "random letters")
	}
	hashedApiKey, err := c.DigestToken(randomApiKey)
	return randomApiKey, hashedApiKey, err
}

func (c *ApiKeyHelper) DigestToken(token string) (string, errors.Error) {
	h := hmac.New(sha256.New, []byte(c.encryptionSecret))
	if _, err := h.Write([]byte(token)); err != nil {
		c.logger.Error(err, "hmac write api key")
		return "", errors.Default.Wrap(err, "hmac write api key")
	}
	hashedApiKey := fmt.Sprintf("%x", h.Sum(nil))
	return hashedApiKey, nil
}

func (c *ApiKeyHelper) DeleteForPlugin(tx dal.Transaction, pluginName string, extra string) errors.Error {
	// delete api key generated by plugin, for example webhook
	var apiKey models.ApiKey
	var clauses []dal.Clause
	if pluginName != "" {
		clauses = append(clauses, dal.Where("type = ?", fmt.Sprintf("plugin:%s", pluginName)))
	}
	if extra != "" {
		clauses = append(clauses, dal.Where("extra = ?", extra))
	}
	if err := tx.First(&apiKey, clauses...); err != nil {
		c.logger.Error(err, "query api key record")
		// if api key doesn't exist, just return success
		if tx.IsErrorNotFound(err.Unwrap()) {
			return nil
		} else {
			return err
		}
	}
	if err := tx.Delete(apiKey); err != nil {
		c.logger.Error(err, "delete api key record")
		return err
	}
	return nil
}

func (c *ApiKeyHelper) CreateForPlugin(tx dal.Transaction, user *common.User, name string, pluginName string, allowedPath string, extra string) (*models.ApiKey, errors.Error) {
	apiKey, hashedApiKey, err := c.digestApiKey()
	if err != nil {
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
	if user != nil {
		apiKeyRecord.Creator = common.Creator{
			Creator:      user.Name,
			CreatorEmail: user.Email,
		}
		apiKeyRecord.Updater = common.Updater{
			Updater:      user.Name,
			UpdaterEmail: user.Email,
		}
	}
	if err := tx.Create(apiKeyRecord); err != nil {
		c.logger.Error(err, "create api key record")
		return nil, err
	}
	apiKeyRecord.ApiKey = apiKey
	return apiKeyRecord, nil
}
