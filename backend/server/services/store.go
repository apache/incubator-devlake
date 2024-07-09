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

func GetStore(storeKey string) (*models.Store, errors.Error) {
	clauses := []dal.Clause{
		dal.From(&models.Store{}),
		dal.Where("store_key = ?", storeKey),
	}

	kvstore := &models.Store{}
	err := db.All(&kvstore, clauses...)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("error finding %s on _devlake_store table", storeKey))
	}

	return kvstore, nil
}

// PutOnboard accepts a project instance and insert it to database
func PutStore(storeKey string, storeValue *models.Store) (*models.Store, errors.Error) {
	// verify input
	if err := VerifyStruct(storeValue); err != nil {
		return nil, err
	}

	err := db.CreateOrUpdate(storeValue, dal.Where("store_key = ?", storeKey))
	if err != nil {
		return nil, err
	}

	return storeValue, nil
}
