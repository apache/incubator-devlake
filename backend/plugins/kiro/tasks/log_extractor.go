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
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/plugins/kiro/models"
)

var (
	_ plugin.SubTaskEntryPoint = ExtractKiroChatLog
	_ plugin.SubTaskEntryPoint = ExtractKiroCompletionLog
)

var ExtractKiroChatLogMeta = plugin.SubTaskMeta{
	Name:             "extractKiroChatLog",
	EntryPoint:       ExtractKiroChatLog,
	EnabledByDefault: true,
	Description:      "Extract Kiro chat interactions from GenerateAssistantResponse logs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
	Dependencies:     []*plugin.SubTaskMeta{&CollectKiroS3FilesMeta},
}

var ExtractKiroCompletionLogMeta = plugin.SubTaskMeta{
	Name:             "extractKiroCompletionLog",
	EntryPoint:       ExtractKiroCompletionLog,
	EnabledByDefault: true,
	Description:      "Extract Kiro inline suggestions from GenerateCompletions logs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
	Dependencies:     []*plugin.SubTaskMeta{&CollectKiroS3FilesMeta},
}

// ExtractKiroChatLog loads the chat interaction logs for this scope.
//
// This is the high-volume stream: one object per interaction, roughly 700 bytes
// each, which for a single active user reaches several hundred objects a day.
// The cost is request count rather than bytes, which is what the worker pool
// addresses.
func ExtractKiroChatLog(taskCtx plugin.SubTaskContext) errors.Error {
	return extractFiles(taskCtx, models.FileTypeChatLog, parseChatLogRows)
}

// ExtractKiroCompletionLog loads the inline suggestion logs for this scope.
//
// This stream is dormant under agentic usage - the sampled history stops in
// March - but the objects remain in S3, so a backfill picks them up and a team
// that enables IDE inline completion starts producing them again.
func ExtractKiroCompletionLog(taskCtx plugin.SubTaskContext) errors.Error {
	return extractFiles(taskCtx, models.FileTypeCompletionLog, parseCompletionLogRows)
}

// Both adapters return a single typed batch: the slice keeps its concrete
// element type so GORM can resolve the table.
func parseChatLogRows(data []byte, connectionId uint64, scopeId string) ([]rowBatch, errors.Error) {
	logs, err := ParseChatLog(data, connectionId, scopeId)
	if err != nil {
		return nil, err
	}
	return []rowBatch{{rows: logs, count: len(logs)}}, nil
}

func parseCompletionLogRows(data []byte, connectionId uint64, scopeId string) ([]rowBatch, errors.Error) {
	logs, err := ParseCompletionLog(data, connectionId, scopeId)
	if err != nil {
		return nil, err
	}
	return []rowBatch{{rows: logs, count: len(logs)}}, nil
}
