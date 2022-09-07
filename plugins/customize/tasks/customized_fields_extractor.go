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

package tasks

import (
	"context"
	"fmt"
	"strings"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/customize/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/tidwall/gjson"
)

var _ core.SubTaskEntryPoint = ExtractCustomizedFields

var ExtractCustomizedFieldsMeta = core.SubTaskMeta{Name: "extractCustomizedFields",
	EntryPoint:       ExtractCustomizedFields,
	EnabledByDefault: true,
	Description:      "extract customized fields",
}

// ExtractCustomizedFields extracts fields from raw data tables and assigns to domain layer tables
func ExtractCustomizedFields(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TaskData)
	if data == nil || data.Options == nil {
		return nil
	}
	d := taskCtx.GetDal()
	var err error
	for _, rule := range data.Options.TransformationRules {
		err = extractCustomizedFields(taskCtx.GetContext(), d, rule.Table, rule.RawDataTable, rule.RawDataParams, rule.Mapping)
		if err != nil {
			return err
		}
	}
	return nil
}

func extractCustomizedFields(ctx context.Context, d dal.Dal, table, rawTable, rawDataParams string, extractor map[string]string) error {
	pkFields, err := dal.GetPrimarykeyColumns(d, &models.Table{table})
	if err != nil {
		return err
	}
	fields := []string{"_raw_data_id"}
	for _, field := range pkFields {
		fields = append(fields, field.Name())
	}
	clauses := []dal.Clause{
		dal.Select(strings.Join(fields, ", ")),
		dal.From(table),
		dal.Where("_raw_data_table = ? AND _raw_data_params = ?", rawTable, rawDataParams),
	}
	rows, err := d.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		data := make(map[string]interface{})
		updates := make(map[string]string)
		err = d.Fetch(rows, &data)
		if err != nil {
			return err
		}
		var raw helper.RawData
		cls := []dal.Clause{
			dal.Select("*"),
			dal.From(rawTable),
			dal.Where("id = ?", data["_raw_data_id"]),
		}
		err = d.First(&raw, cls...)
		if err != nil {
			return err
		}
		for field, path := range extractor {
			updates[field] = gjson.GetBytes(raw.Data, path).String()
		}
		if len(updates) > 0 {
			delete(data, "_raw_data_id")
			query, params := mkUpdate(table, updates, data)
			err = d.Exec(query, params...)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func mkUpdate(table string, updates map[string]string, pk map[string]interface{}) (string, []interface{}) {
	var params []interface{}
	stat := fmt.Sprintf("UPDATE %s SET ", table)
	var uu []string
	for field, value := range updates {
		uu = append(uu, fmt.Sprintf("%s = ?", field))
		params = append(params, value)
	}
	var ww []string
	for field, value := range pk {
		ww = append(ww, fmt.Sprintf("%s = ?", field))
		params = append(params, value)
	}
	return stat + strings.Join(uu, ", ") + " WHERE " + strings.Join(ww, " AND "), params
}
