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

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/log"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/dbhelper"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/webhook/models"
	"github.com/go-playground/validator/v10"
)

type WebhookGenericReq struct {
	Title       string                 `mapstructure:"title" validate:"required"`
	Description string                 `mapstructure:"description"`
	Json        map[string]interface{} `mapstructure:"json"`

	CreatedDate *time.Time `mapstructure:"createdDate"`
}

type Generic struct {
	Id          int
	Title       string
	Description string
	CreatedDate *time.Time
	data        string
}

func PostGeneric(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}

	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	// get request
	request := &WebhookGenericReq{}

	err = api.DecodeMapStruct(input.Body, request, true)
	if err != nil {
		return &plugin.ApiResourceOutput{Body: err.Error(), Status: http.StatusBadRequest}, nil
	}

	// validate
	vld = validator.New()
	err = errors.Convert(vld.Struct(request))
	if err != nil {
		return nil, errors.BadInput.Wrap(vld.Struct(request), `input json error`)
	}

	txHelper := dbhelper.NewTxHelper(basicRes, &err)
	defer txHelper.End()
	tx := txHelper.Begin()
	if err := Create(connection, request, tx, logger); err != nil {
		logger.Error(err, "create generic")
		return nil, err
	}

	return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}

func Create(connection *models.WebhookConnection, request *WebhookGenericReq, db dal.Transaction, logger log.Logger) errors.Error {
	// validation
	if request == nil {
		return errors.BadInput.New("request body is nil")
	}
	if len(request.Json) == 0 {
		return errors.BadInput.New("json payload is empty")
	}
	createdDate := time.Now()
	if request.CreatedDate != nil {
		createdDate = *request.CreatedDate
	}
	if request.CreatedDate == nil {
		request.CreatedDate = &createdDate
	}

	fmt.Println(request.Json)

	jsonString, err := json.Marshal(request.Json)
	fmt.Println(jsonString)
	if err != nil {
		logger.Error(err, "Error marshaling JSON:")
	}

	generic := new(Generic)
	generic.Title = request.Title
	generic.Description = request.Description
	generic.CreatedDate = request.CreatedDate
	generic.data = string(jsonString)

	if !db.HasTable(generic) {
		db.AutoMigrate(generic)
	}
	db.AutoMigrate(generic)
	err = db.Create(generic)
	if err != nil {
		return nil
	}
	db.All(generic)

	// if err := tx.CreateOrUpdate(generic); err != nil {
	// 	logger.Error(err, "failed to generic json data to disk")
	// 	return err
	// }

	return nil
}
