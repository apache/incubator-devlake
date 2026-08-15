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
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

// issueFieldMapping is the resolved set of GitHub issue field names that should populate
// domain issue columns, taken from the scope config.
//
// The mapping is applied when converting to the domain layer rather than written back into
// _tool_github_issues: the collector iterates that table to build its request URLs, so a
// subtask that both reads and writes it would be a cycle in the subtask graph.
type issueFieldMapping struct {
	priority   string
	severity   string
	component  string
	storyPoint string
	dueDate    string
}

func resolveIssueFieldMapping(scopeConfig *models.GithubScopeConfig) issueFieldMapping {
	if scopeConfig == nil {
		return issueFieldMapping{}
	}
	return issueFieldMapping{
		priority:   strings.TrimSpace(scopeConfig.IssueFieldPriority),
		severity:   strings.TrimSpace(scopeConfig.IssueFieldSeverity),
		component:  strings.TrimSpace(scopeConfig.IssueFieldComponent),
		storyPoint: strings.TrimSpace(scopeConfig.IssueFieldStoryPoint),
		dueDate:    strings.TrimSpace(scopeConfig.IssueFieldDueDate),
	}
}

func (m issueFieldMapping) isEmpty() bool {
	return m.priority == "" && m.severity == "" && m.component == "" &&
		m.storyPoint == "" && m.dueDate == ""
}

// names returns the distinct, lower-cased field names the mapping refers to.
func (m issueFieldMapping) names() []string {
	seen := map[string]bool{}
	var out []string
	for _, name := range []string{m.priority, m.severity, m.component, m.storyPoint, m.dueDate} {
		if name == "" {
			continue
		}
		lower := strings.ToLower(name)
		if !seen[lower] {
			seen[lower] = true
			out = append(out, lower)
		}
	}
	return out
}

// asSubtaskConfig lets the stateful converter re-run in full when a mapping changes, so
// clearing or repointing one is not silently skipped by an incremental run.
func (m issueFieldMapping) asSubtaskConfig() map[string]string {
	return map[string]string{
		"issueFieldPriority":   m.priority,
		"issueFieldSeverity":   m.severity,
		"issueFieldComponent":  m.component,
		"issueFieldStoryPoint": m.storyPoint,
		"issueFieldDueDate":    m.dueDate,
	}
}

// issueFieldValues maps a GitHub issue id to its mapped field values, keyed by lower-cased
// field name.
type issueFieldValues map[int]map[string]string

// loadIssueFieldValues fetches the mapped fields for a connection in one query. The result is
// bounded by (issues x mapped fields), so this stays a single round trip rather than one per
// issue. Returns nil when nothing is mapped.
func loadIssueFieldValues(
	taskCtx plugin.SubTaskContext,
	connectionId uint64,
	mapping issueFieldMapping,
) (issueFieldValues, errors.Error) {
	if mapping.isEmpty() {
		return nil, nil
	}
	db := taskCtx.GetDal()
	var values []models.GithubIssueFieldValue
	err := db.All(&values,
		dal.From(models.GithubIssueFieldValue{}.TableName()),
		dal.Where("connection_id = ? and LOWER(field_name) in ?", connectionId, mapping.names()),
	)
	if err != nil {
		return nil, err
	}
	byIssue := make(issueFieldValues, len(values))
	for _, value := range values {
		if value.Value == "" {
			continue
		}
		fields, ok := byIssue[value.IssueId]
		if !ok {
			fields = map[string]string{}
			byIssue[value.IssueId] = fields
		}
		fields[strings.ToLower(value.FieldName)] = value.Value
	}
	return byIssue, nil
}

// lookup returns the value for a mapped field name on one issue, reporting false when the
// mapping is unset or the issue carries no value for it.
func (v issueFieldValues) lookup(issueId int, fieldName string) (string, bool) {
	if v == nil || fieldName == "" {
		return "", false
	}
	fields, ok := v[issueId]
	if !ok {
		return "", false
	}
	text, ok := fields[strings.ToLower(fieldName)]
	return text, ok
}

// dateLayouts covers the date form the API documents plus the full timestamp form, since a
// date field can round-trip either way.
var dateLayouts = []string{"2006-01-02", time.RFC3339, "2006-01-02T15:04:05Z"}

func parseFieldDate(text string) (*time.Time, error) {
	var lastErr error
	for _, layout := range dateLayouts {
		parsed, err := time.Parse(layout, text)
		if err == nil {
			return &parsed, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

// applyIssueFields overlays the mapped GitHub issue field values onto a domain issue. Values
// that do not parse into their target type are logged and skipped, leaving whatever the label
// regexes produced rather than writing a wrong value.
func applyIssueFields(
	taskCtx plugin.SubTaskContext,
	mapping issueFieldMapping,
	values issueFieldValues,
	issue *models.GithubIssue,
	domainIssue *ticket.Issue,
) {
	if values == nil || mapping.isEmpty() {
		return
	}
	logger := taskCtx.GetLogger()

	if text, ok := values.lookup(issue.GithubId, mapping.priority); ok {
		domainIssue.Priority = text
	}
	if text, ok := values.lookup(issue.GithubId, mapping.severity); ok {
		domainIssue.Severity = text
	}
	if text, ok := values.lookup(issue.GithubId, mapping.component); ok {
		domainIssue.Component = text
	}
	if text, ok := values.lookup(issue.GithubId, mapping.storyPoint); ok {
		if number, err := strconv.ParseFloat(text, 64); err == nil {
			domainIssue.StoryPoint = &number
		} else {
			logger.Warn(nil, "issue #%d: field %q value %q is not a number, story point left unset",
				issue.Number, mapping.storyPoint, text)
		}
	}
	if text, ok := values.lookup(issue.GithubId, mapping.dueDate); ok {
		if parsed, err := parseFieldDate(text); err == nil {
			domainIssue.DueDate = parsed
		} else {
			logger.Warn(nil, "issue #%d: field %q value %q is not a date, due date left unset",
				issue.Number, mapping.dueDate, text)
		}
	}
}
