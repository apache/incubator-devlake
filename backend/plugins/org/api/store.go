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
	"reflect"
	"sort"
	"strings"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type store interface {
	findAllUsers() ([]user, errors.Error)
	findAllTeams() ([]team, errors.Error)
	findTeamsPaginated(page, pageSize int, nameFilter string, grouped bool) ([]teamTree, int64, errors.Error)
	findAllAccounts() ([]account, errors.Error)
	findAllUserAccounts() ([]userAccount, errors.Error)
	findAllProjectMapping() ([]projectMapping, errors.Error)
	deleteAll(i interface{}) errors.Error
	save(items []interface{}) errors.Error
	findTeamById(id string) (*crossdomain.Team, errors.Error)
	createTeam(t *crossdomain.Team) errors.Error
	updateTeam(t *crossdomain.Team) errors.Error
	deleteTeam(id string) errors.Error
	findUsersPaginated(page, pageSize int, emailFilter string) ([]userWithTeams, int64, errors.Error)
	findUsersByIds(userIds []string) ([]crossdomain.User, errors.Error)
	findTeamsByIds(teamIds []string) ([]crossdomain.Team, errors.Error)
	findUserById(id string) (*crossdomain.User, errors.Error)
	createUser(u *crossdomain.User) errors.Error
	updateUser(u *crossdomain.User) errors.Error
	deleteUser(id string) errors.Error
	findTeamUsersByTeamId(teamId string) ([]crossdomain.TeamUser, errors.Error)
	findTeamUsersByUserId(userId string) ([]crossdomain.TeamUser, errors.Error)
	replaceTeamUsersForTeam(teamId string, userIds []string) errors.Error
	replaceTeamUsersForUser(userId string, teamIds []string) errors.Error
}

type dbStore struct {
	db     dal.Dal
	driver *helper.BatchSaveDivider
}

func NewDbStore(db dal.Dal, basicRes context.BasicRes) *dbStore {
	driver := helper.NewBatchSaveDivider(basicRes, 1000, "", "")
	return &dbStore{db: db, driver: driver}
}

func (d *dbStore) findAllUsers() ([]user, errors.Error) {
	var u *user
	var uu []crossdomain.User
	err := d.db.All(&uu)
	if err != nil {
		return nil, err
	}
	var tus []crossdomain.TeamUser
	err = d.db.All(&tus)
	if err != nil {
		return nil, err
	}
	return u.fromDomainLayer(uu, tus), nil
}

func (d *dbStore) findAllTeams() ([]team, errors.Error) {
	var tt []crossdomain.Team
	err := d.db.All(&tt)
	if err != nil {
		return nil, err
	}
	var t *team
	return t.fromDomainLayer(tt), nil
}

func (d *dbStore) findAllAccounts() ([]account, errors.Error) {
	var aa []crossdomain.Account
	err := d.db.All(&aa)
	if err != nil {
		return nil, err
	}
	var ua []crossdomain.UserAccount
	err = d.db.All(&ua)
	if err != nil {
		return nil, err
	}
	var a *account
	return a.fromDomainLayer(aa, ua), nil
}

func (d *dbStore) findAllUserAccounts() ([]userAccount, errors.Error) {
	var uas []crossdomain.UserAccount
	err := d.db.All(&uas)
	if err != nil {
		return nil, err
	}

	var au *userAccount
	return au.fromDomainLayer(uas), nil
}

func (d *dbStore) findAllProjectMapping() ([]projectMapping, errors.Error) {
	var mapping []crossdomain.ProjectMapping
	err := d.db.All(&mapping)
	if err != nil {
		return nil, err
	}
	var pm *projectMapping
	return pm.fromDomainLayer(mapping), nil
}
func (d *dbStore) deleteAll(i interface{}) errors.Error {
	return d.db.Delete(i, dal.Where("1=1"))
}

func (d *dbStore) save(items []interface{}) errors.Error {
	for _, item := range items {
		batch, err := d.driver.ForType(reflect.TypeOf(item))
		if err != nil {
			return err
		}
		err = batch.Add(item)
		if err != nil {
			return err
		}
	}
	d.driver.Close()
	return nil
}

func (d *dbStore) findTeamsPaginated(page, pageSize int, nameFilter string, grouped bool) ([]teamTree, int64, errors.Error) {
	if !grouped {
		clauses := []dal.Clause{dal.From(&crossdomain.Team{})}
		if nameFilter != "" {
			clauses = append(clauses, dal.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(nameFilter)+"%"))
		}

		count, err := d.db.Count(clauses...)
		if err != nil {
			return nil, 0, err
		}

		var domainTeams []crossdomain.Team
		err = d.db.All(&domainTeams,
			append(clauses,
				dal.Orderby("sorting_index ASC, created_at DESC, name ASC"),
				dal.Offset((page-1)*pageSize),
				dal.Limit(pageSize),
			)...,
		)
		if err != nil {
			return nil, 0, err
		}

		var t *team
		flatTeams := t.fromDomainLayer(domainTeams)
		teamIds := make([]string, 0, len(flatTeams))
		for _, item := range flatTeams {
			teamIds = append(teamIds, item.Id)
		}
		teamUserCounts, err := d.findTeamUserCounts(teamIds)
		if err != nil {
			return nil, 0, err
		}

		result := make([]teamTree, 0, len(flatTeams))
		for _, item := range flatTeams {
			result = append(result, teamTree{
				Id:           item.Id,
				Name:         item.Name,
				Alias:        item.Alias,
				ParentId:     item.ParentId,
				SortingIndex: item.SortingIndex,
				UserCount:    teamUserCounts[item.Id],
			})
		}
		return result, count, nil
	}

	parentClauses := []dal.Clause{
		dal.From(&crossdomain.Team{}),
		dal.Where("(parent_id IS NULL OR parent_id = '')"),
	}
	if nameFilter != "" {
		parentClauses = append(parentClauses, dal.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(nameFilter)+"%"))
	}

	count, err := d.db.Count(parentClauses...)
	if err != nil {
		return nil, 0, err
	}

	var parentDomainTeams []crossdomain.Team

	if nameFilter == "" {
		// Original behavior: in grouped mode without a name filter, only parents are
		// considered, and pagination is applied directly on them.
		parentClauses := []dal.Clause{
			dal.From(&crossdomain.Team{}),
			dal.Where("(parent_id IS NULL OR parent_id = '')"),
		}

		err = d.db.All(&parentDomainTeams,
			append(parentClauses,
				dal.Orderby("sorting_index ASC, name ASC"),
				dal.Offset((page-1)*pageSize),
				dal.Limit(pageSize),
			)...,
		)
		if err != nil {
			return nil, 0, err
		}
	} else {
		// When a name filter is provided in grouped mode, we want to return parents
		// whose own names match, as well as parents of children whose names match.
		lowerFilter := strings.ToLower(nameFilter)
		like := "%" + lowerFilter + "%"
		// 1. Parents whose own names match the filter.
		var filteredParents []crossdomain.Team
		err := d.db.All(&filteredParents,
			dal.From(&crossdomain.Team{}),
			dal.Where("(parent_id IS NULL OR parent_id = '')"),
			dal.Where("LOWER(name) LIKE ?", like),
		)
		if err != nil {
			return nil, 0, err
		}
		// 2. Children whose names match the filter.
		var filteredChildren []crossdomain.Team
		err = d.db.All(&filteredChildren,
			dal.From(&crossdomain.Team{}),
			dal.Where("parent_id IS NOT NULL AND parent_id <> ''"),
			dal.Where("LOWER(name) LIKE ?", like),
		)
		if err != nil {
			return nil, 0, err
		}
		// 3. Collect distinct parent IDs from both matching parents and parents of
		// matching children.
		parentIdSet := make(map[string]struct{})
		for _, p := range filteredParents {
			if p.Id != "" {
				parentIdSet[p.Id] = struct{}{}
			}
		}
		for _, c := range filteredChildren {
			if c.ParentId != "" {
				parentIdSet[c.ParentId] = struct{}{}
			}
		}
		if len(parentIdSet) == 0 {
			// No parents (or children) matched the filter.
			return []teamTree{}, 0, nil
		}
		parentIds := make([]string, 0, len(parentIdSet))
		for id := range parentIdSet {
			parentIds = append(parentIds, id)
		}
		// The total count is the number of distinct parent IDs.
		count = int64(len(parentIds))
		// 4. Load the paginated parents based on the collected IDs, preserving the
		// original ordering.
		err = d.db.All(
			&parentDomainTeams,
			dal.From(&crossdomain.Team{}),
			dal.Where("(parent_id IS NULL OR parent_id = '')"),
			dal.Where("id IN (?)", parentIds),
			dal.Orderby("sorting_index ASC, name ASC"),
			dal.Offset((page-1)*pageSize),
			dal.Limit(pageSize),
		)
	}
	if err != nil {
		return nil, 0, err
	}

	if len(parentDomainTeams) == 0 {
		return []teamTree{}, count, nil
	}

	parentIds := make([]string, 0, len(parentDomainTeams))
	for _, parent := range parentDomainTeams {
		parentIds = append(parentIds, parent.Id)
	}

	var childDomainTeams []crossdomain.Team
	err := d.db.All(
		&childDomainTeams,
		dal.From(&crossdomain.Team{}),
		dal.Where("parent_id IN ?", parentIds),
		dal.Orderby("sorting_index ASC, name ASC"),
	)
	if err != nil {
		return nil, 0, err
	}

	var t *team
	parentTeams := t.fromDomainLayer(parentDomainTeams)
	childTeams := t.fromDomainLayer(childDomainTeams)
	teamIds := make([]string, 0, len(parentTeams)+len(childTeams))
	for _, parent := range parentTeams {
		teamIds = append(teamIds, parent.Id)
	}
	for _, child := range childTeams {
		teamIds = append(teamIds, child.Id)
	}
	teamUserCounts, err := d.findTeamUserCounts(teamIds)
	if err != nil {
		return nil, 0, err
	}

	childrenByParent := make(map[string][]teamTree, len(parentTeams))
	for _, child := range childTeams {
		childrenByParent[child.ParentId] = append(childrenByParent[child.ParentId], teamTree{
			Id:           child.Id,
			Name:         child.Name,
			Alias:        child.Alias,
			ParentId:     child.ParentId,
			SortingIndex: child.SortingIndex,
			UserCount:    teamUserCounts[child.Id],
		})
	}

	result := make([]teamTree, 0, len(parentTeams))
	for _, parent := range parentTeams {
		result = append(result, teamTree{
			Id:           parent.Id,
			Name:         parent.Name,
			Alias:        parent.Alias,
			ParentId:     parent.ParentId,
			SortingIndex: parent.SortingIndex,
			UserCount:    teamUserCounts[parent.Id],
			Children:     childrenByParent[parent.Id],
		})
	}

	return result, count, nil
}

func (d *dbStore) findTeamUserCounts(teamIds []string) (map[string]int, errors.Error) {
	counts := make(map[string]int, len(teamIds))
	if len(teamIds) == 0 {
		return counts, nil
	}

	var teamUsers []crossdomain.TeamUser
	err := d.db.All(&teamUsers, dal.Where("team_id IN ?", teamIds))
	if err != nil {
		return nil, err
	}

	usersByTeam := make(map[string]map[string]struct{}, len(teamIds))
	for _, teamUser := range teamUsers {
		if usersByTeam[teamUser.TeamId] == nil {
			usersByTeam[teamUser.TeamId] = make(map[string]struct{})
		}
		usersByTeam[teamUser.TeamId][teamUser.UserId] = struct{}{}
	}

	for teamId, users := range usersByTeam {
		counts[teamId] = len(users)
	}

	return counts, nil
}

func (d *dbStore) findTeamById(id string) (*crossdomain.Team, errors.Error) {
	var t crossdomain.Team
	err := d.db.First(&t, dal.Where("id = ?", id))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (d *dbStore) createTeam(t *crossdomain.Team) errors.Error {
	return d.db.Create(t)
}

func (d *dbStore) updateTeam(t *crossdomain.Team) errors.Error {
	return d.db.Update(t)
}

func (d *dbStore) deleteTeam(id string) errors.Error {
	err := d.db.Delete(&crossdomain.TeamUser{}, dal.Where("team_id = ?", id))
	if err != nil {
		return err
	}
	return d.db.Delete(&crossdomain.Team{}, dal.Where("id = ?", id))
}

func (d *dbStore) findUsersPaginated(page, pageSize int, emailFilter string) ([]userWithTeams, int64, errors.Error) {
	clauses := []dal.Clause{dal.From(&crossdomain.User{})}
	if emailFilter != "" {
		clauses = append(clauses, dal.Where("LOWER(email) LIKE ?", "%"+strings.ToLower(emailFilter)+"%"))
	}
	count, err := d.db.Count(clauses...)
	if err != nil {
		return nil, 0, err
	}
	var uu []crossdomain.User
	err = d.db.All(&uu,
		append(clauses,
			dal.Orderby("name ASC"),
			dal.Offset((page-1)*pageSize),
			dal.Limit(pageSize),
		)...,
	)
	if err != nil {
		return nil, 0, err
	}
	// fetch team associations for the returned users
	userIds := make([]string, 0, len(uu))
	var tus []crossdomain.TeamUser
	if len(uu) > 0 {
		for _, u := range uu {
			userIds = append(userIds, u.Id)
		}
		err = d.db.All(&tus, dal.Where("user_id IN ?", userIds))
		if err != nil {
			return nil, 0, err
		}
	}

	// fetch account associations for the returned users
	var uas []crossdomain.UserAccount
	if len(userIds) > 0 {
		err = d.db.All(&uas, dal.Where("user_id IN ?", userIds))
		if err != nil {
			return nil, 0, err
		}
	}

	teamIdsSet := make(map[string]struct{}, len(tus))
	teamsByUser := make(map[string]map[string]struct{}, len(uu))
	for _, teamUser := range tus {
		if teamUser.UserId == "" || teamUser.TeamId == "" {
			continue
		}
		teamIdsSet[teamUser.TeamId] = struct{}{}
		if teamsByUser[teamUser.UserId] == nil {
			teamsByUser[teamUser.UserId] = make(map[string]struct{})
		}
		teamsByUser[teamUser.UserId][teamUser.TeamId] = struct{}{}
	}

	teamNameById := make(map[string]string, len(teamIdsSet))
	if len(teamIdsSet) > 0 {
		teamIds := make([]string, 0, len(teamIdsSet))
		for teamId := range teamIdsSet {
			teamIds = append(teamIds, teamId)
		}

		var teams []crossdomain.Team
		err = d.db.All(&teams, dal.Where("id IN ?", teamIds))
		if err != nil {
			return nil, 0, err
		}
		for _, team := range teams {
			teamNameById[team.Id] = team.Name
		}
	}

	accountsByUser := make(map[string]map[string]struct{}, len(uu))
	for _, userAccount := range uas {
		if userAccount.UserId == "" || userAccount.AccountId == "" {
			continue
		}
		if accountsByUser[userAccount.UserId] == nil {
			accountsByUser[userAccount.UserId] = make(map[string]struct{})
		}
		accountsByUser[userAccount.UserId][userAccount.AccountId] = struct{}{}
	}

	result := make([]userWithTeams, 0, len(uu))
	for _, domainUser := range uu {
		userTeamIdsSet := teamsByUser[domainUser.Id]
		userTeamIds := make([]string, 0, len(userTeamIdsSet))
		for teamId := range userTeamIdsSet {
			userTeamIds = append(userTeamIds, teamId)
		}
		sort.Strings(userTeamIds)

		userTeamNames := make([]string, 0, len(userTeamIds))
		for _, teamId := range userTeamIds {
			if teamName, exists := teamNameById[teamId]; exists && teamName != "" {
				userTeamNames = append(userTeamNames, teamName)
			}
		}
		sort.Strings(userTeamNames)

		userAccountIdsSet := accountsByUser[domainUser.Id]
		userAccountIds := make([]string, 0, len(userAccountIdsSet))
		for accountId := range userAccountIdsSet {
			userAccountIds = append(userAccountIds, accountId)
		}
		sort.Strings(userAccountIds)
		accountSources := extractAccountSources(userAccountIds)

		result = append(result, userWithTeams{
			Id:             domainUser.Id,
			Name:           domainUser.Name,
			Email:          domainUser.Email,
			TeamIds:        strings.Join(userTeamIds, ";"),
			TeamCount:      len(userTeamIds),
			TeamNames:      userTeamNames,
			AccountCount:   len(userAccountIds),
			AccountSources: accountSources,
		})
	}

	return result, count, nil
}

func extractAccountSources(accountIds []string) []string {
	if len(accountIds) == 0 {
		return []string{}
	}

	sourceSet := make(map[string]struct{}, len(accountIds))
	for _, accountId := range accountIds {
		source := accountSourceFromAccountId(accountId)
		if source == "" {
			continue
		}
		sourceSet[source] = struct{}{}
	}

	sources := make([]string, 0, len(sourceSet))
	for source := range sourceSet {
		sources = append(sources, source)
	}
	sort.Strings(sources)
	return sources
}

func accountSourceFromAccountId(accountId string) string {
	if accountId == "" {
		return ""
	}

	separatorIndex := strings.Index(accountId, ":")
	if separatorIndex <= 0 {
		return "Unknown"
	}

	pluginName := strings.ToLower(strings.TrimSpace(accountId[:separatorIndex]))
	switch pluginName {
	case "github":
		return "GitHub"
	case "gitlab":
		return "GitLab"
	case "jira":
		return "Jira"
	case "azuredevops", "azuredevops_go":
		return "Azure DevOps"
	case "bitbucket":
		return "Bitbucket"
	case "bitbucket_server":
		return "Bitbucket Server"
	case "gh-copilot":
		return "GitHub Copilot"
	case "sonarqube":
		return "SonarQube"
	case "tapd":
		return "TAPD"
	}

	parts := strings.FieldsFunc(pluginName, func(r rune) bool {
		return r == '_' || r == '-'
	})
	if len(parts) == 0 {
		return "Unknown"
	}

	titled := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		titled = append(titled, titleCaseWord(part))
	}
	if len(titled) == 0 {
		return "Unknown"
	}

	return strings.Join(titled, " ")
}

func titleCaseWord(word string) string {
	if word == "" {
		return ""
	}
	if len(word) == 1 {
		return strings.ToUpper(word)
	}
	return strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
}

func (d *dbStore) findUsersByIds(userIds []string) ([]crossdomain.User, errors.Error) {
	if len(userIds) == 0 {
		return []crossdomain.User{}, nil
	}

	var users []crossdomain.User
	err := d.db.All(&users, dal.Where("id IN ?", userIds))
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (d *dbStore) findTeamsByIds(teamIds []string) ([]crossdomain.Team, errors.Error) {
	if len(teamIds) == 0 {
		return []crossdomain.Team{}, nil
	}

	var teams []crossdomain.Team
	err := d.db.All(&teams, dal.Where("id IN ?", teamIds))
	if err != nil {
		return nil, err
	}

	return teams, nil
}

func (d *dbStore) findUserById(id string) (*crossdomain.User, errors.Error) {
	var u crossdomain.User
	err := d.db.First(&u, dal.Where("id = ?", id))
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *dbStore) createUser(u *crossdomain.User) errors.Error {
	return d.db.Create(u)
}

func (d *dbStore) updateUser(u *crossdomain.User) errors.Error {
	return d.db.Update(u)
}

func (d *dbStore) deleteUser(id string) errors.Error {
	err := d.db.Delete(&crossdomain.TeamUser{}, dal.Where("user_id = ?", id))
	if err != nil {
		return err
	}
	err = d.db.Delete(&crossdomain.UserAccount{}, dal.Where("user_id = ?", id))
	if err != nil {
		return err
	}
	return d.db.Delete(&crossdomain.User{}, dal.Where("id = ?", id))
}

func (d *dbStore) findTeamUsersByUserId(userId string) ([]crossdomain.TeamUser, errors.Error) {
	var tus []crossdomain.TeamUser
	err := d.db.All(&tus, dal.Where("user_id = ?", userId))
	if err != nil {
		return nil, err
	}
	return tus, nil
}

func (d *dbStore) findTeamUsersByTeamId(teamId string) ([]crossdomain.TeamUser, errors.Error) {
	var tus []crossdomain.TeamUser
	err := d.db.All(&tus, dal.Where("team_id = ?", teamId))
	if err != nil {
		return nil, err
	}
	return tus, nil
}

func (d *dbStore) replaceTeamUsersForTeam(teamId string, userIds []string) errors.Error {
	err := d.db.Delete(&crossdomain.TeamUser{}, dal.Where("team_id = ?", teamId))
	if err != nil {
		return err
	}

	seen := make(map[string]struct{}, len(userIds))
	for _, userId := range userIds {
		if userId == "" {
			continue
		}
		if _, exists := seen[userId]; exists {
			continue
		}
		seen[userId] = struct{}{}

		err = d.db.Create(&crossdomain.TeamUser{
			TeamId: teamId,
			UserId: userId,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *dbStore) replaceTeamUsersForUser(userId string, teamIds []string) errors.Error {
	err := d.db.Delete(&crossdomain.TeamUser{}, dal.Where("user_id = ?", userId))
	if err != nil {
		return err
	}
	for _, teamId := range teamIds {
		if teamId == "" {
			continue
		}
		err = d.db.Create(&crossdomain.TeamUser{
			TeamId: teamId,
			UserId: userId,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
