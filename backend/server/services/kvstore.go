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
	"fmt"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
)

// GetProjects returns a paginated list of Projects based on `query`
func GetKvstore(storeKey string) (*models.Kvstore, errors.Error) {
	clauses := []dal.Clause{
		dal.From(&models.Kvstore{}),
		dal.Where("store_key = ?", storeKey),
	}

	kvstore := &models.Kvstore{}
	err := db.First(&kvstore, clauses...)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error finding %s on _devlake_kvstore table", storeKey))
	}

	return kvstore, nil
}

// PutOnboard accepts a project instance and insert it to database
func PutKvstore(storeKey string, storeValue *models.Kvstore) (*models.Kvstore, errors.Error) {
	// verify input
	if err := VerifyStruct(storeValue); err != nil {
		return nil, err
	}

	// create transaction to updte multiple tables
	var err errors.Error
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			err = tx.Rollback()
			if err != nil {
				logger.Error(err, "PutKvstore: failed to rollback")
			}
		}
	}()

	result := &models.Kvstore{}
	err = tx.CreateOrUpdate(storeValue, dal.Where("store_key = ?", storeKey))
	if err != nil {
		return nil, err
	}

	// all good, commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}
