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
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/customize/models"
	"gorm.io/gorm/clause"
)

type field struct {
	ColumnName string `json:"columnName"`
	ColumnType string `json:"columnType"`
}

func getFields(d dal.Dal, tbl string) ([]field, error) {
	columns, err := d.GetColumns(&models.Table{tbl}, func(columnMeta dal.ColumnMeta) bool {
		return strings.HasPrefix(columnMeta.Name(), "x_")
	})
	if err != nil {
		return nil, err
	}
	var result []field
	for _, col := range columns {
		result = append(result, field{
			ColumnName: col.Name(),
			ColumnType: "VARCHAR(255)",
		})
	}
	return result, nil
}
func checkField(d dal.Dal, table, field string) (bool, error) {
	if !strings.HasPrefix(field, "x_") {
		return false, errors.New("column name should start with `x_`")
	}
	fields, err := getFields(d, table)
	if err != nil {
		return false, err
	}
	for _, fld := range fields {
		if fld.ColumnName == field {
			return true, nil
		}
	}
	return false, nil
}

func CreateField(d dal.Dal, table, field string) error {
	exists, err := checkField(d, table, field)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return d.Exec("ALTER TABLE ? ADD ? VARCHAR(255)", clause.Table{Name: table}, clause.Column{Name: field})
}

func deleteField(d dal.Dal, table, field string) error {
	exists, err := checkField(d, table, field)
	if err != nil {
		return err
	}
	if !exists {
		return nil
	}
	return d.Exec("ALTER TABLE ? DROP COLUMN ?", clause.Table{Name: table}, clause.Column{Name: field})
}

type input struct {
	Name string `json:"name" example:"x_new_column"`
}
type Handlers struct {
	dal dal.Dal
}

func NewHandlers(dal dal.Dal) *Handlers {
	return &Handlers{dal: dal}
}

// ListFields return all customized fields
// @Summary return all customized fields
// @Description return all customized fields
// @Tags plugins/customize
// @Success 200  {object} shared.ApiBody "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/customize/{table}/fields [GET]
func (h *Handlers) ListFields(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	fields, err := getFields(h.dal, input.Params["table"])
	if err != nil {
		return &core.ApiResourceOutput{Status: http.StatusBadRequest}, err
	}
	return &core.ApiResourceOutput{Body: fields, Status: http.StatusOK}, nil
}

// CreateFields create a customized field
// @Summary create a customized field
// @Description create a customized field
// @Tags plugins/customize
// @Param request body input true "request body"
// @Success 200  {object} shared.ApiBody "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/customize/{table}/fields [POST]
func (h *Handlers) CreateFields(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	table := input.Params["table"]
	fld, ok := input.Body["name"].(string)
	if !ok {
		return &core.ApiResourceOutput{Status: http.StatusBadRequest}, fmt.Errorf("the name is not string")
	}
	err := CreateField(h.dal, table, fld)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: field{fld, "varchar(255)"}, Status: http.StatusOK}, nil
}

// DeleteField delete a customized fields
// @Summary return all customized fields
// @Description return all customized fields
// @Tags plugins/customize
// @Success 200  {object} shared.ApiBody "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/customize/{table}/fields [DELETE]
func (h *Handlers) DeleteField(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	table := input.Params["table"]
	fld := input.Params["field"]
	err := deleteField(h.dal, table, fld)
	if err != nil {
		return &core.ApiResourceOutput{Status: http.StatusBadRequest}, err
	}
	return &core.ApiResourceOutput{Status: http.StatusOK}, nil
}
