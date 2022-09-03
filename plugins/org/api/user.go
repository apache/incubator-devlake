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
	"github.com/apache/incubator-devlake/errors"
	"net/http"

	"github.com/apache/incubator-devlake/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/gocarina/gocsv"
)

// GetUser returns all user in csv format
// @Summary      Get users.csv file
// @Description  get users.csv file
// @Tags 		 plugins/org
// @Produce      text/csv
// @Param        fake_data    query     bool  false  "return fake data or not"
// @Success      200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/users.csv [get]
func (h *Handlers) GetUser(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var users []user
	var u *user
	var err errors.Error
	if input.Query.Get("fake_data") == "true" {
		users = u.fakeData()
	} else {
		users, err = h.store.findAllUsers()
		if err != nil {
			return nil, err
		}
	}
	blob, err := errors.Convert01(gocsv.MarshalBytes(users))
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{
		Body:   nil,
		Status: http.StatusOK,
		File: &core.OutputFile{
			ContentType: "text/csv",
			Data:        blob,
		},
	}, nil
}

// CreateUser accepts a CSV file containing user information mapping and saves it to the database
// @Summary      Upload users.csv file
// @Description  upload users.csv file
// @Tags 		 plugins/org
// @Accept       multipart/form-data
// @Param        file formData file true "select file to upload"
// @Produce      json
// @Success      200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/users.csv [put]
func (h *Handlers) CreateUser(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var uu []user
	err := h.unmarshal(input.Request, &uu)
	if err != nil {
		return nil, err
	}
	var u *user
	var items []interface{}
	users, teamUsers := u.toDomainLayer(uu)
	for _, user := range users {
		items = append(items, user)
	}
	for _, teamUser := range teamUsers {
		items = append(items, teamUser)
	}
	err = h.store.deleteAll(&crossdomain.User{})
	if err != nil {
		return nil, err
	}
	err = h.store.deleteAll(&crossdomain.TeamUser{})
	if err != nil {
		return nil, err
	}
	err = h.store.save(items)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Status: http.StatusOK}, nil
}
