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

package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

// ScopeDuplicateConnection is a connection that shares a project scope.
type ScopeDuplicateConnection struct {
	ConnectionId   uint64 `json:"connectionId"`
	ConnectionName string `json:"connectionName"`
}

// ScopeDuplicateGroup is one project that appears under multiple connections
// (diagnostics) or already exists under another connection (pre-add check).
type ScopeDuplicateGroup struct {
	GitlabId    int                        `json:"gitlabId"`
	HTMLUrl     string                     `json:"htmlUrl"`
	FullName    string                     `json:"fullName"`
	Connections []ScopeDuplicateConnection `json:"connections"`
}

// ScopeDuplicatesOutput is the response body for GetScopeDuplicates.
type ScopeDuplicatesOutput struct {
	Duplicates []ScopeDuplicateGroup `json:"duplicates"`
}

// scopeDuplicateRow is one joined row from the scoped SQL query.
type scopeDuplicateRow struct {
	GitlabId       int    `gorm:"column:gitlab_id"`
	HTMLUrl        string `gorm:"column:html_url"`
	FullName       string `gorm:"column:full_name"`
	ConnectionId   uint64 `gorm:"column:connection_id"`
	ConnectionName string `gorm:"column:connection_name"`
}

// GetScopeDuplicates returns GitLab projects registered under more than one
// connection, or (with connectionId + gitlabIds) candidates already present on
// other connections.
// @Summary Find GitLab scopes duplicated across connections
// @Description Diagnostics: groups where the same gitlabId appears on more than one connection.
// @Description Pre-add check: pass connectionId and gitlabIds to find candidates already registered elsewhere.
// @Tags plugins/gitlab
// @Param connectionId query int false "Current connection id (pre-add check)"
// @Param gitlabIds query string false "Comma-separated GitLab project ids to check (pre-add check)"
// @Success 200 {object} ScopeDuplicatesOutput
// @Failure 400 {object} shared.ApiBody "Bad Request"
// @Failure 500 {object} shared.ApiBody "Internal Error"
// @Router /plugins/gitlab/scope-duplicates [GET]
func GetScopeDuplicates(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, gitlabIds, err := parseScopeDuplicateQuery(input)
	if err != nil {
		return nil, err
	}

	// Pre-add check with an empty selection: nothing to warn about.
	if connectionId != nil && len(gitlabIds) == 0 {
		return &plugin.ApiResourceOutput{
			Body:   ScopeDuplicatesOutput{Duplicates: []ScopeDuplicateGroup{}},
			Status: http.StatusOK,
		}, nil
	}

	rows, err := queryScopeDuplicateRows(basicRes.GetDal(), connectionId, gitlabIds)
	if err != nil {
		return nil, err
	}

	return &plugin.ApiResourceOutput{
		Body:   ScopeDuplicatesOutput{Duplicates: groupScopeDuplicateRows(rows)},
		Status: http.StatusOK,
	}, nil
}

func parseScopeDuplicateQuery(input *plugin.ApiResourceInput) (*uint64, []int, errors.Error) {
	var connectionId *uint64
	if v := input.Query.Get("connectionId"); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, nil, errors.BadInput.Wrap(err, "invalid connectionId")
		}
		connectionId = &id
	}

	var gitlabIds []int
	if v := input.Query.Get("gitlabIds"); v != "" {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			id, err := strconv.Atoi(part)
			if err != nil {
				return nil, nil, errors.BadInput.Wrap(err, "invalid gitlabIds")
			}
			gitlabIds = append(gitlabIds, id)
		}
	}

	if len(gitlabIds) > 0 && connectionId == nil {
		return nil, nil, errors.BadInput.New("connectionId is required when gitlabIds is provided")
	}

	return connectionId, gitlabIds, nil
}

// queryScopeDuplicateRows loads only the rows needed for the requested mode.
// Check mode: selected gitlabIds on any connection other than connectionId.
// Diagnostics: gitlabIds that already appear on more than one connection.
func queryScopeDuplicateRows(db dal.Dal, connectionId *uint64, gitlabIds []int) ([]scopeDuplicateRow, errors.Error) {
	clauses := []dal.Clause{
		dal.Select("r.gitlab_id, r.web_url AS html_url, r.path_with_namespace AS full_name, r.connection_id, c.name AS connection_name"),
		dal.From("_tool_gitlab_projects r"),
		dal.Join("INNER JOIN _tool_gitlab_connections c ON c.id = r.connection_id"),
		dal.Orderby("r.gitlab_id ASC, r.connection_id ASC"),
	}

	if connectionId != nil {
		clauses = append(clauses, dal.Where(
			"r.gitlab_id IN ? AND r.connection_id != ?",
			gitlabIds,
			*connectionId,
		))
	} else {
		clauses = append(clauses, dal.Where(`r.gitlab_id IN (
			SELECT gitlab_id FROM _tool_gitlab_projects
			GROUP BY gitlab_id
			HAVING COUNT(DISTINCT connection_id) > 1
		)`))
	}

	var rows []scopeDuplicateRow
	if err := db.All(&rows, clauses...); err != nil {
		return nil, err
	}
	return rows, nil
}

// groupScopeDuplicateRows collapses already-filtered SQL rows into API groups.
func groupScopeDuplicateRows(rows []scopeDuplicateRow) []ScopeDuplicateGroup {
	if len(rows) == 0 {
		return []ScopeDuplicateGroup{}
	}

	result := make([]ScopeDuplicateGroup, 0)
	var current *ScopeDuplicateGroup
	seenConns := make(map[uint64]struct{})

	flush := func() {
		if current != nil {
			result = append(result, *current)
		}
	}

	for _, row := range rows {
		if current == nil || current.GitlabId != row.GitlabId {
			flush()
			current = &ScopeDuplicateGroup{
				GitlabId:    row.GitlabId,
				HTMLUrl:     row.HTMLUrl,
				FullName:    row.FullName,
				Connections: make([]ScopeDuplicateConnection, 0, 2),
			}
			seenConns = make(map[uint64]struct{})
		}
		if current.HTMLUrl == "" && row.HTMLUrl != "" {
			current.HTMLUrl = row.HTMLUrl
		}
		if current.FullName == "" && row.FullName != "" {
			current.FullName = row.FullName
		}
		if _, ok := seenConns[row.ConnectionId]; ok {
			continue
		}
		seenConns[row.ConnectionId] = struct{}{}
		current.Connections = append(current.Connections, ScopeDuplicateConnection{
			ConnectionId:   row.ConnectionId,
			ConnectionName: row.ConnectionName,
		})
	}
	flush()
	return result
}
