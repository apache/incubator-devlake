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
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

// parseClickUpTime parses a ClickUp millisecond-epoch timestamp. ClickUp encodes
// timestamps as strings of milliseconds since the Unix epoch (e.g. "1567780450202").
// It returns nil for empty / zero / unparseable values so callers can leave the
// corresponding *time.Time unset rather than storing a bogus 1970 date.
func parseClickUpTime(ms string) *time.Time {
	ms = strings.TrimSpace(ms)
	if ms == "" {
		return nil
	}
	millis, err := strconv.ParseInt(ms, 10, 64)
	if err != nil || millis <= 0 {
		return nil
	}
	t := time.UnixMilli(millis).UTC()
	return &t
}

// statusFromType maps a ClickUp status.type onto a DevLake standard issue
// status. ClickUp's status types are standardized:
//
//	open, unstarted -> TODO
//	custom          -> IN_PROGRESS
//	done, closed    -> DONE
//
// Any unrecognized type falls back to OTHER so unexpected API values surface
// rather than silently masquerading as a known status.
func statusFromType(statusType string) string {
	switch strings.ToLower(statusType) {
	case "open", "unstarted":
		return ticket.TODO
	case "custom":
		return ticket.IN_PROGRESS
	case "done", "closed":
		return ticket.DONE
	default:
		return ticket.OTHER
	}
}

// statusMapper resolves a ClickUp task's domain status. When the scope config
// supplies explicit status name lists, a matching raw status name wins;
// otherwise the status.type-derived default is used.
type statusMapper struct {
	byName  map[string]string
	hasList bool
}

func newStatusMapper(sc *models.ClickUpScopeConfig) *statusMapper {
	m := &statusMapper{byName: map[string]string{}}
	if sc == nil {
		return m
	}
	add := func(names []string, domain string) {
		for _, n := range names {
			n = strings.ToLower(strings.TrimSpace(n))
			if n != "" {
				m.byName[n] = domain
				m.hasList = true
			}
		}
	}
	add(sc.IssueStatusTodo, ticket.TODO)
	add(sc.IssueStatusInProgress, ticket.IN_PROGRESS)
	add(sc.IssueStatusDone, ticket.DONE)
	return m
}

// statusOf returns the domain status for a raw status name / type. A user-configured
// name mapping takes precedence over the type-derived default.
func (m *statusMapper) statusOf(rawStatus, statusType string) string {
	if m.hasList {
		if domain, ok := m.byName[strings.ToLower(strings.TrimSpace(rawStatus))]; ok {
			return domain
		}
	}
	return statusFromType(statusType)
}

// issueTypeMatcher derives the domain ticket.Issue.Type from a task's derived
// type string using the scope config's regex patterns. Precedence is
// INCIDENT > BUG > REQUIREMENT; a task matching none defaults to REQUIREMENT.
//
// When no patterns are configured a sensible default (bug -> BUG) is applied so
// bug-typed tasks feed DORA change-failure-rate out of the box.
type issueTypeMatcher struct {
	incident    *regexp.Regexp
	bug         *regexp.Regexp
	requirement *regexp.Regexp
}

// defaultBugPattern matches ClickUp's built-in "Bug" custom task type name
// (case-insensitive) when the scope config leaves IssueTypeBug empty.
const defaultBugPattern = "(?i)^bug$"

func newIssueTypeMatcher(sc *models.ClickUpScopeConfig) (*issueTypeMatcher, errors.Error) {
	m := &issueTypeMatcher{}
	bugPattern := defaultBugPattern
	var incidentPattern, requirementPattern string
	if sc != nil {
		if sc.IssueTypeBug != "" {
			bugPattern = sc.IssueTypeBug
		}
		incidentPattern = sc.IssueTypeIncident
		requirementPattern = sc.IssueTypeRequirement
	}
	for _, p := range []struct {
		pattern string
		field   string
		out     **regexp.Regexp
	}{
		{incidentPattern, "issueTypeIncident", &m.incident},
		{bugPattern, "issueTypeBug", &m.bug},
		{requirementPattern, "issueTypeRequirement", &m.requirement},
	} {
		if p.pattern == "" {
			continue
		}
		re, err := errors.Convert01(regexp.Compile(p.pattern))
		if err != nil {
			return nil, errors.Default.Wrap(err, "invalid "+p.field+" pattern")
		}
		*p.out = re
	}
	return m, nil
}

// typeOf returns the domain issue type for a task's derived type string.
func (m *issueTypeMatcher) typeOf(taskType string) string {
	for _, c := range []struct {
		pattern *regexp.Regexp
		typ     string
	}{
		{m.incident, ticket.INCIDENT},
		{m.bug, ticket.BUG},
		{m.requirement, ticket.REQUIREMENT},
	} {
		if c.pattern != nil && c.pattern.MatchString(taskType) {
			return c.typ
		}
	}
	return ticket.REQUIREMENT
}
