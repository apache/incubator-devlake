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
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

// ImportQaApis accepts a CSV file, parses and saves it to the database
// @Summary      Upload qa_apis.csv file
// @Description  Upload qa_apis.csv file.
// @Tags 		 plugins/customize
// @Accept       multipart/form-data
// @Param        qaProjectId formData string true "the ID of the QA project"
// @Param        file formData file true "select file to upload"
// @Param        incremental formData bool false "incremental import"
// @Produce      json
// @Success      200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/customize/csvfiles/qa_apis.csv [post]
func (h *Handlers) ImportQaApis(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	file, err := h.extractFile(input)
	if err != nil {
		return nil, err
	}
	// nolint
	defer file.Close()

	incremental := false
	if input.Request.FormValue("incremental") == "true" {
		incremental = true
	}

	qaProjectId := strings.TrimSpace(input.Request.FormValue("qaProjectId"))
	if qaProjectId == "" {
		return nil, errors.BadInput.New("empty qaProjectId")
	}

	return nil, h.svc.ImportQaApis(qaProjectId, file, incremental)

}

// ImportQaTestCases accepts a CSV file, parses and saves it to the database
// @Summary      Upload qa_test_cases.csv file
// @Description  Upload qa_test_cases.csv file.
// @Tags 		 plugins/customize
// @Accept       multipart/form-data
// @Param        qaProjectId formData string true "the ID of the QA project"
// @Param        qaProjectName formData string true "the name of the QA project"
// @Param        file formData file true "select file to upload"
// @Param        incremental formData bool false "incremental update"
// @Produce      json
// @Success      200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/customize/csvfiles/qa_test_cases.csv [post]
func (h *Handlers) ImportQaTestCases(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	file, err := h.extractFile(input)
	if err != nil {
		return nil, err
	}
	// nolint
	defer file.Close()

	incremental := false
	if input.Request.FormValue("incremental") == "true" {
		incremental = true
	}

	qaProjectId := strings.TrimSpace(input.Request.FormValue("qaProjectId"))
	if qaProjectId == "" {
		return nil, errors.BadInput.New("empty qaProjectId")
	}
	qaProjectName := strings.TrimSpace(input.Request.FormValue("qaProjectName"))
	if qaProjectName == "" {
		return nil, errors.BadInput.New("empty qaProjectName")
	}
	return nil, h.svc.ImportQaTestCases(qaProjectId, qaProjectName, file, incremental) // records contains the CSV data
}

// ImportQaTestCaseExecutions accepts a CSV file, parses and saves it to the database
// @Summary      Upload qa_test_case_executions.csv file
// @Description  Upload qa_test_case_executions.csv file.
// @Tags 		 plugins/customize
// @Accept       multipart/form-data
// @Param        qaProjectId formData string true "the ID of the QA project"
// @Param        file formData file true "select file to upload"
// @Param        incremental formData bool false "incremental update"
// @Produce      json
// @Success      200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/customize/csvfiles/qa_test_case_executions.csv [post]
func (h *Handlers) ImportQaTestCaseExecutions(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	file, err := h.extractFile(input)
	if err != nil {
		return nil, err
	}
	// nolint
	defer file.Close()

	incremental := false
	if input.Request.FormValue("incremental") == "true" {
		incremental = true
	}

	qaProjectId := strings.TrimSpace(input.Request.FormValue("qaProjectId"))
	if qaProjectId == "" {
		return nil, errors.BadInput.New("empty qaProjectId")
	}

	return nil, h.svc.ImportQaTestCaseExecutions(qaProjectId, file, incremental)
}
