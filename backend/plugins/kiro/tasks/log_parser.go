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
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/plugins/kiro/models"
)

// Markers used by the prompt heuristics. These read a prompt's text for signs
// of how the developer invoked Kiro; they are only meaningful when a prompt is
// actually present.
const (
	steeringMarker = ".kiro/steering"
	specMarker     = ".kiro/specs"
)

// Every field below is a pointer or checked for presence because Kiro omits
// empty keys entirely rather than emitting a null. Across 4530 sampled records,
// five documented fields never appeared at all and modelId appeared on only 45%,
// so "absent" and "zero" have to stay distinguishable.
type chatLogFile struct {
	Records []struct {
		Request *struct {
			Prompt           *string `json:"prompt"`
			ChatTriggerType  *string `json:"chatTriggerType"`
			UserId           *string `json:"userId"`
			TimeStamp        *string `json:"timeStamp"`
			ModelId          *string `json:"modelId"`
			CustomizationArn *string `json:"customizationArn"`
		} `json:"generateAssistantResponseEventRequest"`
		Response *struct {
			AssistantResponse *string `json:"assistantResponse"`
			FollowupPrompts   *string `json:"followupPrompts"`
			RequestId         *string `json:"requestId"`
			MessageMetadata   *struct {
				ConversationId *string `json:"conversationId"`
				UtteranceId    *string `json:"utteranceId"`
			} `json:"messageMetadata"`
		} `json:"generateAssistantResponseEventResponse"`
	} `json:"records"`
}

type completionLogFile struct {
	Records []struct {
		Request *struct {
			LeftContext      *string `json:"leftContext"`
			RightContext     *string `json:"rightContext"`
			FileName         *string `json:"fileName"`
			UserId           *string `json:"userId"`
			TimeStamp        *string `json:"timeStamp"`
			CustomizationArn *string `json:"customizationArn"`
		} `json:"generateCompletionsEventRequest"`
		Response *struct {
			Completions []string `json:"completions"`
			RequestId   *string  `json:"requestId"`
		} `json:"generateCompletionsEventResponse"`
	} `json:"records"`
}

// ParseChatLog parses a gzipped GenerateAssistantResponse log.
//
// Neither the prompt nor the assistant response text is persisted: derived
// features are computed here and the originals are dropped, which keeps
// proprietary code and personal content out of the warehouse.
//
// A file holds one or two records, so the array is always walked.
func ParseChatLog(gzData []byte, connectionId uint64, scopeId string) ([]*models.KiroChatLog, errors.Error) {
	raw, err := gunzip(gzData)
	if err != nil {
		return nil, err
	}

	var parsed chatLogFile
	if jsonErr := json.Unmarshal(raw, &parsed); jsonErr != nil {
		return nil, errors.Default.Wrap(jsonErr, "failed to unmarshal chat log")
	}

	var result []*models.KiroChatLog
	for _, record := range parsed.Records {
		if record.Request == nil || record.Response == nil {
			// Without both halves there is no usable interaction.
			continue
		}
		requestId := deref(record.Response.RequestId)
		if requestId == "" {
			// requestId is the primary key; a record without one cannot be
			// stored or deduplicated.
			continue
		}

		userId, identityStoreId := SplitUserId(deref(record.Request.UserId))

		timestamp, tsErr := ParseKiroTime(deref(record.Request.TimeStamp))
		if tsErr != nil {
			return nil, tsErr
		}

		log := &models.KiroChatLog{
			ConnectionId:    connectionId,
			ScopeId:         scopeId,
			RequestId:       requestId,
			UserId:          userId,
			IdentityStoreId: identityStoreId,
			Timestamp:       timestamp,
			ChatTriggerType: deref(record.Request.ChatTriggerType),
			// Only ~45% of records carry a model id, so it stays nil when
			// absent. Attributing model usage from this column would skew any
			// share-of-usage figure; user_model_messages is authoritative.
			ModelId:        record.Request.ModelId,
			ResponseLength: len(deref(record.Response.AssistantResponse)),
			// Presence, not content: the follow-up text itself is not stored.
			HasFollowupPrompts: record.Response.FollowupPrompts != nil,
		}

		prompt := deref(record.Request.Prompt)
		if prompt != "" {
			// A non-empty prompt means the user spoke this turn.
			sum := sha256.Sum256([]byte(prompt))
			hash := hex.EncodeToString(sum[:])
			hasSteering := strings.Contains(prompt, steeringMarker)
			isSpecMode := strings.Contains(prompt, specMarker)

			log.HasPrompt = true
			log.PromptLength = len(prompt)
			log.PromptSha256 = &hash
			log.HasSteering = &hasSteering
			log.IsSpecMode = &isSpecMode
		}
		// When the prompt is empty the agent continued on its own. The hash and
		// both heuristics stay nil: hashing the empty string would give roughly
		// 71% of rows one identical hash and destroy the rework signal, and a
		// false heuristic would be indistinguishable from a real negative.

		if md := record.Response.MessageMetadata; md != nil {
			log.ConversationId = md.ConversationId
			log.UtteranceId = md.UtteranceId
		}

		result = append(result, log)
	}

	return result, nil
}

// ParseCompletionLog parses a gzipped GenerateCompletions log.
//
// The counters are named Returned* rather than Accepted*: a record is written
// when the suggestion is requested, and an empty completions array is common, so
// these measure what Kiro offered rather than what was taken. Records with no
// completions are still stored - they are the denominator.
func ParseCompletionLog(gzData []byte, connectionId uint64, scopeId string) ([]*models.KiroCompletionLog, errors.Error) {
	raw, err := gunzip(gzData)
	if err != nil {
		return nil, err
	}

	var parsed completionLogFile
	if jsonErr := json.Unmarshal(raw, &parsed); jsonErr != nil {
		return nil, errors.Default.Wrap(jsonErr, "failed to unmarshal completion log")
	}

	var result []*models.KiroCompletionLog
	for _, record := range parsed.Records {
		if record.Request == nil || record.Response == nil {
			continue
		}
		requestId := deref(record.Response.RequestId)
		if requestId == "" {
			continue
		}

		userId, identityStoreId := SplitUserId(deref(record.Request.UserId))

		timestamp, tsErr := ParseKiroTime(deref(record.Request.TimeStamp))
		if tsErr != nil {
			return nil, tsErr
		}

		// fileName is a bare file name with no directory path, so it cannot be
		// resolved to a unique repository file.
		fileName := deref(record.Request.FileName)

		charCount := 0
		lineCount := 0
		for _, completion := range record.Response.Completions {
			charCount += len(completion)
			lineCount += strings.Count(completion, "\n") + 1
		}

		result = append(result, &models.KiroCompletionLog{
			ConnectionId:    connectionId,
			ScopeId:         scopeId,
			RequestId:       requestId,
			UserId:          userId,
			IdentityStoreId: identityStoreId,
			Timestamp:       timestamp,
			FileName:        fileName,
			FileExtension:   strings.TrimPrefix(filepath.Ext(fileName), "."),
			// Present on completion records but never on chat records.
			HasCustomization:   record.Request.CustomizationArn != nil,
			CompletionsCount:   len(record.Response.Completions),
			ReturnedCharCount:  charCount,
			ReturnedLineCount:  lineCount,
			LeftContextLength:  len(deref(record.Request.LeftContext)),
			RightContextLength: len(deref(record.Request.RightContext)),
		})
	}

	return result, nil
}

func gunzip(gzData []byte) ([]byte, errors.Error) {
	reader, err := gzip.NewReader(bytes.NewReader(gzData))
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to open gzip reader")
	}
	defer reader.Close()

	raw, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to decompress log file")
	}
	return raw, nil
}

func deref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
