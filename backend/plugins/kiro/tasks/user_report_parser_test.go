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
	"os"
	"path/filepath"
	"testing"

	"github.com/apache/incubator-devlake/plugins/kiro/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testConnectionId = uint64(1)
	testScopeId      = "123456789012_2026_07"
)

// Fixtures are real report files with only the user id, email and profile
// suffix replaced. Every numeric value is untouched, so the assertions below
// double as a record of what Kiro actually emits.
func loadReportFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", "user_report", name))
	require.NoError(t, err)
	return data
}

func parseReportFixture(t *testing.T, name string) ([]*models.KiroUserReport, []*models.KiroUserModelMessage) {
	t.Helper()
	reports, modelMessages, err := ParseUserReport(loadReportFixture(t, name), testConnectionId, testScopeId)
	require.Nil(t, err)
	return reports, modelMessages
}

// The earliest report layout has neither User_Email nor New_User nor any model
// column. Both identity fields must be NULL rather than empty, because a blank
// email would read as a failed identity join instead of an unavailable one.
func TestParseUserReport_EarliestVariant(t *testing.T) {
	reports, modelMessages := parseReportFixture(t, "01_earliest_no_email_no_newuser.csv")

	require.Len(t, reports, 1)
	r := reports[0]
	assert.Nil(t, r.UserEmail, "User_Email column is absent, so the field must be NULL")
	assert.Nil(t, r.IsNewUser, "New_User column is absent, so the field must be NULL")
	assert.Equal(t, "KIRO_CLI", r.ClientType)
	assert.Equal(t, "POWER", r.SubscriptionTier)
	assert.Equal(t, 4, r.ChatConversations)
	assert.Equal(t, 28, r.TotalMessages)
	assert.InDelta(t, 6.421250216662884, r.CreditsUsed, 1e-12)
	assert.Empty(t, modelMessages, "no model columns in this variant")
}

func TestParseUserReport_StandardCliAndIde(t *testing.T) {
	// Same person, same day, two client types - which is why ClientType is part
	// of the primary key.
	cli, cliModels := parseReportFixture(t, "02_standard_cli.csv")
	require.Len(t, cli, 1)
	assert.Equal(t, "KIRO_CLI", cli[0].ClientType)
	assert.Equal(t, 618, cli[0].TotalMessages)
	assert.Equal(t, 39, cli[0].ChatConversations)
	assert.InDelta(t, 308.97581546497514, cli[0].CreditsUsed, 1e-12)
	require.NotNil(t, cli[0].UserEmail)
	assert.Equal(t, "dev1@example.com", *cli[0].UserEmail)
	require.NotNil(t, cli[0].IsNewUser)
	assert.False(t, *cli[0].IsNewUser)
	require.Len(t, cliModels, 1)
	assert.Equal(t, "claude_opus_5", cliModels[0].ModelName)
	assert.Equal(t, 618, cliModels[0].MessageCount)

	ide, ideModels := parseReportFixture(t, "03_standard_ide.csv")
	require.Len(t, ide, 1)
	assert.Equal(t, "KIRO_IDE", ide[0].ClientType)
	assert.Equal(t, 30, ide[0].TotalMessages)
	assert.InDelta(t, 19.59644039641791, ide[0].CreditsUsed, 1e-12)
	require.Len(t, ideModels, 1)
	assert.Equal(t, 30, ideModels[0].MessageCount)

	// The two rows share a user and date but differ in client type.
	assert.Equal(t, cli[0].UserId, ide[0].UserId)
	assert.Equal(t, cli[0].Date, ide[0].Date)
	assert.NotEqual(t, cli[0].ClientType, ide[0].ClientType)
}

// In this layout New_User comes after the model columns rather than before
// them, which is why nothing may be located by position.
func TestParseUserReport_NewUserAfterModelColumns(t *testing.T) {
	reports, modelMessages := parseReportFixture(t, "04_newuser_after_models.csv")

	require.Len(t, reports, 1)
	require.NotNil(t, reports[0].IsNewUser)
	assert.False(t, *reports[0].IsNewUser)
	assert.Equal(t, 612, reports[0].TotalMessages)

	require.Len(t, modelMessages, 2)
	counts := map[string]int{}
	for _, m := range modelMessages {
		counts[m.ModelName] = m.MessageCount
	}
	assert.Equal(t, 550, counts["claude_opus_4.6"])
	assert.Equal(t, 62, counts["claude_opus_4.7"])
	// Total_Messages must never surface as a model.
	assert.NotContains(t, counts, "Total")
	assert.NotContains(t, counts, "Total_Messages")
}

// This file carries the identity-store prefix in the CSV's UserId column and
// spells its boolean in upper case. Without prefix stripping the same person
// would exist under two ids and every per-person aggregate would be wrong.
func TestParseUserReport_PrefixedUserIdAndUpperCaseBool(t *testing.T) {
	reports, modelMessages := parseReportFixture(t, "05_prefixed_userid_upper_bool.csv")

	require.Len(t, reports, 1)
	r := reports[0]
	assert.Equal(t, "11111111-1111-4111-8111-111111111111", r.UserId,
		"the d-... prefix must be stripped from the CSV path too")
	assert.Equal(t, "d-1234567890", r.IdentityStoreId,
		"the prefix is retained separately for auditability")
	assert.True(t, r.OverageEnabled, "TRUE in upper case must parse")
	assert.Equal(t, 240, r.TotalMessages)

	require.Len(t, modelMessages, 2)
	counts := map[string]int{}
	for _, m := range modelMessages {
		counts[m.ModelName] = m.MessageCount
		// Model rows must carry the normalized id, or they will not join to the
		// report row.
		assert.Equal(t, r.UserId, m.UserId)
	}
	assert.Equal(t, 212, counts["claude_opus_4.6"])
	// "simple_task" is a routing mode rather than a model, but upstream exposes
	// it through the same column pattern so it is stored the same way.
	assert.Equal(t, 28, counts["simple_task"])
}

func TestParseUserReport_MultipleUsers(t *testing.T) {
	reports, _ := parseReportFixture(t, "06_two_users.csv")

	require.Len(t, reports, 2, "a single file can hold more than one user")
	ids := []string{reports[0].UserId, reports[1].UserId}
	assert.NotEqual(t, ids[0], ids[1])
	for _, r := range reports {
		assert.Equal(t, "KIRO_IDE", r.ClientType)
		assert.NotEmpty(t, r.UserId)
	}
}

// KIRO_WEB is absent from the published documentation but present in real
// exports, so client type is never validated against a fixed set. This file
// also carries a non-zero overage figure.
func TestParseUserReport_KiroWebAndNonZeroOverage(t *testing.T) {
	reports, _ := parseReportFixture(t, "07_kiro_web_nonzero_overage.csv")

	require.Len(t, reports, 1)
	r := reports[0]
	assert.Equal(t, "KIRO_WEB", r.ClientType)
	assert.InDelta(t, 1.281141613432836, r.OverageCreditsUsed, 1e-12,
		"overage credits are genuinely non-zero in real data")
	assert.InDelta(t, 1.281141613432836, r.CreditsUsed, 1e-12)
}

func TestParseUserReport_PluginMultiModel(t *testing.T) {
	reports, modelMessages := parseReportFixture(t, "08_plugin_multi_model.csv")

	require.Len(t, reports, 1)
	assert.Equal(t, "PLUGIN", reports[0].ClientType)

	require.Len(t, modelMessages, 2)
	counts := map[string]int{}
	for _, m := range modelMessages {
		counts[m.ModelName] = m.MessageCount
	}
	assert.Equal(t, 9, counts["auto"])
	assert.Equal(t, 7, counts["claude_sonnet_4.6"])
}

// Every fixture must satisfy the invariant that the per-model counts add up to
// Total_Messages. This catches a parser that drops or double-counts a column.
func TestParseUserReport_ModelCountsSumToTotal(t *testing.T) {
	fixtures := []string{
		"02_standard_cli.csv",
		"03_standard_ide.csv",
		"04_newuser_after_models.csv",
		"05_prefixed_userid_upper_bool.csv",
		"07_kiro_web_nonzero_overage.csv",
		"08_plugin_multi_model.csv",
	}
	for _, name := range fixtures {
		t.Run(name, func(t *testing.T) {
			reports, modelMessages := parseReportFixture(t, name)
			require.Len(t, reports, 1)

			sum := 0
			for _, m := range modelMessages {
				sum += m.MessageCount
			}
			assert.Equal(t, reports[0].TotalMessages, sum,
				"per-model counts must add up to Total_Messages")
		})
	}
}

func TestParseUserReport_EdgeCases(t *testing.T) {
	t.Run("empty input yields no rows and no error", func(t *testing.T) {
		reports, modelMessages, err := ParseUserReport(nil, testConnectionId, testScopeId)
		assert.Nil(t, err)
		assert.Empty(t, reports)
		assert.Empty(t, modelMessages)
	})

	t.Run("header only yields no rows", func(t *testing.T) {
		csvData := []byte("Date,UserId,Client_Type,Total_Messages\n")
		reports, modelMessages, err := ParseUserReport(csvData, testConnectionId, testScopeId)
		assert.Nil(t, err)
		assert.Empty(t, reports)
		assert.Empty(t, modelMessages)
	})

	t.Run("blank trailing line is skipped", func(t *testing.T) {
		// Real exports end with a trailing newline, which the CSV reader
		// surfaces as an empty record on some inputs.
		csvData := []byte("Date,UserId,Client_Type,Total_Messages\n2026-07-27,u1,KIRO_CLI,5\n\n")
		reports, _, err := ParseUserReport(csvData, testConnectionId, testScopeId)
		assert.Nil(t, err)
		assert.Len(t, reports, 1)
	})

	// A brand new column must not break collection - that is the whole point of
	// name-based mapping.
	t.Run("unknown columns are ignored", func(t *testing.T) {
		csvData := []byte("Date,UserId,Client_Type,Total_Messages,Some_Future_Column\n" +
			"2026-07-27,u1,KIRO_CLI,5,whatever\n")
		reports, _, err := ParseUserReport(csvData, testConnectionId, testScopeId)
		assert.Nil(t, err)
		require.Len(t, reports, 1)
		assert.Equal(t, 5, reports[0].TotalMessages)
	})

	t.Run("short row leaves trailing columns absent", func(t *testing.T) {
		csvData := []byte("Date,UserId,Client_Type,Total_Messages,User_Email\n" +
			"2026-07-27,u1,KIRO_CLI,5\n")
		reports, _, err := ParseUserReport(csvData, testConnectionId, testScopeId)
		assert.Nil(t, err)
		require.Len(t, reports, 1)
		assert.Nil(t, reports[0].UserEmail)
	})

	t.Run("empty email stays NULL", func(t *testing.T) {
		csvData := []byte("Date,UserId,Client_Type,Total_Messages,User_Email\n" +
			"2026-07-27,u1,KIRO_CLI,5,\n")
		reports, _, err := ParseUserReport(csvData, testConnectionId, testScopeId)
		assert.Nil(t, err)
		require.Len(t, reports, 1)
		assert.Nil(t, reports[0].UserEmail)
	})

	t.Run("zero model count is still recorded", func(t *testing.T) {
		csvData := []byte("Date,UserId,Client_Type,Total_Messages,auto_messages\n" +
			"2026-07-27,u1,KIRO_CLI,0,0\n")
		_, modelMessages, err := ParseUserReport(csvData, testConnectionId, testScopeId)
		assert.Nil(t, err)
		require.Len(t, modelMessages, 1)
		assert.Equal(t, 0, modelMessages[0].MessageCount)
	})

	t.Run("bad date is an error", func(t *testing.T) {
		csvData := []byte("Date,UserId,Client_Type\n07-27-2026,u1,KIRO_CLI\n")
		_, _, err := ParseUserReport(csvData, testConnectionId, testScopeId)
		assert.NotNil(t, err, "the legacy MM-DD-YYYY format must not parse silently")
	})

	t.Run("bad number is an error", func(t *testing.T) {
		csvData := []byte("Date,UserId,Client_Type,Total_Messages\n2026-07-27,u1,KIRO_CLI,not-a-number\n")
		_, _, err := ParseUserReport(csvData, testConnectionId, testScopeId)
		assert.NotNil(t, err)
	})

	t.Run("bad boolean is an error", func(t *testing.T) {
		csvData := []byte("Date,UserId,Client_Type,Overage_Enabled\n2026-07-27,u1,KIRO_CLI,maybe\n")
		_, _, err := ParseUserReport(csvData, testConnectionId, testScopeId)
		assert.NotNil(t, err)
	})

	t.Run("connection and scope are propagated", func(t *testing.T) {
		csvData := []byte("Date,UserId,Client_Type,Total_Messages,auto_messages\n" +
			"2026-07-27,u1,KIRO_CLI,5,5\n")
		reports, modelMessages, err := ParseUserReport(csvData, 42, "scope-x")
		assert.Nil(t, err)
		require.Len(t, reports, 1)
		require.Len(t, modelMessages, 1)
		assert.Equal(t, uint64(42), reports[0].ConnectionId)
		assert.Equal(t, "scope-x", reports[0].ScopeId)
		assert.Equal(t, uint64(42), modelMessages[0].ConnectionId)
		assert.Equal(t, "scope-x", modelMessages[0].ScopeId)
	})
}
