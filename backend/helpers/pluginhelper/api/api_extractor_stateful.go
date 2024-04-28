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
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	plugin "github.com/apache/incubator-devlake/core/plugin"
)

// StatefulApiExtractorArgs is a struct that contains the arguments for a stateful api extractor
type StatefulApiExtractorArgs struct {
	*SubtaskCommonArgs
	Extract func(row *RawData) ([]any, errors.Error)
}

type StatefulApiExtractor struct {
	*StatefulApiExtractorArgs
	*SubtaskStateManager
}

// NewStatefulApiExtractor creates a new StatefulApiExtractor
func NewStatefulApiExtractor(args *StatefulApiExtractorArgs) (*StatefulApiExtractor, errors.Error) {
	stateManager, err := NewSubtaskStateManager(args.SubtaskCommonArgs)
	if err != nil {
		return nil, err
	}
	return &StatefulApiExtractor{
		StatefulApiExtractorArgs: args,
		SubtaskStateManager:      stateManager,
	}, nil
}

// Execute sub-task
func (extractor *StatefulApiExtractor) Execute() errors.Error {
	// load data from database
	db := extractor.GetDal()
	logger := extractor.GetLogger()
	table := extractor.GetRawDataTable()
	params := extractor.GetRawDataParams()
	if !db.HasTable(table) {
		return nil
	}
	clauses := []dal.Clause{
		dal.From(table),
		dal.Where("params = ?", params),
		dal.Orderby("id ASC"),
	}

	if extractor.IsIncremental() {
		since := extractor.GetSince()
		if since != nil {
			clauses = append(clauses, dal.Where("created_at >= ? ", since))
		}
	}
	clauses = append(clauses, dal.Where("created_at < ? ", extractor.GetUntil()))

	count, err := db.Count(clauses...)
	if err != nil {
		return errors.Default.Wrap(err, "error getting count of clauses")
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return errors.Default.Wrap(err, "error running DB query")
	}
	logger.Info("get data from %s where params=%s and got %d", table, params, count)
	defer cursor.Close()
	// batch save divider
	divider := NewBatchSaveDivider(extractor.SubTaskContext, extractor.GetBatchSize(), table, params)
	divider.SetIncrementalMode(extractor.IsIncremental())

	// progress
	extractor.SetProgress(0, -1)
	ctx := extractor.GetContext()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.Convert(ctx.Err())
		default:
		}
		row := &RawData{}
		err = db.Fetch(cursor, row)
		if err != nil {
			return errors.Default.Wrap(err, "error fetching row")
		}

		results, err := extractor.Extract(row)
		if err != nil {
			return errors.Default.Wrap(err, "error calling plugin Extract implementation")
		}
		for _, result := range results {
			// get the batch operator for the specific type
			batch, err := divider.ForType(reflect.TypeOf(result))
			if err != nil {
				return errors.Default.Wrap(err, "error getting batch from result")
			}
			// set raw data origin field
			setRawDataOrigin(result, common.RawDataOrigin{
				RawDataTable:  table,
				RawDataParams: params,
				RawDataId:     row.ID,
			})
			// records get saved into db when slots were max outed
			err = batch.Add(result)
			if err != nil {
				return errors.Default.Wrap(err, "error adding result to batch")
			}
		}
		extractor.IncProgress(1)
	}

	// save the last batches
	err = divider.Close()
	if err != nil {
		return err
	}
	// save the incremantal state
	return extractor.SubtaskStateManager.Close()
}

var _ plugin.SubTask = (*StatefulApiExtractor)(nil)
