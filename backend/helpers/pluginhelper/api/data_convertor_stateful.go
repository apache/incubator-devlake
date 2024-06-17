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
	plugin "github.com/apache/incubator-devlake/core/plugin"
)

type StatefulDataConverterArgs[InputType any] struct {
	*SubtaskCommonArgs
	Input     func(*SubtaskStateManager) (dal.Rows, errors.Error)
	Convert   func(row *InputType) ([]any, errors.Error)
	BatchSize int
}

type StatefulDataConverter[InputType any] struct {
	*StatefulDataConverterArgs[InputType]
	*SubtaskStateManager
}

func NewStatefulDataConverter[
	OptType any,
	InputType any,
](
	args *StatefulDataConverterArgs[InputType],
) (*StatefulDataConverter[InputType], errors.Error) {
	// process args
	if args.BatchSize == 0 {
		args.BatchSize = 500
	}
	stateManager, err := NewSubtaskStateManager(args.SubtaskCommonArgs)
	if err != nil {
		return nil, err
	}
	return &StatefulDataConverter[InputType]{
		StatefulDataConverterArgs: args,
		SubtaskStateManager:       stateManager,
	}, nil
}

func (converter *StatefulDataConverter[InputType]) Execute() errors.Error {
	// load data from database
	db := converter.GetDal()

	table := converter.GetRawDataTable()
	params := converter.GetRawDataParams()

	// batch save divider
	RAW_DATA_ORIGIN := "RawDataOrigin"
	divider := NewBatchSaveDivider(converter, converter.BatchSize, table, params)
	divider.SetIncrementalMode(converter.IsIncremental())

	// set progress
	converter.SetProgress(0, -1)

	cursor, err := converter.Input(converter.SubtaskStateManager)
	if err != nil {
		return err
	}
	defer cursor.Close()
	ctx := converter.GetContext()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.Convert(ctx.Err())
		default:
		}
		inputRow := new(InputType)
		err := db.Fetch(cursor, inputRow)
		if err != nil {
			return errors.Default.Wrap(err, "error fetching rows")
		}

		results, err := converter.Convert(inputRow)
		if err != nil {
			return errors.Default.Wrap(err, "error calling Converter plugin implementation")
		}

		for _, result := range results {
			// get the batch operator for the specific type
			batch, err := divider.ForType(reflect.TypeOf(result))
			if err != nil {
				return errors.Default.Wrap(err, "error getting batch from result")
			}
			// set raw data origin field
			origin := reflect.ValueOf(result).Elem().FieldByName(RAW_DATA_ORIGIN)
			if origin.IsValid() {
				origin.Set(reflect.ValueOf(inputRow).Elem().FieldByName(RAW_DATA_ORIGIN))
			}
			// records get saved into db when slots were max outed
			err = batch.Add(result)
			if err != nil {
				return errors.Default.Wrap(err, "error adding result to batch")
			}
		}
		converter.IncProgress(1)
	}

	// save the last batches
	err = divider.Close()
	if err != nil {
		return err
	}
	// save the incremantal state
	return converter.SubtaskStateManager.Close()
}

// Check if DataConverter implements SubTask interface
var _ plugin.SubTask = (*StatefulDataConverter[any])(nil)
