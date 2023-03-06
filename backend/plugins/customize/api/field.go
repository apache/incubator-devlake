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
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/customize/models"
	"github.com/apache/incubator-devlake/plugins/customize/service"
)

type fieldResponse struct {
	Field
	IsCustomizedField bool `json:"isCustomizedField" example:"true"`
}

type Field struct {
	ColumnName  string `json:"columnName" example:"x_column_varchar"`
	DisplayName string `json:"displayName" example:"department"`
	DataType    string `json:"dataType" example:"varchar(255)"`
	Description string `json:"description" example:"more details about the column"`
}

func (f *Field) toDBModel(table string) (*models.CustomizedField, errors.Error) {
	if !strings.HasPrefix(f.ColumnName, "x_") {
		return nil, errors.BadInput.New("the columnName should start with x_")
	}
	if f.DisplayName == "" {
		return nil, errors.BadInput.New("the displayName is empty")
	}
	t, ok := dal.ToColumnType(f.DataType)
	if !ok {
		return nil, errors.BadInput.New(fmt.Sprintf("the columnType:%s is unsupported", f.DataType))
	}
	return &models.CustomizedField{
		TbName:      table,
		ColumnName:  f.ColumnName,
		DisplayName: f.DisplayName,
		DataType:    t,
		Description: f.Description,
	}, nil
}

func fromCustomizedField(cf models.CustomizedField) fieldResponse {
	return fieldResponse{
		Field: Field{
			ColumnName:  cf.ColumnName,
			DisplayName: cf.DisplayName,
			DataType:    cf.DataType.String(),
			Description: cf.Description,
		},
		IsCustomizedField: strings.HasPrefix(cf.ColumnName, "x_"),
	}
}

type Handlers struct {
	svc *service.Service
}

func NewHandlers(dal dal.Dal) *Handlers {
	return &Handlers{svc: service.NewService(dal)}
}

// ListFields return all customized fields
// @Summary return all customized fields
// @Description return all customized fieldsh
// @Tags plugins/customize
// @Param table path string true "the table name"
// @Success 200  {object} []fieldResponse "Success"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/customize/{table}/fields [GET]
func (h *Handlers) ListFields(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	customizedFields, err := h.svc.GetFields(input.Params["table"])
	if err != nil {
		return &plugin.ApiResourceOutput{Status: http.StatusBadRequest}, errors.Default.Wrap(err, "getFields error")
	}
	fields := make([]fieldResponse, 0, len(customizedFields))
	for _, cf := range customizedFields {
		fields = append(fields, fromCustomizedField(cf))
	}
	return &plugin.ApiResourceOutput{Body: fields, Status: http.StatusOK}, nil
}

// CreateFields create a customized field
// @Summary create a customized field
// @Description create a customized field
// @Tags plugins/customize
// @Param table path string true "the table name"
// @Param request body Field true "request body"
// @Success 200  {object} fieldResponse "Success"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/customize/{table}/fields [POST]
func (h *Handlers) CreateFields(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	table := input.Params["table"]
	fld := &Field{}
	err := helper.Decode(input.Body, fld, nil)
	if err != nil {
		return &plugin.ApiResourceOutput{Status: http.StatusBadRequest}, err
	}
	customizedField, err := fld.toDBModel(table)
	if err != nil {
		return &plugin.ApiResourceOutput{Status: http.StatusBadRequest}, err
	}
	err = h.svc.CreateField(customizedField)
	if err != nil {
		return nil, errors.Default.Wrap(err, "CreateField error")
	}
	return &plugin.ApiResourceOutput{Body: fieldResponse{*fld, true}, Status: http.StatusOK}, nil
}

// DeleteField delete a customized fields
// @Summary return all customized fields
// @Description return all customized fields
// @Tags plugins/customize
// @Param table path string true "the table name"
// @Param field path string true "the column to be deleted"
// @Success 200  {object} shared.ApiBody "Success"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/customize/{table}/fields/{field} [DELETE]
func (h *Handlers) DeleteField(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	table := input.Params["table"]
	fld := input.Params["field"]
	err := h.svc.DeleteField(table, fld)
	if err != nil {
		return &plugin.ApiResourceOutput{Status: http.StatusBadRequest}, errors.Default.Wrap(err, "deleteField error")
	}
	return &plugin.ApiResourceOutput{Status: http.StatusOK}, nil
}
