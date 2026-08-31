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
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/plugins/kiro/models"
)

// Known column names in the user activity report. Only the first eleven are
// stable; everything after Total_Messages varies by report version.
const (
	colDate               = "Date"
	colUserId             = "UserId"
	colClientType         = "Client_Type"
	colChatConversations  = "Chat_Conversations"
	colCreditsUsed        = "Credits_Used"
	colOverageCap         = "Overage_Cap"
	colOverageCreditsUsed = "Overage_Credits_Used"
	colOverageEnabled     = "Overage_Enabled"
	colProfileId          = "ProfileId"
	colSubscriptionTier   = "Subscription_Tier"
	colTotalMessages      = "Total_Messages"
	colNewUser            = "New_User"
	colUserEmail          = "User_Email"
)

// ParseUserReport parses a user activity report CSV into report rows and their
// per-model message counts.
//
// The parser is deliberately tolerant of schema drift. Across the sampled
// history the report has had 19 distinct header layouts: columns were added
// over time, their order moved, and the per-model columns come and go with the
// models a team uses. So:
//
//   - every field is located by column name, never by position
//   - unknown columns are ignored rather than treated as an error
//   - columns that are absent yield NULL, not a zero value
//
// It is a pure function: no S3, no database, no task context. That keeps it
// testable against real sample bytes.
func ParseUserReport(data []byte, connectionId uint64, scopeId string) ([]*models.KiroUserReport, []*models.KiroUserModelMessage, errors.Error) {
	reader := csv.NewReader(bytes.NewReader(data))
	// Row width varies across report versions, so do not enforce a fixed count.
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			// An empty file is not a failure; it simply has no rows.
			return nil, nil, nil
		}
		return nil, nil, errors.Default.Wrap(err, "failed to read user report CSV header")
	}
	for i := range headers {
		headers[i] = strings.TrimSpace(headers[i])
	}

	var reports []*models.KiroUserReport
	var modelMessages []*models.KiroUserModelMessage

	for {
		record, readErr := reader.Read()
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return nil, nil, errors.Default.Wrap(readErr, "failed to read user report CSV row")
		}
		if isBlankRecord(record) {
			continue
		}

		fields := zipHeaders(headers, record)

		report, rowModels, rowErr := buildReportRow(fields, connectionId, scopeId)
		if rowErr != nil {
			return nil, nil, rowErr
		}
		reports = append(reports, report)
		modelMessages = append(modelMessages, rowModels...)
	}

	return reports, modelMessages, nil
}

// buildReportRow converts one name-keyed CSV row into its models.
func buildReportRow(fields map[string]string, connectionId uint64, scopeId string) (*models.KiroUserReport, []*models.KiroUserModelMessage, errors.Error) {
	date, err := ParseReportDate(fields[colDate])
	if err != nil {
		return nil, nil, err
	}

	// Report CSVs usually carry a bare UUID, but some files carry the
	// identity-store prefix. Normalizing both paths keeps one person from
	// splitting into two ids.
	userId, identityStoreId := SplitUserId(fields[colUserId])
	clientType := strings.TrimSpace(fields[colClientType])

	report := &models.KiroUserReport{
		ConnectionId:    connectionId,
		ScopeId:         scopeId,
		UserId:          userId,
		Date:            date,
		ClientType:      clientType,
		IdentityStoreId: identityStoreId,
		ProfileId:       strings.TrimSpace(fields[colProfileId]),
		// Tiers arrive as UPPER_SNAKE in real data but CamelCase in the docs;
		// normalizing means an unrecognized tier still stores cleanly.
		SubscriptionTier: NormalizeTier(fields[colSubscriptionTier]),
	}

	// These two columns were introduced partway through the report's history.
	// When absent they must stay NULL: a zero value would be indistinguishable
	// from real data, and an empty email would look like a failed identity join
	// rather than an unavailable one.
	if raw, ok := fields[colUserEmail]; ok {
		if email := strings.TrimSpace(raw); email != "" {
			report.UserEmail = &email
		}
	}
	if raw, ok := fields[colNewUser]; ok {
		if trimmed := strings.TrimSpace(raw); trimmed != "" {
			isNew, boolErr := ParseKiroBool(trimmed)
			if boolErr != nil {
				return nil, nil, boolErr
			}
			report.IsNewUser = &isNew
		}
	}

	var intErr errors.Error
	if report.ChatConversations, intErr = optionalInt(fields, colChatConversations); intErr != nil {
		return nil, nil, intErr
	}
	if report.TotalMessages, intErr = optionalInt(fields, colTotalMessages); intErr != nil {
		return nil, nil, intErr
	}

	var floatErr errors.Error
	if report.CreditsUsed, floatErr = optionalFloat(fields, colCreditsUsed); floatErr != nil {
		return nil, nil, floatErr
	}
	if report.OverageCap, floatErr = optionalFloat(fields, colOverageCap); floatErr != nil {
		return nil, nil, floatErr
	}
	// Genuinely non-zero in real data, so it is parsed rather than assumed.
	if report.OverageCreditsUsed, floatErr = optionalFloat(fields, colOverageCreditsUsed); floatErr != nil {
		return nil, nil, floatErr
	}

	if raw, ok := fields[colOverageEnabled]; ok {
		if trimmed := strings.TrimSpace(raw); trimmed != "" {
			// Casing drifts between report versions (TRUE vs true).
			enabled, boolErr := ParseKiroBool(trimmed)
			if boolErr != nil {
				return nil, nil, boolErr
			}
			report.OverageEnabled = enabled
		}
	}

	modelMessages, modelErr := buildModelMessages(fields, connectionId, scopeId, userId, date, clientType)
	if modelErr != nil {
		return nil, nil, modelErr
	}

	return report, modelMessages, nil
}

// buildModelMessages extracts the dynamic per-model message count columns.
//
// Total_Messages matches the same naming pattern but is excluded by
// ParseModelColumn - it is the aggregate, and admitting it would create a
// phantom model whose count equals the sum of the real ones.
func buildModelMessages(
	fields map[string]string,
	connectionId uint64,
	scopeId string,
	userId string,
	date time.Time,
	clientType string,
) ([]*models.KiroUserModelMessage, errors.Error) {
	var result []*models.KiroUserModelMessage
	for header, value := range fields {
		modelName, ok := ParseModelColumn(header)
		if !ok {
			continue
		}
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		count, err := strconv.Atoi(trimmed)
		if err != nil {
			return nil, errors.Default.Wrap(err, "failed to parse message count for model "+modelName)
		}
		result = append(result, &models.KiroUserModelMessage{
			ConnectionId: connectionId,
			ScopeId:      scopeId,
			UserId:       userId,
			Date:         date,
			ClientType:   clientType,
			// Stored with the CSV's underscore spelling; reconciliation with the
			// log's hyphenated form happens at query time.
			ModelName:    modelName,
			MessageCount: count,
		})
	}
	return result, nil
}

// zipHeaders pairs header names with row values.
//
// Short rows are tolerated: a missing trailing column is simply absent from the
// map, which downstream reads as NULL rather than as a zero value.
func zipHeaders(headers, record []string) map[string]string {
	fields := make(map[string]string, len(headers))
	for i, header := range headers {
		if i >= len(record) {
			break
		}
		if header == "" {
			continue
		}
		fields[header] = record[i]
	}
	return fields
}

func isBlankRecord(record []string) bool {
	for _, v := range record {
		if strings.TrimSpace(v) != "" {
			return false
		}
	}
	return true
}

// optionalInt returns 0 when the column is absent or empty. Counters are
// genuinely zero in that case, unlike the nullable identity columns.
func optionalInt(fields map[string]string, column string) (int, errors.Error) {
	raw, ok := fields[column]
	if !ok {
		return 0, nil
	}
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, nil
	}
	v, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0, errors.Default.Wrap(err, "failed to parse integer column "+column)
	}
	return v, nil
}

func optionalFloat(fields map[string]string, column string) (float64, errors.Error) {
	raw, ok := fields[column]
	if !ok {
		return 0, nil
	}
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, nil
	}
	v, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return 0, errors.Default.Wrap(err, "failed to parse float column "+column)
	}
	return v, nil
}
