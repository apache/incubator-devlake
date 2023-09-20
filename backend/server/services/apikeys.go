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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/helpers/apikeyhelper"
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
	for idx := range apiKeys {
		apiKeys[idx].RemoveHashedApiKey() // delete the hashed api key to reduce the attack surface.
	}

	return apiKeys, count, nil
}

func DeleteApiKey(id uint64) errors.Error {
	// verify input
	if id == 0 {
		return errors.BadInput.New("api key's id is missing")
	}

	apiKeyHelper := apikeyhelper.NewApiKeyHelper(basicRes, logger)
	err := apiKeyHelper.Delete(id)
	if err != nil {
		logger.Error(err, "api key helper delete: %d", id)
		return err
	}
	return nil
}

func PutApiKey(user *common.User, id uint64) (*models.ApiOutputApiKey, errors.Error) {
	// verify input
	if id == 0 {
		return nil, errors.BadInput.New("api key's id is missing")
	}
	apiKeyHelper := apikeyhelper.NewApiKeyHelper(basicRes, logger)
	apiKey, err := apiKeyHelper.Put(user, id)
	if err != nil {
		logger.Error(err, "api key helper put: %d", id)
		return nil, err
	}
	return apiKey, nil
}

// CreateApiKey accepts an api key instance and insert it to database
func CreateApiKey(user *common.User, apiKeyInput *models.ApiInputApiKey) (*models.ApiOutputApiKey, errors.Error) {
	// verify input
	if err := VerifyStruct(apiKeyInput); err != nil {
		logger.Error(err, "verify: %+v", apiKeyInput)
		return nil, err
	}

	apiKeyHelper := apikeyhelper.NewApiKeyHelper(basicRes, logger)
	tx := basicRes.GetDal().Begin()
	apiKey, err := apiKeyHelper.Create(tx, user, apiKeyInput.Name, apiKeyInput.ExpiredAt, apiKeyInput.AllowedPath, apiKeyInput.Type, "")
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.Error(err, "transaction Rollback")
		}
		logger.Error(err, "api key helper create")
		return nil, errors.Default.Wrap(err, "random letters")
	}
	if err := tx.Commit(); err != nil {
		logger.Info("transaction commit: %s", err)
	}
	return apiKey, nil
}
