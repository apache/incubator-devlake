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
	Url         string `mapstructure:"url"`
	IssueKey    string `mapstructure:"issueKey" validate:"required"`
	Title       string `mapstructure:"title" validate:"required"`
	Description string `mapstructure:"description"`
	Name        string `mapstructure:"name"`
	Json        string `mapstructure:"json"`

	CreatedDate  *time.Time `mapstructure:"createdDate"`
	StartedDate  *time.Time `mapstructure:"startedDate" validate:"required"`
	FinishedDate *time.Time `mapstructure:"finishedDate" validate:"required"`
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
		logger.Error(err, "create deployments")
		return nil, err
	}

	return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
}

func Create(connection *models.WebhookConnection, request *WebhookGenericReq, tx dal.Transaction, logger log.Logger) errors.Error {
	// validation
	if request == nil {
		return errors.BadInput.New("request body is nil")
	}
	if len(request.Json) == 0 {
		return errors.BadInput.New("json payload is empty")
	}
	if request.CreatedDate == nil {
		request.CreatedDate = request.StartedDate
	}
	if request.FinishedDate == nil {
		now := time.Now()
		request.FinishedDate = &now
	}
	createdDate := time.Now()
	if request.CreatedDate != nil {
		createdDate = *request.CreatedDate
	} else if request.StartedDate != nil {
		createdDate = *request.StartedDate
	}
	if request.CreatedDate == nil {
		request.CreatedDate = &createdDate
	}

	if err := tx.CreateOrUpdate(request.Json); err != nil {
		logger.Error(err, "failed to generic json data to disk")
		return err
	}

	return nil
}
