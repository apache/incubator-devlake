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
	"regexp"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	plugin "github.com/apache/incubator-devlake/core/plugin"
)

// DataEnrichHandler Accepts row from the Input and produces arbitrary records.
// you are free to modify given `row` in place and include it in returned result for it to be saved.
type DataEnrichHandler[InputRowType any] func(row *InputRowType) ([]interface{}, errors.Error)

// DataEnricherArgs includes the arguments needed for data enrichment
type DataEnricherArgs[InputRowType any] struct {
	Ctx       plugin.SubTaskContext
	Name      string // Enricher name, which will be put into _raw_data_remark
	Input     dal.Rows
	Enrich    DataEnrichHandler[InputRowType]
	BatchSize int
}

// DataEnricher helps you enrich Data with Cancellation and BatchSave supports
type DataEnricher[InputRowType any] struct {
	args *DataEnricherArgs[InputRowType]
}

var dataEnricherNamePattern = regexp.MustCompile(`^\w+$`)

// NewDataEnricher creates a new DataEnricher
func NewDataEnricher[InputRowType any](args DataEnricherArgs[InputRowType]) (*DataEnricher[InputRowType], errors.Error) {
	// process args
	if args.Name == "" || !dataEnricherNamePattern.MatchString(args.Name) {
		return nil, errors.Default.New("DataEnricher: Name is require and should contain only word characters(a-zA-Z0-9_)")
	}
	if args.BatchSize == 0 {
		args.BatchSize = 500
	}
	return &DataEnricher[InputRowType]{
		args: &args,
	}, nil
}

func (enricher *DataEnricher[InputRowType]) Execute() errors.Error {
	// load data from database
	db := enricher.args.Ctx.GetDal()

	// batch save divider
	divider := NewBatchSaveDivider(enricher.args.Ctx, enricher.args.BatchSize, "", "")

	// set progress
	enricher.args.Ctx.SetProgress(0, -1)

	cursor := enricher.args.Input
	defer cursor.Close()
	ctx := enricher.args.Ctx.GetContext()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.Convert(ctx.Err())
		default:
		}
		inputRow := new(InputRowType)
		err := db.Fetch(cursor, inputRow)
		if err != nil {
			return errors.Default.Wrap(err, "error fetching rows")
		}

		results, err := enricher.args.Enrich(inputRow)
		if err != nil {
			return errors.Default.Wrap(err, "error calling plugin implementation")
		}

		for _, result := range results {
			// get the batch operator for the specific type
			batch, err := divider.ForType(reflect.TypeOf(result))
			if err != nil {
				return errors.Default.Wrap(err, "error getting batch from result")
			}
			// append enricher to data origin remark
			if getRawDataOrigin, ok := result.(common.GetRawDataOrigin); ok {
				origin := getRawDataOrigin.GetRawDataOrigin()
				enricherComponent := enricher.args.Name + "," // name is word characters only
				if !strings.Contains(origin.RawDataRemark, enricherComponent) {
					origin.RawDataRemark += enricherComponent
				}
			}
			// records get saved into db when slots were max outed
			err = batch.Add(result)
			if err != nil {
				return errors.Default.Wrap(err, "error adding result to batch")
			}
		}
		enricher.args.Ctx.IncProgress(1)
	}

	// save the last batches
	return divider.Close()
}

// Check if DataEnricher implements SubTask interface
var _ plugin.SubTask = (*DataEnricher[any])(nil)
