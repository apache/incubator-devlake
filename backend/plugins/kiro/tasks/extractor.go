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
	"time"

	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
	"github.com/apache/devlake/plugins/kiro/models"
)

// parseFunc turns one downloaded object into batches of rows ready for
// insertion. It is a pure function so it can run on a worker goroutine without
// touching the database.
//
// The result is a slice of batches rather than one flat slice because GORM
// derives the target table from the element type: a mixed []interface{} fails
// with "Table not set". A report CSV yields two different models, so each model
// gets its own homogeneous batch.
type parseFunc func(data []byte, connectionId uint64, scopeId string) ([]rowBatch, errors.Error)

// rowBatch is a set of rows that all belong to the same table.
type rowBatch struct {
	// rows must be a typed slice (e.g. []*models.KiroChatLog), not
	// []interface{}, so GORM can resolve the table.
	rows  interface{}
	count int
}

// parsedFile is a worker's output, handed to the main goroutine for persistence.
type parsedFile struct {
	meta    *models.KiroS3FileMeta
	batches []rowBatch
	err     errors.Error
}

// rowCount totals the rows across every batch, for the cursor's record count.
func (p parsedFile) rowCount() int {
	total := 0
	for _, batch := range p.batches {
		total += batch.count
	}
	return total
}

// pendingFiles returns the files still awaiting extraction for a scope.
//
// Files that have already failed MaxAttempts times are excluded. Without that
// bound, one permanently malformed object is retried on every run: the log fills
// with the same error and the scope never reaches a finished state.
func pendingFiles(db dal.Dal, connectionId uint64, scopeId string, fileType string) ([]models.KiroS3FileMeta, errors.Error) {
	var files []models.KiroS3FileMeta
	err := db.All(&files,
		dal.From(&models.KiroS3FileMeta{}),
		dal.Where("connection_id = ? AND scope_id = ? AND file_type = ? AND processed = ? AND attempt_count < ?",
			connectionId, scopeId, fileType, false, models.MaxAttempts),
	)
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to query pending kiro files")
	}
	return files, nil
}

// extractFiles downloads and parses every pending file of one type, then stores
// the results.
//
// Downloads run concurrently but all database writes happen on the calling
// goroutine. That split is a correctness requirement, not an optimization:
// twenty workers inserting into the same table concurrently deadlock on MySQL
// gap locks, and DevLake's retry layer would turn those deadlocks into
// intermittent failures that only appear under load.
func extractFiles(taskCtx plugin.SubTaskContext, fileType string, parse parseFunc) errors.Error {
	data := taskCtx.GetData().(*KiroTaskData)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()

	files, err := pendingFiles(db, data.Options.ConnectionId, data.Options.ScopeId, fileType)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		logger.Info("no pending %s files for scope %s", fileType, data.Options.ScopeId)
		return nil
	}

	logger.Info("extracting %d %s files", len(files), fileType)
	taskCtx.SetProgress(0, len(files))

	client := data.S3Clients.ForFileType(fileType)
	results := make(chan parsedFile, len(files))

	scheduler, err := helper.NewWorkerScheduler(
		taskCtx.GetContext(),
		data.WorkerCount(),
		// The tick is a global rate limiter, not a per-worker one: every task
		// waits for one tick before running, so this value caps total
		// throughput regardless of pool size. A one-second tick would pin the
		// whole run to one file per second and leave the pool idle - measured
		// at 13 seconds for 13 files. Keep it well below the per-object fetch
		// time so the workers, not the ticker, set the pace.
		time.Millisecond,
		logger,
	)
	if err != nil {
		return err
	}
	defer scheduler.Release()

	for i := range files {
		meta := files[i]
		scheduler.SubmitBlocking(func() errors.Error {
			// Workers only fetch and parse. Anything touching the database
			// happens after WaitAsync below.
			body, getErr := client.GetObjectBytes(meta.S3Path)
			if getErr != nil {
				results <- parsedFile{meta: &meta, err: getErr}
				return nil
			}
			batches, parseErr := parse(body, data.Options.ConnectionId, data.Options.ScopeId)
			results <- parsedFile{meta: &meta, batches: batches, err: parseErr}
			return nil
		})
	}

	if waitErr := scheduler.WaitAsync(); waitErr != nil {
		return waitErr
	}
	close(results)

	for result := range results {
		if saveErr := persistFile(db, result); saveErr != nil {
			return saveErr
		}
		taskCtx.IncProgress(1)
	}

	return nil
}

// persistFile stores one file's rows and updates its cursor entry.
//
// A parse or download failure is recorded rather than aborting the run: one bad
// object should not block the rest of the scope. The reason is written to
// error_message so a short month can be explained with a query instead of
// guesswork.
func persistFile(db dal.Dal, result parsedFile) errors.Error {
	meta := result.meta

	if result.err != nil {
		meta.AttemptCount++
		meta.ErrorMessage = result.err.Error()
		// Processed stays false so the file is retried, up to MaxAttempts.
		if err := db.Update(meta); err != nil {
			return errors.Default.Wrap(err, "failed to record kiro extraction failure")
		}
		return nil
	}

	// Each batch is written separately so GORM sees a typed slice and can
	// resolve the target table.
	for _, batch := range result.batches {
		if batch.count == 0 {
			continue
		}
		if err := db.CreateOrUpdate(batch.rows); err != nil {
			// A write failure counts as an attempt too, otherwise a row that
			// cannot be stored would be fetched forever.
			meta.AttemptCount++
			meta.ErrorMessage = err.Error()
			if updateErr := db.Update(meta); updateErr != nil {
				return errors.Default.Wrap(updateErr, "failed to record kiro write failure")
			}
			return nil
		}
	}

	now := time.Now()
	meta.Processed = true
	meta.ProcessedTime = &now
	meta.RecordCount = result.rowCount()
	meta.ErrorMessage = ""
	if err := db.Update(meta); err != nil {
		return errors.Default.Wrap(err, "failed to mark kiro file processed")
	}
	return nil
}
