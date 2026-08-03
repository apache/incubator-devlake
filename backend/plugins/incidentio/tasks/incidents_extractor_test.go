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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/apache/incubator-devlake/plugins/incidentio/models"
)

const baseHappyPathActive = `{
	"id": "01FCNDV6P870EA6S7TK1DSYDG0",
	"reference": "INC-61",
	"name": "db outage",
	"summary": "replica lag blew past threshold",
	"permalink": "https://app.incident.io/example/incidents/61",
	"mode": "standard",
	"created_at": "2026-05-10T09:55:00Z",
	"updated_at": "2026-05-10T10:05:00Z",
	"incident_status": {"name": "Investigating", "category": "active"},
	"severity": {"id": "sev-uuid-1", "name": "Major", "rank": 2},
	"incident_type": {"id": "type_01", "name": "Default"},
	"creator": {"user": {"id": "usr_100", "name": "Reporter One", "email": "reporter@example.com"}},
	"incident_timestamp_values": [
		{"incident_timestamp": {"id": "ts_01", "name": "Reported at"}, "value": {"value": "2026-05-10T09:55:00Z"}},
		{"incident_timestamp": {"id": "ts_02", "name": "Declared at"}, "value": {"value": "2026-05-10T10:00:00Z"}}
	]
}`

func newTestOptions() *IncidentioOptions {
	return &IncidentioOptions{
		ConnectionId:   7,
		IncidentTypeId: "type_01",
	}
}

func TestExtractIncidentioIncident_HappyPathActive(t *testing.T) {
	op := newTestOptions()
	results, err := extractIncidentioIncident([]byte(baseHappyPathActive), op)
	require.NoError(t, err)
	require.Len(t, results, 1)

	incident, ok := results[0].(*models.Incident)
	require.True(t, ok, "first result should be *models.Incident")
	assert.Equal(t, uint64(7), incident.ConnectionId)
	assert.Equal(t, "01FCNDV6P870EA6S7TK1DSYDG0", incident.Id)
	assert.Equal(t, "INC-61", incident.Reference)
	assert.Equal(t, "db outage", incident.Name)
	assert.Equal(t, "replica lag blew past threshold", incident.Summary)
	assert.Equal(t, "https://app.incident.io/example/incidents/61", incident.Url)
	assert.Equal(t, "standard", incident.Mode)
	assert.Equal(t, "Investigating", incident.StatusName)
	assert.Equal(t, "active", incident.StatusCategory)
	assert.Equal(t, "Major", incident.SeverityName)
	assert.Equal(t, int64(2), incident.SeverityRank)
	assert.Equal(t, "type_01", incident.IncidentTypeId)
	assert.Equal(t, time.Date(2026, 5, 10, 9, 55, 0, 0, time.UTC), incident.CreatedDate)
	assert.Equal(t, time.Date(2026, 5, 10, 10, 5, 0, 0, time.UTC), incident.UpdatedDate)
	assert.Equal(t, time.Date(2026, 5, 10, 10, 0, 0, 0, time.UTC), incident.DeclaredDate)
	assert.Nil(t, incident.ResolvedDate)
}

func TestExtractIncidentioIncident_Resolved(t *testing.T) {
	raw := []byte(`{
		"id": "inc_02",
		"reference": "INC-62",
		"name": "cache cleared",
		"mode": "standard",
		"created_at": "2026-05-09T07:55:00Z",
		"updated_at": "2026-05-09T09:01:00Z",
		"incident_status": {"name": "Closed", "category": "closed"},
		"severity": {"id": "sev-uuid-3", "name": "Minor", "rank": 1},
		"incident_type": {"id": "type_01", "name": "Default"},
		"incident_timestamp_values": [
			{"incident_timestamp": {"id": "ts_02", "name": "Declared at"}, "value": {"value": "2026-05-09T08:00:00Z"}},
			{"incident_timestamp": {"id": "ts_03", "name": "Fixed at"}, "value": {"value": "2026-05-09T08:30:00Z"}},
			{"incident_timestamp": {"id": "ts_04", "name": "Resolved at"}, "value": {"value": "2026-05-09T09:00:00Z"}},
			{"incident_timestamp": {"id": "ts_05", "name": "Closed at"}, "value": {"value": "2026-05-09T09:01:00Z"}}
		]
	}`)
	op := newTestOptions()
	results, err := extractIncidentioIncident(raw, op)
	require.NoError(t, err)
	require.Len(t, results, 1)

	incident := results[0].(*models.Incident)
	assert.Equal(t, "Closed", incident.StatusName)
	assert.Equal(t, "closed", incident.StatusCategory)
	assert.Equal(t, time.Date(2026, 5, 9, 8, 0, 0, 0, time.UTC), incident.DeclaredDate)
	require.NotNil(t, incident.ResolvedDate)
	assert.Equal(t, time.Date(2026, 5, 9, 9, 0, 0, 0, time.UTC), *incident.ResolvedDate)
}

func TestExtractIncidentioIncident_ClosedAtFallback(t *testing.T) {
	raw := []byte(`{
		"id": "inc_03",
		"reference": "INC-63",
		"name": "closed without resolved timestamp",
		"mode": "standard",
		"created_at": "2026-05-09T07:55:00Z",
		"updated_at": "2026-05-09T09:01:00Z",
		"incident_status": {"name": "Closed", "category": "closed"},
		"incident_type": {"id": "type_01", "name": "Default"},
		"incident_timestamp_values": [
			{"incident_timestamp": {"id": "ts_04", "name": "Resolved at"}, "value": null},
			{"incident_timestamp": {"id": "ts_05", "name": "Closed at"}, "value": {"value": "2026-05-09T09:01:00Z"}}
		]
	}`)
	op := newTestOptions()
	results, err := extractIncidentioIncident(raw, op)
	require.NoError(t, err)
	require.Len(t, results, 1)

	incident := results[0].(*models.Incident)
	require.NotNil(t, incident.ResolvedDate)
	assert.Equal(t, time.Date(2026, 5, 9, 9, 1, 0, 0, time.UTC), *incident.ResolvedDate)
}

func TestExtractIncidentioIncident_DeclaredAtFallsBackToCreatedAt(t *testing.T) {
	raw := []byte(`{
		"id": "inc_04",
		"reference": "INC-64",
		"name": "no declared timestamp",
		"mode": "standard",
		"created_at": "2026-05-10T12:00:00Z",
		"updated_at": "2026-05-10T12:05:00Z",
		"incident_status": {"name": "Investigating", "category": "active"},
		"incident_type": {"id": "type_01", "name": "Default"}
	}`)
	op := newTestOptions()
	results, err := extractIncidentioIncident(raw, op)
	require.NoError(t, err)
	require.Len(t, results, 1)
	incident := results[0].(*models.Incident)
	assert.Equal(t, time.Date(2026, 5, 10, 12, 0, 0, 0, time.UTC), incident.DeclaredDate)
	assert.Nil(t, incident.ResolvedDate)
}

func TestExtractIncidentioIncident_NullSeverity(t *testing.T) {
	raw := []byte(`{
		"id": "inc_05",
		"reference": "INC-65",
		"name": "no sev yet",
		"mode": "standard",
		"created_at": "2026-05-10T14:00:00Z",
		"updated_at": "2026-05-10T14:05:00Z",
		"incident_status": {"name": "Triage", "category": "triage"},
		"severity": null,
		"incident_type": {"id": "type_01", "name": "Default"}
	}`)
	op := newTestOptions()
	results, err := extractIncidentioIncident(raw, op)
	require.NoError(t, err)
	require.Len(t, results, 1)
	incident := results[0].(*models.Incident)
	assert.Equal(t, "", incident.SeverityName)
	assert.Equal(t, int64(0), incident.SeverityRank)
}

func TestExtractIncidentioIncident_TestModeSkipped(t *testing.T) {
	raw := []byte(`{
		"id": "inc_test",
		"reference": "INC-66",
		"name": "practice run",
		"mode": "test",
		"created_at": "2026-05-10T15:00:00Z",
		"updated_at": "2026-05-10T15:05:00Z",
		"incident_type": {"id": "type_01", "name": "Default"}
	}`)
	op := newTestOptions()
	results, err := extractIncidentioIncident(raw, op)
	require.NoError(t, err)
	assert.Empty(t, results, "test-mode incident should produce no rows")
}

func TestExtractIncidentioIncident_TutorialModeSkipped(t *testing.T) {
	raw := []byte(`{
		"id": "inc_tutorial",
		"reference": "INC-67",
		"name": "onboarding walkthrough",
		"mode": "tutorial",
		"created_at": "2026-05-10T15:00:00Z",
		"updated_at": "2026-05-10T15:05:00Z",
		"incident_type": {"id": "type_01", "name": "Default"}
	}`)
	op := newTestOptions()
	results, err := extractIncidentioIncident(raw, op)
	require.NoError(t, err)
	assert.Empty(t, results, "tutorial-mode incident should produce no rows")
}

func TestExtractIncidentioIncident_RetrospectiveModeKept(t *testing.T) {
	raw := []byte(`{
		"id": "inc_retro",
		"reference": "INC-68",
		"name": "backfilled incident",
		"mode": "retrospective",
		"created_at": "2026-05-10T16:00:00Z",
		"updated_at": "2026-05-10T16:05:00Z",
		"incident_status": {"name": "Closed", "category": "closed"},
		"incident_type": {"id": "type_01", "name": "Default"}
	}`)
	op := newTestOptions()
	results, err := extractIncidentioIncident(raw, op)
	require.NoError(t, err)
	require.Len(t, results, 1)
	incident := results[0].(*models.Incident)
	assert.Equal(t, "retrospective", incident.Mode)
}

func TestExtractIncidentioIncident_WrongIncidentTypeSkipped(t *testing.T) {
	raw := []byte(`{
		"id": "inc_wrong_type",
		"reference": "INC-69",
		"name": "other type",
		"mode": "standard",
		"created_at": "2026-05-10T18:00:00Z",
		"updated_at": "2026-05-10T18:05:00Z",
		"incident_type": {"id": "type_99", "name": "Security"}
	}`)
	op := newTestOptions()
	results, err := extractIncidentioIncident(raw, op)
	require.NoError(t, err)
	assert.Empty(t, results, "incident for unrelated incident type should produce no rows")
}

func TestExtractIncidentioIncident_MissingIncidentTypeSkippedWhenScoped(t *testing.T) {
	raw := []byte(`{
		"id": "inc_no_type",
		"reference": "INC-70",
		"name": "type omitted",
		"mode": "standard",
		"created_at": "2026-05-10T19:00:00Z",
		"updated_at": "2026-05-10T19:05:00Z"
	}`)
	op := newTestOptions()
	results, err := extractIncidentioIncident(raw, op)
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestExtractIncidentioIncident_GlobalCollectionKeepsAllTypes(t *testing.T) {
	raw := []byte(`{
		"id": "inc_any_type",
		"reference": "INC-71",
		"name": "any type",
		"mode": "standard",
		"created_at": "2026-05-10T20:00:00Z",
		"updated_at": "2026-05-10T20:05:00Z",
		"incident_type": {"id": "type_99", "name": "Security"}
	}`)
	op := &IncidentioOptions{ConnectionId: 7}
	results, err := extractIncidentioIncident(raw, op)
	require.NoError(t, err)
	require.Len(t, results, 1)
}

func TestExtractIncidentioIncident_MissingCreatedAtReturnsError(t *testing.T) {
	raw := []byte(`{
		"id": "inc_bad",
		"reference": "INC-72",
		"name": "bad row",
		"mode": "standard",
		"updated_at": "2026-05-10T20:05:00Z",
		"incident_type": {"id": "type_01", "name": "Default"}
	}`)
	op := newTestOptions()
	_, err := extractIncidentioIncident(raw, op)
	assert.Error(t, err)
}
