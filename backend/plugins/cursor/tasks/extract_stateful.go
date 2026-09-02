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
	"encoding/json"
	"reflect"
	"time"

	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/models/common"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
)

// cursorExtractRepairGapWarnThreshold is the unpromoted raw-id gap above which
// extract logs a Warn (small gaps are normal right after collect in the same pipeline).
const cursorExtractRepairGapWarnThreshold = 100

// cursorStatefulExtractorArgs mirrors helper.StatefulApiExtractorArgs with an
// extra tool-table repair path for unpromoted raw rows.
type cursorStatefulExtractorArgs[InputType any] struct {
	*helper.SubtaskCommonArgs
	Extract      func(body *InputType, row *helper.RawData) ([]any, errors.Error)
	ConnectionId uint64
	ScopeId      string
	ToolTable    string
}

type cursorStatefulExtractor[InputType any] struct {
	*cursorStatefulExtractorArgs[InputType]
	*helper.SubtaskStateManager
}

func newCursorStatefulExtractor[InputType any](args *cursorStatefulExtractorArgs[InputType]) (*cursorStatefulExtractor[InputType], errors.Error) {
	stateManager, err := helper.NewSubtaskStateManager(args.SubtaskCommonArgs)
	if err != nil {
		return nil, err
	}
	return &cursorStatefulExtractor[InputType]{
		cursorStatefulExtractorArgs: args,
		SubtaskStateManager:         stateManager,
	}, nil
}

// maxPromotedRawId returns COALESCE(MAX(_raw_data_id), 0) for the tool table
// scoped to connection/scope. Used to repair raw rows that were collected but
// never promoted to the tool layer.
func maxPromotedRawId(db dal.Dal, toolTable string, connectionId uint64, scopeId string) (uint64, errors.Error) {
	if toolTable == "" || !db.HasTable(toolTable) {
		return 0, nil
	}
	var maxIds []uint64
	err := db.Pluck(
		"COALESCE(MAX(_raw_data_id), 0)",
		&maxIds,
		dal.From(toolTable),
		dal.Where("connection_id = ? AND scope_id = ?", connectionId, scopeId),
	)
	if err != nil {
		return 0, errors.Default.Wrap(err, "failed to query max promoted raw id")
	}
	if len(maxIds) == 0 {
		return 0, nil
	}
	return maxIds[0], nil
}

// maxRawId returns COALESCE(MAX(id), 0) for raw rows matching params.
func maxRawId(db dal.Dal, rawTable, params string) (uint64, errors.Error) {
	if rawTable == "" || !db.HasTable(rawTable) {
		return 0, nil
	}
	var maxIds []uint64
	err := db.Pluck(
		"COALESCE(MAX(id), 0)",
		&maxIds,
		dal.From(rawTable),
		dal.Where("params = ?", params),
	)
	if err != nil {
		return 0, errors.Default.Wrap(err, "failed to query max raw id")
	}
	if len(maxIds) == 0 {
		return 0, nil
	}
	return maxIds[0], nil
}

// extractRepairGap returns how far raw MAX(id) is ahead of the tool's MAX(_raw_data_id).
func extractRepairGap(maxRaw, maxPromoted uint64) uint64 {
	if maxRaw <= maxPromoted {
		return 0
	}
	return maxRaw - maxPromoted
}

// countIncrementalRepairBreakdown returns counts for repair-only backlog rows
// (id > maxPromoted AND created_at < since) and normal since-window rows.
func countIncrementalRepairBreakdown(db dal.Dal, rawTable, params string, since *time.Time, until *time.Time, maxPromoted uint64) (repairOnlyRows, sinceRows int64, err errors.Error) {
	base := []dal.Clause{
		dal.From(rawTable),
		dal.Where("params = ?", params),
	}
	if until != nil {
		base = append(base, dal.Where("created_at < ?", *until))
	}

	if since != nil {
		sinceClauses := append([]dal.Clause{}, base...)
		sinceClauses = append(sinceClauses, dal.Where("created_at >= ?", *since))
		sinceRows, err = db.Count(sinceClauses...)
		if err != nil {
			return 0, 0, errors.Default.Wrap(err, "failed to count since-window raw rows")
		}

		repairClauses := append([]dal.Clause{}, base...)
		repairClauses = append(repairClauses, dal.Where("id > ? AND created_at < ?", maxPromoted, *since))
		repairOnlyRows, err = db.Count(repairClauses...)
		if err != nil {
			return 0, 0, errors.Default.Wrap(err, "failed to count repair-only raw rows")
		}
		return repairOnlyRows, sinceRows, nil
	}

	// No since: repair path is id > maxPromoted; sinceRows is 0.
	repairClauses := append([]dal.Clause{}, base...)
	repairClauses = append(repairClauses, dal.Where("id > ?", maxPromoted))
	repairOnlyRows, err = db.Count(repairClauses...)
	if err != nil {
		return 0, 0, errors.Default.Wrap(err, "failed to count repair-only raw rows")
	}
	return repairOnlyRows, 0, nil
}

// incrementalRawWhere builds the incremental raw-row filter:
// created_at >= since OR id > maxPromotedRawId (repair gap).
// Exported for unit tests via the package-level function below.
func incrementalRawWhere(since *time.Time, maxPromoted uint64) (clause string, args []any) {
	if since != nil {
		return "(created_at >= ? OR id > ?)", []any{*since, maxPromoted}
	}
	return "id > ?", []any{maxPromoted}
}

func setCursorRawDataOrigin(result any, table, params string, id uint64) {
	if origin, ok := result.(common.GetRawDataOrigin); ok {
		o := origin.GetRawDataOrigin()
		o.RawDataTable = table
		o.RawDataParams = params
		o.RawDataId = id
	}
}

func (extractor *cursorStatefulExtractor[InputType]) Execute() errors.Error {
	db := extractor.GetDal()
	logger := extractor.GetLogger()
	table := extractor.GetRawDataTable()
	params := extractor.GetRawDataParams()
	if !db.HasTable(table) {
		return nil
	}

	if !extractor.IsIncremental() {
		logger.Info(
			"cursor extract full sync: tool=%s raw=%s params=%s (first run, blueprint full refresh, or subtask config change)",
			extractor.ToolTable,
			table,
			params,
		)
	}

	clauses := []dal.Clause{
		dal.Select("id"),
		dal.From(table),
		dal.Where("params = ?", params),
		dal.Orderby("id ASC"),
	}

	var maxPromoted uint64
	if extractor.IsIncremental() {
		var err errors.Error
		maxPromoted, err = maxPromotedRawId(db, extractor.ToolTable, extractor.ConnectionId, extractor.ScopeId)
		if err != nil {
			return err
		}
		whereSQL, whereArgs := incrementalRawWhere(extractor.GetSince(), maxPromoted)
		clauses = append(clauses, dal.Where(whereSQL, whereArgs...))
	}
	clauses = append(clauses, dal.Where("created_at < ? ", extractor.GetUntil()))

	count, err := db.Count(clauses...)
	if err != nil {
		return errors.Default.Wrap(err, "error getting count of records")
	}
	logger.Info("get data from %s where params=%s and got %d with clauses %+v", table, params, count, clauses)

	if extractor.IsIncremental() {
		rawMax, rawErr := maxRawId(db, table, params)
		if rawErr != nil {
			return rawErr
		}
		gap := extractRepairGap(rawMax, maxPromoted)
		repairOnlyRows, sinceRows, breakdownErr := countIncrementalRepairBreakdown(
			db, table, params, extractor.GetSince(), extractor.GetUntil(), maxPromoted,
		)
		if breakdownErr != nil {
			return breakdownErr
		}
		sinceStr := "nil"
		if since := extractor.GetSince(); since != nil {
			sinceStr = since.UTC().Format(time.RFC3339)
		}
		logger.Info(
			"cursor incremental extract: tool=%s maxPromotedRawId=%d maxRawId=%d gap=%d since=%s rowsToProcess=%d repairOnlyRows=%d sinceRows=%d",
			extractor.ToolTable,
			maxPromoted,
			rawMax,
			gap,
			sinceStr,
			count,
			repairOnlyRows,
			sinceRows,
		)
		if gap >= cursorExtractRepairGapWarnThreshold {
			logger.Warn(
				nil,
				"cursor extract repair gap: tool=%s maxPromoted=%d maxRawId=%d unpromoted=%d (threshold=%d)",
				extractor.ToolTable,
				maxPromoted,
				rawMax,
				gap,
				cursorExtractRepairGapWarnThreshold,
			)
		}
	}

	var ids []uint64
	err = db.Pluck("id", &ids, clauses...)
	if err != nil {
		return errors.Default.Wrap(err, "error getting IDs")
	}

	divider := helper.NewBatchSaveDivider(extractor.SubTaskContext, extractor.GetBatchSize(), table, params)
	divider.SetIncrementalMode(extractor.IsIncremental())

	extractor.SetProgress(0, -1)
	ctx := extractor.GetContext()

	for _, id := range ids {
		select {
		case <-ctx.Done():
			return errors.Convert(ctx.Err())
		default:
		}

		row := &helper.RawData{}
		err := db.First(row, dal.From(table), dal.Where("id = ?", id))
		if err != nil {
			return errors.Default.Wrap(err, "error loading full row by ID")
		}

		body := new(InputType)
		err = errors.Convert(json.Unmarshal(row.Data, body))
		if err != nil {
			return err
		}

		results, err := extractor.Extract(body, row)
		if err != nil {
			return errors.Default.Wrap(err, "error calling plugin Extract implementation")
		}

		for _, result := range results {
			batch, err := divider.ForType(reflect.TypeOf(result))
			if err != nil {
				return errors.Default.Wrap(err, "error getting batch from result")
			}
			setCursorRawDataOrigin(result, table, params, row.ID)
			err = batch.Add(result)
			if err != nil {
				return errors.Default.Wrap(err, "error adding result to batch")
			}
		}
		extractor.IncProgress(1)
	}

	err = divider.Close()
	if err != nil {
		return err
	}
	return extractor.SubtaskStateManager.Close()
}

var _ plugin.SubTask = (*cursorStatefulExtractor[any])(nil)
