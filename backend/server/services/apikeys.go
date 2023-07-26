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

package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/utils"
	"time"
)

// ApiKeysQuery used to query api keys as the api key input
type ApiKeysQuery struct {
	Pagination
}

// GetApiKeys returns a paginated list of api keys based on `query`
func GetApiKeys(query *ApiKeysQuery) ([]*models.ApiKey, int64, errors.Error) {
	// verify input
	if err := VerifyStruct(query); err != nil {
		return nil, 0, err
	}
	clauses := []dal.Clause{
		dal.From(&models.ApiKey{}),
		dal.Where("type = ?", "devlake"),
	}

	logger.Info("query: %+v", query)
	count, err := db.Count(clauses...)
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error getting DB count of api key")
	}

	clauses = append(clauses,
		dal.Orderby("created_at DESC"),
		dal.Offset(query.GetSkip()),
		dal.Limit(query.GetPageSize()),
	)
	apiKeys := make([]*models.ApiKey, 0)
	err = db.All(&apiKeys, clauses...)
	if err != nil {
		return nil, 0, errors.Default.Wrap(err, "error finding DB api key")
	}

	return apiKeys, count, nil
}

func getApiKeyById(tx dal.Dal, id uint64, additionalClauses ...dal.Clause) (*models.ApiKey, errors.Error) {
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

func DeleteApiKey(id uint64) errors.Error {
	// verify input
	if id == 0 {
		return errors.BadInput.New("api key's id is missing")
	}
	// verify exists
	_, err := getApiKeyById(db, id)
	if err != nil {
		logger.Error(err, "get api key by id: %d", id)
		return err
	}
	err = db.Delete(&models.ApiKey{}, dal.Where("id = ?", id))
	if err != nil {
		logger.Error(err, "delete api key, id: %d", id)
		return errors.Default.Wrap(err, "error deleting project")
	}
	return nil
}

func PutApiKey(id uint64) (*models.ApiOutputApiKey, errors.Error) {
	// verify input
	if id == 0 {
		return nil, errors.BadInput.New("api key's id is missing")
	}
	// verify exists
	apiKey, err := getApiKeyById(db, id)
	if err != nil {
		logger.Error(err, "get api key by id: %d", id)
		return nil, err
	}
	apiKeyStr, hashApiKey, err := utils.GenerateApiKey(context.Background())
	if err != nil {
		logger.Error(err, "GenerateApiKey")
		return nil, err
	}
	apiKey.ApiKey = hashApiKey
	apiKey.UpdatedAt = time.Now()
	apiKey.Updater = common.Updater{
		Updater:      "",
		UpdaterEmail: "",
	}
	if err = db.Update(apiKey); err != nil {
		logger.Error(err, "delete api key, id: %d", id)
		return nil, errors.Default.Wrap(err, "error deleting project")
	}

	apiOutputApiKey := &models.ApiOutputApiKey{ApiKey: *apiKey}
	apiOutputApiKey.ApiKey.ApiKey = apiKeyStr
	return apiOutputApiKey, nil
}

// CreateApiKey accepts an api key instance and insert it to database
func CreateApiKey(apiKeyInput *models.ApiInputApiKey) (*models.ApiOutputApiKey, errors.Error) {
	logger.Info("api key: %+v", apiKeyInput)

	// verify input
	if err := VerifyStruct(apiKeyInput); err != nil {
		logger.Error(err, "verify: %+v", apiKeyInput)
		return nil, err
	}

	apiKeyLen := 128 // notice
	randomApiKey, randomLetterErr := utils.RandLetterBytes(apiKeyLen)
	if randomLetterErr != nil {
		logger.Error(randomLetterErr, "RandLetterBytes with length: %d", apiKeyLen)
		return nil, errors.Default.Wrap(randomLetterErr, "random letters")
	}
	encodeKeyEnv, exist := utils.GetEncodeKeyEnv(context.Background())
	if !exist {
		err := errors.Default.New("encode key env doesn't exist")
		logger.Warn(err, "GetEncodeKeyEnv error")
		return nil, err
	}
	h := hmac.New(sha256.New, []byte(encodeKeyEnv))
	if _, err := h.Write([]byte(randomApiKey)); err != nil {
		logger.Error(err, "hmac write %s", randomApiKey)
		return nil, errors.Default.Wrap(err, "hmac write api key")
	}
	hashedApiKey := fmt.Sprintf("%x", h.Sum(nil))

	// create transaction to update multiple tables
	var err errors.Error
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			err = tx.Rollback()
			if err != nil {
				logger.Error(err, "CreateApiKey: failed to rollback")
			}
		}
	}()

	extra, jsonMarshalErr := json.Marshal(apiKeyInput.Extra)
	if jsonMarshalErr != nil {
		logger.Error(jsonMarshalErr, "json marshal %s", apiKeyInput.Extra)
		return nil, errors.Default.Wrap(err, "json marshal")
	}
	// create project first
	apiKey := &models.ApiKey{
		Model: common.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Creator: common.Creator{
			Creator:      "", // fixme
			CreatorEmail: "", // fixme
		},
		Updater: common.Updater{
			Updater:      "", // fixme
			UpdaterEmail: "", // fixme
		},
		Name:        apiKeyInput.Name,
		ApiKey:      hashedApiKey,
		ExpiredAt:   apiKeyInput.ExpiredAt,
		AllowedPath: apiKeyInput.AllowedPath,
		Type:        apiKeyInput.Type,
		Extra:       extra,
	}
	err = db.Create(apiKey)
	if err != nil {
		if db.IsDuplicationError(err) {
			return nil, errors.BadInput.New(fmt.Sprintf("An api key with name [%s] has already exists", apiKeyInput.Name))
		}
		return nil, errors.Default.Wrap(err, "error creating DB api key")
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	apiOutputApiKey := &models.ApiOutputApiKey{
		ApiKey: *apiKey,
	}
	apiKey.ApiKey = randomApiKey
	return apiOutputApiKey, nil
}
