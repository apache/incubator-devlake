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

	"github.com/apache/devlake/plugins/github/models"
	"github.com/stretchr/testify/assert"
)

func TestResolveIssueFieldMappingTrimsAndTolueratesNilConfig(t *testing.T) {
	assert.True(t, resolveIssueFieldMapping(nil).isEmpty())

	mapping := resolveIssueFieldMapping(&models.GithubScopeConfig{
		IssueFieldPriority:   "  Priority  ",
		IssueFieldStoryPoint: "Effort",
	})
	assert.Equal(t, "Priority", mapping.priority)
	assert.Equal(t, "Effort", mapping.storyPoint)
	assert.False(t, mapping.isEmpty())
}

func TestIssueFieldMappingIsEmpty(t *testing.T) {
	assert.True(t, issueFieldMapping{}.isEmpty())
	assert.False(t, issueFieldMapping{priority: "Priority"}.isEmpty())
	assert.False(t, issueFieldMapping{dueDate: "Target date"}.isEmpty())
}

func TestIssueFieldMappingNamesAreLowercasedAndDeduplicated(t *testing.T) {
	mapping := issueFieldMapping{
		priority:   "Priority",
		severity:   "priority", // one field driving two columns
		component:  "",
		storyPoint: "Effort",
		dueDate:    "Target Date",
	}
	assert.Equal(t, []string{"priority", "effort", "target date"}, mapping.names())
}

func TestIssueFieldValuesLookupIsCaseInsensitiveAndNilSafe(t *testing.T) {
	var missing issueFieldValues
	_, ok := missing.lookup(1, "Priority")
	assert.False(t, ok, "nil map must not panic")

	values := issueFieldValues{7: {"priority": "Critical"}}

	text, ok := values.lookup(7, "PRIORITY")
	assert.True(t, ok)
	assert.Equal(t, "Critical", text)

	_, ok = values.lookup(7, "")
	assert.False(t, ok, "an unset mapping must never match")

	_, ok = values.lookup(8, "Priority")
	assert.False(t, ok, "another issue's values must not leak")
}

func TestParseFieldDateAcceptsTheDocumentedAndTimestampForms(t *testing.T) {
	expected := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	for _, text := range []string{"2024-12-31", "2024-12-31T00:00:00Z", "2024-12-31T00:00:00+00:00"} {
		parsed, err := parseFieldDate(text)
		assert.NoError(t, err, text)
		assert.Equal(t, expected, parsed.UTC(), text)
	}
}

func TestParseFieldDateRejectsNonDates(t *testing.T) {
	_, err := parseFieldDate("soon")
	assert.Error(t, err)

	_, err = parseFieldDate("")
	assert.Error(t, err)
}
