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
	"encoding/json"
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	plugin "github.com/apache/incubator-devlake/core/plugin"
)

// StatefulApiExtractorArgs is a struct that contains the arguments for a stateful api extractor
type StatefulApiExtractorArgs[InputType any] struct {
	*SubtaskCommonArgs
	BeforeExtract func(issue *InputType, stateManager *SubtaskStateManager) errors.Error
	Extract       func(body *InputType, row *RawData) ([]any, errors.Error)
}

// StatefulApiExtractor is a struct that manages the stateful API extraction process.
// It facilitates extracting data from a single _raw_data table and saving it into multiple Tool Layer tables.
// By default, the extractor operates in Incremental Mode, processing only new records added to the raw table since the previous run.
// This approach reduces the amount of data to process, significantly decreasing the execution time.
// The extractor automatically detects if the configuration has changed since the last run. If a change is detected,
// it will automatically switch to Full-Sync mode.
//
// Example:
//
//	extractor, err := api.NewStatefulApiExtractor(&api.StatefulApiExtractorArgs[apiv2models.Issue]{
//	  SubtaskCommonArgs: &api.SubtaskCommonArgs{
//	    SubTaskContext: subtaskCtx,
//	    Table:          RAW_ISSUE_TABLE,
//	    Params: JiraApiParams{
//	      ConnectionId: data.Options.ConnectionId,
//	      BoardId:      data.Options.BoardId,
//	    },
//	    SubtaskConfig: config,  // The helper stores this configuration in the state and compares it with the previous one
//	                            // to determine the operating mode (Incremental/FullSync).
//	                            // Ensure that the configuration is serializable and contains only public fields.
//	                            // It is also recommended that the configuration includes only the necessary fields used by the extractor.
//	..},
//	  BeforeExtract: func(body *IssuesResponse, stateManager *api.SubtaskStateManager) errors.Error {
//	    if stateManager.IsIncremental() {
//	      // It is important to delete all existing child-records under DiffSync Mode
//	      err := db.Delete(
//	        &models.JiraIssueLabel{},
//	        dal.Where("connection_id = ? AND issue_id = ?", data.Options.ConnectionId, body.Id),
//	      )
//	    }
//	    return nil
//	  },
//	  Extract: func(apiIssue *apiv2models.Issue, row *api.RawData) ([]interface{}, errors.Error) {
//	  },
//	})
//
//	if err != nil {
//	  return err
//	}
//
// return extractor.Execute()
type StatefulApiExtractor[InputType any] struct {
	*StatefulApiExtractorArgs[InputType]
	*SubtaskStateManager
}

// NewStatefulApiExtractor creates a new StatefulApiExtractor
func NewStatefulApiExtractor[InputType any](args *StatefulApiExtractorArgs[InputType]) (*StatefulApiExtractor[InputType], errors.Error) {
	stateManager, err := NewSubtaskStateManager(args.SubtaskCommonArgs)
	if err != nil {
		return nil, err
	}
	return &StatefulApiExtractor[InputType]{
		StatefulApiExtractorArgs: args,
		SubtaskStateManager:      stateManager,
	}, nil
}

// Execute sub-task
func (extractor *StatefulApiExtractor[InputType]) Execute() errors.Error {
	// load data from database
	db := extractor.GetDal()
	logger := extractor.GetLogger()
	table := extractor.GetRawDataTable()
	params := extractor.GetRawDataParams()
	if !db.HasTable(table) {
		return nil
	}

	clauses := []dal.Clause{
		dal.Select("id"),
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

	// first get total count for progress tracking
	count, err := db.Count(clauses...)
	if err != nil {
		return errors.Default.Wrap(err, "error getting count of records")
	}
	logger.Info("get data from %s where params=%s and got %d with clauses %+v", table, params, count, clauses)

	// get all IDs
	var ids []uint64
	err = db.Pluck("id", &ids, clauses...)
	if err != nil {
		return errors.Default.Wrap(err, "error getting IDs")
	}

	// batch save divider
	divider := NewBatchSaveDivider(extractor.SubTaskContext, extractor.GetBatchSize(), table, params)
	divider.SetIncrementalMode(extractor.IsIncremental())

	// progress
	extractor.SetProgress(0, -1)
	ctx := extractor.GetContext()

	// process each record individually by ID
	for _, id := range ids {
		select {
		case <-ctx.Done():
			return errors.Convert(ctx.Err())
		default:
		}

		// load full record by ID
		row := &RawData{}
		err := db.First(row, dal.From(table), dal.Where("id = ?", id))
		if err != nil {
			return errors.Default.Wrap(err, "error loading full row by ID")
		}

		body := new(InputType)
		err = errors.Convert(json.Unmarshal(row.Data, body))
		if err != nil {
			return err
		}

		if extractor.BeforeExtract != nil {
			err = extractor.BeforeExtract(body, extractor.SubtaskStateManager)
			if err != nil {
				return err
			}
		}

		results, err := extractor.Extract(body, row)
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
	// save the incremental state
	return extractor.SubtaskStateManager.Close()
}

var _ plugin.SubTask = (*StatefulApiExtractor[any])(nil)
