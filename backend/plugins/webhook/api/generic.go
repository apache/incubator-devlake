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
	Data        map[string]interface{} `mapstructure:"json"`

	CreatedDate *time.Time `mapstructure:"createdDate"`
}

type Generic struct {
	Title       string
	Description string
	CreatedDate *time.Time
	Data        string
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
	if err := SaveGeneric(connection, request, tx, logger); err != nil {
		logger.Error(err, "create generic")
		return nil, err
	}

	return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}

func SaveGeneric(connection *models.WebhookConnection, request *WebhookGenericReq, db dal.Transaction, logger log.Logger) errors.Error {
	// validation
	if request == nil {
		return errors.BadInput.New("request body is nil")
	}
	if len(request.Data) == 0 {
		return errors.BadInput.New("json payload is empty")
	}
	createdDate := time.Now()
	if request.CreatedDate != nil {
		createdDate = *request.CreatedDate
	}
	if request.CreatedDate == nil {
		request.CreatedDate = &createdDate
	}

	generic := new(Generic)
	generic.Title = request.Title
	generic.Description = request.Description
	generic.CreatedDate = request.CreatedDate
	generic.Data = fmt.Sprintf("%s", request.Data)

	db.AutoMigrate(generic)
	err := db.CreateOrUpdate(generic)
	if err != nil {
		return err
	}

	return nil
}
