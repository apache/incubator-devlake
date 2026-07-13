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
	"testing"
)

func Test_BoardConfiguration_UnmarshalSubQuery(t *testing.T) {
	tests := []struct {
		name             string
		raw              string
		wantSubQuery     string
		wantFilterID     string
		wantID           int
		wantName         string
		wantType         string
		wantColumnCount  int
		wantRankFieldID  int
	}{
		{
			name: "kanban board with sub-filter object",
			raw: `{"id":1201,"name":"Squad 5","type":"kanban",` +
				`"self":"https://example.atlassian.net/rest/agile/1.0/board/1201/configuration",` +
				`"filter":{"id":"17696","self":"https://example.atlassian.net/rest/api/2/filter/17696"},` +
				`"subQuery":{"query":"fixVersion in unreleasedVersions() OR fixVersion is EMPTY"},` +
				`"columnConfig":{"columns":[{"name":"Backlog","statuses":[{"id":"1","self":"https://example.atlassian.net/rest/api/2/status/1"}]},` +
				`{"name":"Done","statuses":[{"id":"10037","self":"https://example.atlassian.net/rest/api/2/status/10037"}]}],` +
				`"constraintType":"issueCount"},"ranking":{"rankCustomFieldId":10019}}`,
			wantSubQuery:    "fixVersion in unreleasedVersions() OR fixVersion is EMPTY",
			wantFilterID:    "17696",
			wantID:          1201,
			wantName:        "Squad 5",
			wantType:        "kanban",
			wantColumnCount: 2,
			wantRankFieldID: 10019,
		},
		{
			name: "board without subQuery field",
			raw: `{"id":500,"name":"No SubFilter Board","type":"scrum",` +
				`"self":"https://example.atlassian.net/rest/agile/1.0/board/500/configuration",` +
				`"filter":{"id":"99999","self":"https://example.atlassian.net/rest/api/2/filter/99999"},` +
				`"columnConfig":{"columns":[],"constraintType":"issueCount"},"ranking":{"rankCustomFieldId":10019}}`,
			wantSubQuery:    "",
			wantFilterID:    "99999",
			wantID:          500,
			wantName:        "No SubFilter Board",
			wantType:        "scrum",
			wantColumnCount: 0,
			wantRankFieldID: 10019,
		},
		{
			name: "board with empty subQuery object",
			raw: `{"id":600,"name":"Empty SubQuery Board","type":"kanban",` +
				`"self":"https://example.atlassian.net/rest/agile/1.0/board/600/configuration",` +
				`"filter":{"id":"11111","self":"https://example.atlassian.net/rest/api/2/filter/11111"},` +
				`"subQuery":{},` +
				`"columnConfig":{"columns":[],"constraintType":"issueCount"},"ranking":{"rankCustomFieldId":10019}}`,
			wantSubQuery:    "",
			wantFilterID:    "11111",
			wantID:          600,
			wantName:        "Empty SubQuery Board",
			wantType:        "kanban",
			wantColumnCount: 0,
			wantRankFieldID: 10019,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bc BoardConfiguration
			if err := json.Unmarshal([]byte(tt.raw), &bc); err != nil {
				t.Fatalf("failed to unmarshal BoardConfiguration: %v", err)
			}
			if bc.SubQuery.Query != tt.wantSubQuery {
				t.Errorf("SubQuery.Query = %q, want %q", bc.SubQuery.Query, tt.wantSubQuery)
			}
			if bc.Filter.ID != tt.wantFilterID {
				t.Errorf("Filter.ID = %q, want %q", bc.Filter.ID, tt.wantFilterID)
			}
			if bc.ID != tt.wantID {
				t.Errorf("ID = %d, want %d", bc.ID, tt.wantID)
			}
			if bc.Name != tt.wantName {
				t.Errorf("Name = %q, want %q", bc.Name, tt.wantName)
			}
			if bc.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", bc.Type, tt.wantType)
			}
			if len(bc.ColumnConfig.Columns) != tt.wantColumnCount {
				t.Errorf("ColumnConfig.Columns length = %d, want %d", len(bc.ColumnConfig.Columns), tt.wantColumnCount)
			}
			if bc.Ranking.RankCustomFieldID != tt.wantRankFieldID {
				t.Errorf("Ranking.RankCustomFieldID = %d, want %d", bc.Ranking.RankCustomFieldID, tt.wantRankFieldID)
			}
		})
	}
}

func Test_BoardConfiguration_FullJiraCloudResponse(t *testing.T) {
	// Exact response payload from Jira Cloud for Board 1201 (Squad 5)
	raw := `{"id":1201,"name":"Squad 5","type":"kanban","self":"https://rakutenadvertising.atlassian.net/rest/agile/1.0/board/1201/configuration","location":{"type":"user","id":"62d8159bb2e6b1992b5be875","self":"https://rakutenadvertising.atlassian.net/rest/api/2/user?accountId=62d8159bb2e6b1992b5be875"},"filter":{"id":"17696","self":"https://rakutenadvertising.atlassian.net/rest/api/2/filter/17696"},"subQuery":{"query":"fixVersion in unreleasedVersions() OR fixVersion is EMPTY"},"columnConfig":{"columns":[{"name":"Backlog","statuses":[{"id":"1","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/1"},{"id":"4","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/4"},{"id":"10016","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/10016"},{"id":"10003","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/10003"}]},{"name":"To Do","statuses":[{"id":"10054","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/10054"}]},{"name":"Blocked","statuses":[{"id":"10019","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/10019"}]},{"name":"In Development","statuses":[{"id":"10017","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/10017"},{"id":"10177","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/10177"},{"id":"10038","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/10038"}]},{"name":"Code Review","statuses":[{"id":"10024","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/10024"}]},{"name":"Ready for QA","statuses":[{"id":"10029","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/10029"},{"id":"10033","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/10033"}]},{"name":"In QA","statuses":[{"id":"10018","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/10018"},{"id":"10158","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/10158"}]},{"name":"Done","statuses":[{"id":"10037","self":"https://rakutenadvertising.atlassian.net/rest/api/2/status/10037"}]}],"constraintType":"issueCount"},"ranking":{"rankCustomFieldId":10019}}`

	var bc BoardConfiguration
	if err := json.Unmarshal([]byte(raw), &bc); err != nil {
		t.Fatalf("failed to unmarshal real Jira Cloud response: %v", err)
	}

	if bc.SubQuery.Query != "fixVersion in unreleasedVersions() OR fixVersion is EMPTY" {
		t.Errorf("SubQuery.Query = %q, want the fixVersion sub-filter", bc.SubQuery.Query)
	}
	if bc.Filter.ID != "17696" {
		t.Errorf("Filter.ID = %q, want %q", bc.Filter.ID, "17696")
	}
	if len(bc.ColumnConfig.Columns) != 8 {
		t.Errorf("ColumnConfig.Columns length = %d, want 8", len(bc.ColumnConfig.Columns))
	}
	if bc.Location.ID != "62d8159bb2e6b1992b5be875" {
		t.Errorf("Location.ID = %q, want %q", bc.Location.ID, "62d8159bb2e6b1992b5be875")
	}
}
