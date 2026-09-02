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
	"fmt"
	"strconv"
	"strings"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/plugin"
	"github.com/apache/devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/devlake/helpers/pluginhelper/api/models"
	"github.com/apache/devlake/plugins/kiro/models"
	"github.com/apache/devlake/plugins/kiro/tasks"
)

// listKiroRemoteScopes browses the export layout as a tree.
//
// Three levels, mirroring Kiro's own S3 partitioning:
//
//	(root)        -> one group per AWS account with exported data
//	{account}     -> one group per year, plus a whole-year scope
//	{account}/{y} -> one selectable scope per month
//
// Everything comes from S3 rather than user input. That is the point: a
// hand-typed prefix cannot be validated from the outcome, because a typo and a
// month with no data both produce a successful run that collects nothing.
func listKiroRemoteScopes(connection *models.KiroConnection, groupId string) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice],
	err errors.Error,
) {
	if connection == nil {
		return nil, errors.BadInput.New("connection is required")
	}

	discovery, err := tasks.NewDiscovery(connection)
	if err != nil {
		return nil, err
	}

	accountId, year, err := parseGroupId(groupId)
	if err != nil {
		return nil, err
	}

	switch {
	case accountId == "":
		return listAccountGroups(discovery)
	case year == 0:
		return listYearGroups(discovery, accountId)
	default:
		return listMonthScopes(discovery, accountId, year)
	}
}

// listAccountGroups is the tree root: the accounts that actually have exports.
func listAccountGroups(discovery *tasks.Discovery) (
	[]dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice], errors.Error,
) {
	accounts, err := discovery.ListAccounts()
	if err != nil {
		return nil, err
	}

	entries := make([]dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice], 0, len(accounts))
	for _, accountId := range accounts {
		entries = append(entries, dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice]{
			Type:     api.RAS_ENTRY_TYPE_GROUP,
			Id:       accountId,
			Name:     accountId,
			FullName: accountId,
		})
	}
	return entries, nil
}

// listYearGroups lists the years under an account.
//
// Each year is offered both as a group to expand and as a directly selectable
// scope, because a nil month means "collect the whole year" - which is how a
// year-long backfill is expressed without creating twelve scopes by hand.
func listYearGroups(discovery *tasks.Discovery, accountId string) (
	[]dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice], errors.Error,
) {
	years, err := discovery.ListYears(accountId)
	if err != nil {
		return nil, err
	}

	entries := make([]dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice], 0, len(years)*2)
	for _, year := range years {
		groupId := fmt.Sprintf("%s/%04d", accountId, year)
		parent := accountId

		entries = append(entries, dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice]{
			Type:     api.RAS_ENTRY_TYPE_GROUP,
			ParentId: &parent,
			Id:       groupId,
			Name:     fmt.Sprintf("%04d", year),
			FullName: groupId,
		})

		wholeYear := &models.KiroS3Slice{AccountId: accountId, Year: year}
		*wholeYear = wholeYear.Sanitize()
		entries = append(entries, dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			ParentId: &parent,
			Id:       wholeYear.Id,
			Name:     fmt.Sprintf("%04d (whole year)", year),
			FullName: wholeYear.ScopeName(),
			Data:     wholeYear,
		})
	}
	return entries, nil
}

// listMonthScopes lists the months that hold data for an account and year.
func listMonthScopes(discovery *tasks.Discovery, accountId string, year int) (
	[]dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice], errors.Error,
) {
	months, err := discovery.ListMonths(accountId, year)
	if err != nil {
		return nil, err
	}

	parent := fmt.Sprintf("%s/%04d", accountId, year)
	entries := make([]dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice], 0, len(months))
	for _, month := range months {
		m := month
		slice := &models.KiroS3Slice{AccountId: accountId, Year: year, Month: &m}
		*slice = slice.Sanitize()

		entries = append(entries, dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			ParentId: &parent,
			Id:       slice.Id,
			Name:     fmt.Sprintf("%04d-%02d", year, month),
			FullName: slice.ScopeName(),
			Data:     slice,
		})
	}
	return entries, nil
}

// searchKiroRemoteScopes filters the discovered months by substring.
//
// Matching is against "{account} {year}-{month}", so "2026-07" or an account
// number both work. The search space is one listing per year, small enough to
// scan without an index.
func searchKiroRemoteScopes(
	connection *models.KiroConnection,
	query string,
	page int,
	pageSize int,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice],
	err errors.Error,
) {
	empty := []dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice]{}
	if connection == nil {
		return empty, nil
	}
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return empty, nil
	}

	discovery, err := tasks.NewDiscovery(connection)
	if err != nil {
		return nil, err
	}
	accounts, err := discovery.ListAccounts()
	if err != nil {
		return nil, err
	}

	matches := make([]dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice], 0)
	for _, accountId := range accounts {
		years, yearErr := discovery.ListYears(accountId)
		if yearErr != nil {
			return nil, yearErr
		}
		for _, year := range years {
			months, monthErr := discovery.ListMonths(accountId, year)
			if monthErr != nil {
				return nil, monthErr
			}
			for _, month := range months {
				m := month
				slice := &models.KiroS3Slice{AccountId: accountId, Year: year, Month: &m}
				*slice = slice.Sanitize()

				if !strings.Contains(strings.ToLower(slice.ScopeName()), query) &&
					!strings.Contains(strings.ToLower(slice.Id), query) {
					continue
				}
				matches = append(matches, dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice]{
					Type:     api.RAS_ENTRY_TYPE_SCOPE,
					Id:       slice.Id,
					Name:     slice.ScopeName(),
					FullName: slice.ScopeName(),
					Data:     slice,
				})
			}
		}
	}

	return paginate(matches, page, pageSize), nil
}

// paginate applies the requested page window to an in-memory result set.
func paginate(
	entries []dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice],
	page int, pageSize int,
) []dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice] {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	start := (page - 1) * pageSize
	if start >= len(entries) {
		return []dsmodels.DsRemoteApiScopeListEntry[models.KiroS3Slice]{}
	}
	end := start + pageSize
	if end > len(entries) {
		end = len(entries)
	}
	return entries[start:end]
}

// parseGroupId splits a tree node id into its parts.
//
// "" is the root, "{account}" is an account node, "{account}/{year}" is a year
// node.
func parseGroupId(groupId string) (accountId string, year int, err errors.Error) {
	trimmed := strings.Trim(strings.TrimSpace(groupId), "/")
	if trimmed == "" {
		return "", 0, nil
	}

	parts := strings.Split(trimmed, "/")
	switch len(parts) {
	case 1:
		return parts[0], 0, nil
	case 2:
		parsedYear, convErr := strconv.Atoi(parts[1])
		if convErr != nil {
			return "", 0, errors.BadInput.New("invalid year in groupId: " + groupId)
		}
		return parts[0], parsedYear, nil
	default:
		return "", 0, errors.BadInput.New("unrecognized groupId: " + groupId)
	}
}

// RemoteScopes browses the Kiro export layout in S3.
//
// Implemented directly rather than through the shared scope-list helper. That
// helper builds an HTTP client from the connection first, and its constructor
// runs a DNS check on the endpoint - which fails here, because a bucket name is
// not a hostname. The helper is built for HTTP data sources; this one is S3.
// @Summary list available kiro scopes discovered from S3
// @Description Browse accounts, years and months that actually have exported data
// @Tags plugins/kiro
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param groupId query string false "account id, or account/year"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.KiroS3Slice]
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.KiroConnection{}
	if err := connectionHelper.First(connection, input.Params); err != nil {
		return nil, err
	}

	children, err := listKiroRemoteScopes(connection, input.Query.Get("groupId"))
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body: dsmodels.DsRemoteApiScopeList[models.KiroS3Slice]{Children: children},
	}, nil
}

// SearchRemoteScopes finds discovered scopes by substring.
//
// Implemented directly rather than through the shared search helper: that
// helper's callback receives only an HTTP ApiClient, and discovery here needs
// the connection itself to build an S3 client.
// @Summary search kiro scopes discovered from S3
// @Description Search the discovered months by account or year-month
// @Tags plugins/kiro
// @Accept application/json
// @Param connectionId path int true "connection ID"
// @Param search query string false "search"
// @Param page query int false "page number"
// @Param pageSize query int false "page size per page"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.KiroS3Slice] "the parentIds are always null"
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/kiro/connections/{connectionId}/search-remote-scopes [GET]
func SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.KiroConnection{}
	if err := connectionHelper.First(connection, input.Params); err != nil {
		return nil, err
	}

	page, _ := strconv.Atoi(input.Query.Get("page"))
	pageSize, _ := strconv.Atoi(input.Query.Get("pageSize"))

	children, err := searchKiroRemoteScopes(connection, input.Query.Get("search"), page, pageSize)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body: dsmodels.DsRemoteApiScopeList[models.KiroS3Slice]{Children: children},
	}, nil
}
