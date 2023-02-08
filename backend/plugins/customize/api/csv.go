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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"strings"
)

const maxMemory = 32 << 20 // 32 MB

// ImportCSVFile accepts a CSV file, parses and saves it to the database
// @Summary      Upload CSV file
// @Description  Upload CSV file
// @Tags 		 plugins/customize
// @Accept       multipart/form-data
// @Param        file formData file true "select file to upload"
// @Param        table formData string true "the table name, only issues and issue_commits are supported"
// @Param        rawDataParams formData string true "the value of _raw_data_params"
// @Produce      json
// @Success      200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/customize/csvfile [post]
func (h *Handlers) ImportCSVFile(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	if input.Request == nil {
		return nil, errors.Default.New("request is nil")
	}
	if input.Request.MultipartForm == nil {
		if err := input.Request.ParseMultipartForm(maxMemory); err != nil {
			return nil, errors.Convert(err)
		}
	}
	f, fh, err := input.Request.FormFile("file")
	if err != nil {
		return nil, errors.Convert(err)
	}
	// nolint
	f.Close()
	file, err := fh.Open()
	if err != nil {
		return nil, errors.Convert(err)
	}
	// nolint
	defer file.Close()
	table := strings.TrimSpace(input.Request.FormValue("table"))
	if table == "" {
		return nil, errors.Default.New("empty table")
	}
	rawDataParams := strings.TrimSpace(input.Request.FormValue("rawDataParams"))
	if rawDataParams == "" {
		return nil, errors.Default.New("empty rawDataParams")
	}
	return nil, h.svc.ImportCSV(table, rawDataParams, file)
}
