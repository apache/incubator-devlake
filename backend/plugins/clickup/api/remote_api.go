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
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

// ClickUpRemotePagination is a placeholder: the ClickUp v2 hierarchy endpoints
// used here (team/space/folder/list) are not paginated, so there is never a
// next page. It exists to satisfy the DsRemoteApiScopeListHelper generic.
type ClickUpRemotePagination struct{}

// groupId prefixes encode which level of the ClickUp hierarchy a group entry
// points at, so a single list function can walk Team -> Space -> Folder -> List.
const (
	groupTeamPrefix   = "team:"
	groupSpacePrefix  = "space:"
	groupFolderPrefix = "folder:"
)

type clickUpNamedEntities struct {
	Teams   []clickUpNamedEntity `json:"teams"`
	Spaces  []clickUpNamedEntity `json:"spaces"`
	Folders []clickUpNamedEntity `json:"folders"`
	Lists   []clickUpNamedEntity `json:"lists"`
}

type clickUpNamedEntity struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// listClickUpRemoteScopes walks the ClickUp hierarchy one level at a time. The
// config UI drives it via `groupId`:
//
//	""             -> workspaces (Team) as groups
//	team:{id}      -> spaces as groups
//	space:{id}     -> folders as groups + folderless lists as selectable scopes
//	folder:{id}    -> lists as selectable scopes
func listClickUpRemoteScopes(
	_ *models.ClickUpConnection,
	apiClient plugin.ApiClient,
	groupId string,
	_ ClickUpRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.ClickUpList],
	nextPage *ClickUpRemotePagination,
	err errors.Error,
) {
	switch {
	case groupId == "":
		return listTeamsAsGroups(apiClient)
	case strings.HasPrefix(groupId, groupTeamPrefix):
		return listSpacesAsGroups(apiClient, strings.TrimPrefix(groupId, groupTeamPrefix))
	case strings.HasPrefix(groupId, groupSpacePrefix):
		return listSpaceChildren(apiClient, strings.TrimPrefix(groupId, groupSpacePrefix))
	case strings.HasPrefix(groupId, groupFolderPrefix):
		return listFolderLists(apiClient, strings.TrimPrefix(groupId, groupFolderPrefix))
	default:
		return nil, nil, errors.BadInput.New(fmt.Sprintf("unrecognized groupId %q", groupId))
	}
}

func getEntities(apiClient plugin.ApiClient, path string) (*clickUpNamedEntities, errors.Error) {
	res, err := apiClient.Get(path, nil, nil)
	if err != nil {
		return nil, errors.Default.Wrap(err, "failed to query ClickUp "+path)
	}
	var body clickUpNamedEntities
	if err := api.UnmarshalResponse(res, &body); err != nil {
		return nil, errors.Default.Wrap(err, "failed to unmarshal ClickUp "+path+" response")
	}
	return &body, nil
}

func groupEntry(id, name string) dsmodels.DsRemoteApiScopeListEntry[models.ClickUpList] {
	return dsmodels.DsRemoteApiScopeListEntry[models.ClickUpList]{
		Type:     api.RAS_ENTRY_TYPE_GROUP,
		ParentId: nil,
		Id:       id,
		Name:     name,
		FullName: name,
	}
}

func listTeamsAsGroups(apiClient plugin.ApiClient) ([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpList], *ClickUpRemotePagination, errors.Error) {
	body, err := getEntities(apiClient, "team")
	if err != nil {
		return nil, nil, err
	}
	children := make([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpList], 0, len(body.Teams))
	for _, team := range body.Teams {
		children = append(children, groupEntry(groupTeamPrefix+team.Id, team.Name))
	}
	return children, nil, nil
}

func listSpacesAsGroups(apiClient plugin.ApiClient, teamId string) ([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpList], *ClickUpRemotePagination, errors.Error) {
	body, err := getEntities(apiClient, fmt.Sprintf("team/%s/space", teamId))
	if err != nil {
		return nil, nil, err
	}
	children := make([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpList], 0, len(body.Spaces))
	for _, space := range body.Spaces {
		children = append(children, groupEntry(groupSpacePrefix+space.Id, space.Name))
	}
	return children, nil, nil
}

func listSpaceChildren(apiClient plugin.ApiClient, spaceId string) ([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpList], *ClickUpRemotePagination, errors.Error) {
	folders, err := getEntities(apiClient, fmt.Sprintf("space/%s/folder", spaceId))
	if err != nil {
		return nil, nil, err
	}
	children := make([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpList], 0)
	for _, folder := range folders.Folders {
		children = append(children, groupEntry(groupFolderPrefix+folder.Id, folder.Name))
	}
	// Folderless lists live directly under the space.
	lists, err := getEntities(apiClient, fmt.Sprintf("space/%s/list", spaceId))
	if err != nil {
		return nil, nil, err
	}
	children = append(children, listsToScopeEntries(lists.Lists, spaceId)...)
	return children, nil, nil
}

func listFolderLists(apiClient plugin.ApiClient, folderId string) ([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpList], *ClickUpRemotePagination, errors.Error) {
	body, err := getEntities(apiClient, fmt.Sprintf("folder/%s/list", folderId))
	if err != nil {
		return nil, nil, err
	}
	return listsToScopeEntries(body.Lists, ""), nil, nil
}

// listsToScopeEntries maps ClickUp lists into selectable (leaf) scope entries.
func listsToScopeEntries(lists []clickUpNamedEntity, spaceId string) []dsmodels.DsRemoteApiScopeListEntry[models.ClickUpList] {
	entries := make([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpList], 0, len(lists))
	for _, list := range lists {
		list := list
		entries = append(entries, dsmodels.DsRemoteApiScopeListEntry[models.ClickUpList]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			ParentId: nil,
			Id:       list.Id,
			Name:     list.Name,
			FullName: list.Name,
			Data: &models.ClickUpList{
				ListId:  list.Id,
				Name:    list.Name,
				SpaceId: spaceId,
			},
		})
	}
	return entries
}

// RemoteScopes lists the ClickUp lists available on the connection so the
// config UI can enumerate selectable scopes.
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// Proxy forwards arbitrary requests to the ClickUp API through the connection.
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
