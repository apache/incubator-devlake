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

	"github.com/apache/devlake/core/dal"
	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
)

// ScopeDuplicateConnection is a connection that shares a repository scope.
type ScopeDuplicateConnection struct {
	ConnectionId   uint64 `json:"connectionId"`
	ConnectionName string `json:"connectionName"`
}

// ScopeDuplicateGroup is one repository that appears under multiple connections
// (diagnostics) or already exists under another connection (pre-add check).
// Unlike GitHub/GitLab, a Bitbucket repo's natural id (BitbucketId) is the
// "owner/repo" full name string rather than a numeric id.
type ScopeDuplicateGroup struct {
	BitbucketId string                     `json:"bitbucketId"`
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
	BitbucketId    string `gorm:"column:bitbucket_id"`
	HTMLUrl        string `gorm:"column:html_url"`
	ConnectionId   uint64 `gorm:"column:connection_id"`
	ConnectionName string `gorm:"column:connection_name"`
}

// GetScopeDuplicates returns Bitbucket repositories registered under more than
// one connection, or (with connectionId + bitbucketIds) candidates already
// present on other connections.
// @Summary Find Bitbucket scopes duplicated across connections
// @Description Diagnostics: groups where the same bitbucketId (owner/repo) appears on more than one connection.
// @Description Pre-add check: pass connectionId and bitbucketIds to find candidates already registered elsewhere.
// @Tags plugins/bitbucket
// @Param connectionId query int false "Current connection id (pre-add check)"
// @Param bitbucketIds query string false "Comma-separated Bitbucket repo full names (owner/repo) to check (pre-add check)"
// @Success 200 {object} ScopeDuplicatesOutput
// @Failure 400 {object} shared.ApiBody "Bad Request"
// @Failure 500 {object} shared.ApiBody "Internal Error"
// @Router /plugins/bitbucket/scope-duplicates [GET]
func GetScopeDuplicates(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, bitbucketIds, err := parseScopeDuplicateQuery(input)
	if err != nil {
		return nil, err
	}

	// Pre-add check with an empty selection: nothing to warn about.
	if connectionId != nil && len(bitbucketIds) == 0 {
		return &plugin.ApiResourceOutput{
			Body:   ScopeDuplicatesOutput{Duplicates: []ScopeDuplicateGroup{}},
			Status: http.StatusOK,
		}, nil
	}

	rows, err := queryScopeDuplicateRows(basicRes.GetDal(), connectionId, bitbucketIds)
	if err != nil {
		return nil, err
	}

	return &plugin.ApiResourceOutput{
		Body:   ScopeDuplicatesOutput{Duplicates: groupScopeDuplicateRows(rows)},
		Status: http.StatusOK,
	}, nil
}

func parseScopeDuplicateQuery(input *plugin.ApiResourceInput) (*uint64, []string, errors.Error) {
	var connectionId *uint64
	if v := input.Query.Get("connectionId"); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, nil, errors.BadInput.Wrap(err, "invalid connectionId")
		}
		connectionId = &id
	}

	var bitbucketIds []string
	if v := input.Query.Get("bitbucketIds"); v != "" {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			bitbucketIds = append(bitbucketIds, part)
		}
	}

	if len(bitbucketIds) > 0 && connectionId == nil {
		return nil, nil, errors.BadInput.New("connectionId is required when bitbucketIds is provided")
	}

	return connectionId, bitbucketIds, nil
}

// queryScopeDuplicateRows loads only the rows needed for the requested mode.
// Check mode: selected bitbucketIds on any connection other than connectionId.
// Diagnostics: bitbucketIds that already appear on more than one connection.
func queryScopeDuplicateRows(db dal.Dal, connectionId *uint64, bitbucketIds []string) ([]scopeDuplicateRow, errors.Error) {
	clauses := []dal.Clause{
		dal.Select("r.bitbucket_id, r.html_url, r.connection_id, c.name AS connection_name"),
		dal.From("_tool_bitbucket_repos r"),
		dal.Join("INNER JOIN _tool_bitbucket_connections c ON c.id = r.connection_id"),
		dal.Orderby("r.bitbucket_id ASC, r.connection_id ASC"),
	}

	if connectionId != nil {
		clauses = append(clauses, dal.Where(
			"r.bitbucket_id IN ? AND r.connection_id != ?",
			bitbucketIds,
			*connectionId,
		))
	} else {
		clauses = append(clauses, dal.Where(`r.bitbucket_id IN (
			SELECT bitbucket_id FROM _tool_bitbucket_repos
			GROUP BY bitbucket_id
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
		if current == nil || current.BitbucketId != row.BitbucketId {
			flush()
			current = &ScopeDuplicateGroup{
				BitbucketId: row.BitbucketId,
				HTMLUrl:     row.HTMLUrl,
				// Bitbucket's "full name" is the owner/repo id itself, see BitbucketRepo.ScopeFullName().
				FullName:    row.BitbucketId,
				Connections: make([]ScopeDuplicateConnection, 0, 2),
			}
			seenConns = make(map[uint64]struct{})
		}
		if current.HTMLUrl == "" && row.HTMLUrl != "" {
			current.HTMLUrl = row.HTMLUrl
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
