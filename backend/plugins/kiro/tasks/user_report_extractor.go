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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/kiro/models"
)

var _ plugin.SubTaskEntryPoint = ExtractKiroUserReport

var ExtractKiroUserReportMeta = plugin.SubTaskMeta{
	Name:             "extractKiroUserReport",
	EntryPoint:       ExtractKiroUserReport,
	EnabledByDefault: true,
	Description:      "Extract daily per-user activity from Kiro report CSVs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
	Dependencies:     []*plugin.SubTaskMeta{&CollectKiroS3FilesMeta},
}

// ExtractKiroUserReport loads the report CSVs discovered for this scope.
//
// Reports are written once per day per client type, so there are only a few
// hundred per year - the concurrency that matters for logs is irrelevant here,
// but reusing extractFiles keeps the retry and bookkeeping behaviour identical
// across all three streams.
func ExtractKiroUserReport(taskCtx plugin.SubTaskContext) errors.Error {
	return extractFiles(taskCtx, models.FileTypeReport, parseUserReportRows)
}

// parseUserReportRows adapts ParseUserReport to the extractor's batch interface.
//
// Both tables come from one parse because the per-model counts are columns of the
// same CSV row; splitting them into two passes would mean reading every file
// twice. They are returned as two batches rather than one mixed slice because
// GORM resolves the target table from the slice's element type.
func parseUserReportRows(data []byte, connectionId uint64, scopeId string) ([]rowBatch, errors.Error) {
	reports, modelMessages, err := ParseUserReport(data, connectionId, scopeId)
	if err != nil {
		return nil, err
	}

	return []rowBatch{
		{rows: reports, count: len(reports)},
		{rows: modelMessages, count: len(modelMessages)},
	}, nil
}
