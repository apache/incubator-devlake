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
	"github.com/apache/incubator-devlake/plugins/core"
	"net/http"

	"github.com/gocarina/gocsv"
)

// GetProjectMapping returns all project mapping in csv format
// @Summary      Get project_mapping.csv file
// @Description  get project_mapping.csv file
// @Tags 		 plugins/org
// @Produce      text/csv
// @Param        fake_data    query     bool  false  "return fake data or not"
// @Success      200 {object} core.ApiResourceOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/project_mapping.csv [get]
func (h *Handlers) GetProjectMapping(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	input.Query.Get("fake_data")
	var mapping []projectMapping
	var err error
	if input.Query.Get("fake_data") == "true" {
		mapping = fakeProjectMapping
	} else {
		mapping, err = h.store.findAllProjectMapping()
		if err != nil {
			return nil, errors.Convert(err)
		}
	}
	blob, err := gocsv.MarshalBytes(mapping)
	if err != nil {
		return nil, errors.Convert(err)
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

// CreateProjectMapping accepts a CSV file containing project mapping and saves it to the database
// @Summary      Upload project_mapping.csv file
// @Description  upload project_mapping.csv file
// @Tags 		 plugins/org
// @Accept       multipart/form-data
// @Param        file formData file true "select file to upload"
// @Produce      json
// @Success      200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/project_mapping.csv [put]
func (h *Handlers) CreateProjectMapping(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var mapping []projectMapping
	err := h.unmarshal(input.Request, &mapping)
	if err != nil {
		return nil, err
	}
	var pm *projectMapping
	var items []interface{}
	for _, tm := range pm.toDomainLayer(mapping) {
		items = append(items, tm)
	}
	err = h.store.save(items)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Status: http.StatusOK}, nil
}
