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

package helper

import (
	"database/sql"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core"
)

// Accept row from source cursor, return list of entities that need to be stored
type DataConvertHandler func(row interface{}) ([]interface{}, error)

type DataConverterArgs struct {
	RawDataSubTaskArgs
	// Domain layer entity Id prefix, i.e. `jira:JiraIssue:1`, `github:GithubIssue`
	InputRowType reflect.Type
	// Cursor to a set of Tool Layer Records
	Input     *sql.Rows
	Convert   DataConvertHandler
	BatchSize int
}

// DataConverter helps you convert Data from Tool Layer Tables to Domain Layer Tables
// It reads rows from specified Iterator, and feed it into `Convter` handler
// you can return arbitrary domain layer entities from this handler, ApiConverter would
// first delete old data by their RawDataOrigin information, and then perform a
// batch save operation for you.
type DataConverter struct {
	*RawDataSubTask
	args *DataConverterArgs
}

func NewDataConverter(args DataConverterArgs) (*DataConverter, error) {
	rawDataSubTask, err := newRawDataSubTask(args.RawDataSubTaskArgs)
	if err != nil {
		return nil, err
	}
	// process args
	if args.BatchSize == 0 {
		args.BatchSize = 500
	}
	return &DataConverter{
		RawDataSubTask: rawDataSubTask,
		args:           &args,
	}, nil
}

func (converter *DataConverter) Execute() error {
	// load data from database
	db := converter.args.Ctx.GetDb()

	// batch save divider
	RAW_DATA_ORIGIN := "RawDataOrigin"
	divider := NewBatchSaveDivider(converter.args.Ctx, converter.args.BatchSize, converter.table, converter.params)

	// prgress
	converter.args.Ctx.SetProgress(0, -1)

	cursor := converter.args.Input
	defer cursor.Close()
	ctx := converter.args.Ctx.GetContext()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		inputRow := reflect.New(converter.args.InputRowType).Interface()
		err := db.ScanRows(cursor, inputRow)
		if err != nil {
			return err
		}

		results, err := converter.args.Convert(inputRow)
		if err != nil {
			return err
		}

		for _, result := range results {
			// get the batch operator for the specific type
			batch, err := divider.ForType(reflect.TypeOf(result))
			if err != nil {
				return err
			}
			// set raw data origin field
			origin := reflect.ValueOf(result).Elem().FieldByName(RAW_DATA_ORIGIN)
			if origin.IsValid() {
				origin.Set(reflect.ValueOf(inputRow).Elem().FieldByName(RAW_DATA_ORIGIN))
			}
			// records get saved into db when slots were max outed
			err = batch.Add(result)
			if err != nil {
				return err
			}
		}
		converter.args.Ctx.IncProgress(1)
	}

	// save the last batches
	return divider.Close()
}

var _ core.SubTask = (*DataConverter)(nil)
