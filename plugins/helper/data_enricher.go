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
	"github.com/apache/incubator-devlake/errors"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core"
)

// DataEnrichHandler Accept row from source cursor, return list of entities that need to be stored
type DataEnrichHandler func(row interface{}) ([]interface{}, error)

// DataEnricherArgs includes the arguments about DataEnricher.
// This will be used in Creating a DataEnricher.
//
//	DataEnricherArgs {
//				InputRowType: 		type of inputRow ,
//				Input:        		dal cursor,
//				RawDataSubTaskArgs: args about raw data task
//				Enrich: 			main function including conversion logic
//				BatchSize: 			batch size
type DataEnricherArgs struct {
	RawDataSubTaskArgs
	// Domain layer entity Id prefix, i.e. `jira:JiraIssue:1`, `github:GithubIssue`
	InputRowType reflect.Type
	Input        *sql.Rows
	Enrich       DataEnrichHandler
	BatchSize    int
}

// DataEnricher helps you convert Data from Tool Layer Tables to Domain Layer Tables
// It reads rows from specified Iterator, and feed it into `Enricher` handler
// you can return arbitrary domain layer entities from this handler, ApiEnricher would
// first delete old data by their RawDataOrigin information, and then perform a
// batch save operation for you.
type DataEnricher struct {
	*RawDataSubTask
	args *DataEnricherArgs
}

// NewDataEnricher function helps you create a DataEnricher using DataEnricherArgs.
// You can see the usage in plugins/github/tasks/pr_issue_convertor.go or other convertor file.
func NewDataEnricher(args DataEnricherArgs) (*DataEnricher, error) {
	rawDataSubTask, err := newRawDataSubTask(args.RawDataSubTaskArgs)
	if err != nil {
		return nil, err
	}
	// process args
	if args.BatchSize == 0 {
		args.BatchSize = 500
	}
	return &DataEnricher{
		RawDataSubTask: rawDataSubTask,
		args:           &args,
	}, nil
}

// Execute function implements Subtask interface.
// It loads data from Tool Layer Tables using `Ctx.GetDal()`, convert Data using `converter.args.Enrich` handler
// Then save data to Domain Layer Tables using BatchSaveDivider
func (enricher *DataEnricher) Execute() error {
	// load data from database
	db := enricher.args.Ctx.GetDal()

	divider := NewBatchUpdateDivider(enricher.args.Ctx, enricher.args.BatchSize, enricher.table, enricher.params)

	// set progress
	enricher.args.Ctx.SetProgress(0, -1)

	cursor := enricher.args.Input
	defer cursor.Close()
	ctx := enricher.args.Ctx.GetContext()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		inputRow := reflect.New(enricher.args.InputRowType).Interface()
		err := db.Fetch(cursor, inputRow)
		if err != nil {
			return errors.Default.Wrap(err, "error fetching rows", errors.UserMessage("Internal Enricher execution error"))
		}

		results, err := enricher.args.Enrich(inputRow)
		if err != nil {
			return errors.Default.Wrap(err, "error calling Enricher plugin implementation", errors.UserMessage("Internal Enricher execution error"))
		}

		for _, result := range results {
			// get the batch operator for the specific type
			batch, err := divider.ForType(reflect.TypeOf(result))
			if err != nil {
				return errors.Default.Wrap(err, "error getting batch from result", errors.UserMessage("Internal Enricher execution error"))
			}
			// records get saved into db when slots were max outed
			err = batch.Add(result)
			if err != nil {
				return errors.Default.Wrap(err, "error updating result to batch", errors.UserMessage("Internal Enricher execution error"))
			}
		}
		enricher.args.Ctx.IncProgress(1)
	}

	// save the last batches
	return divider.Close()
}

// Check if DataEnricher implements SubTask interface
var _ core.SubTask = (*DataEnricher)(nil)
