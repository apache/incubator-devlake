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
// used here (team/space/folder) are not paginated, so there is never a next
// page. It exists to satisfy the DsRemoteApiScopeListHelper generic.
type ClickUpRemotePagination struct{}

// groupId prefixes encode which level of the ClickUp hierarchy a group entry
// points at, so a single list function can walk Team -> Space -> Folder. The
// selectable scope is the Folder (= board); spaces and teams are navigation
// groups only.
const (
	groupTeamPrefix  = "team:"
	groupSpacePrefix = "space:"
)

type clickUpNamedEntities struct {
	Teams   []clickUpNamedEntity `json:"teams"`
	Spaces  []clickUpNamedEntity `json:"spaces"`
	Folders []clickUpNamedEntity `json:"folders"`
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
//	space:{id}     -> folders as selectable scopes (boards)
func listClickUpRemoteScopes(
	_ *models.ClickUpConnection,
	apiClient plugin.ApiClient,
	groupId string,
	_ ClickUpRemotePagination,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.ClickUpFolder],
	nextPage *ClickUpRemotePagination,
	err errors.Error,
) {
	switch {
	case groupId == "":
		return listTeamsAsGroups(apiClient)
	case strings.HasPrefix(groupId, groupTeamPrefix):
		return listSpacesAsGroups(apiClient, strings.TrimPrefix(groupId, groupTeamPrefix), groupId)
	case strings.HasPrefix(groupId, groupSpacePrefix):
		return listFoldersAsScopes(apiClient, strings.TrimPrefix(groupId, groupSpacePrefix), groupId)
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

// groupEntry builds a navigation-group row. parentId is the id of the group
// this row is nested under (nil for the top level); the miller-column UI uses
// it to render children in the next column instead of inline.
func groupEntry(id, name string, parentId *string) dsmodels.DsRemoteApiScopeListEntry[models.ClickUpFolder] {
	return dsmodels.DsRemoteApiScopeListEntry[models.ClickUpFolder]{
		Type:     api.RAS_ENTRY_TYPE_GROUP,
		ParentId: parentId,
		Id:       id,
		Name:     name,
		FullName: name,
	}
}

func listTeamsAsGroups(apiClient plugin.ApiClient) ([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpFolder], *ClickUpRemotePagination, errors.Error) {
	body, err := getEntities(apiClient, "team")
	if err != nil {
		return nil, nil, err
	}
	children := make([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpFolder], 0, len(body.Teams))
	for _, team := range body.Teams {
		children = append(children, groupEntry(groupTeamPrefix+team.Id, team.Name, nil))
	}
	return children, nil, nil
}

func listSpacesAsGroups(apiClient plugin.ApiClient, teamId, parentId string) ([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpFolder], *ClickUpRemotePagination, errors.Error) {
	body, err := getEntities(apiClient, fmt.Sprintf("team/%s/space", teamId))
	if err != nil {
		return nil, nil, err
	}
	children := make([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpFolder], 0, len(body.Spaces))
	for _, space := range body.Spaces {
		children = append(children, groupEntry(groupSpacePrefix+space.Id, space.Name, &parentId))
	}
	return children, nil, nil
}

// listFoldersAsScopes returns a space's folders as selectable (leaf) scope
// entries — the folder is the board a user picks. parentId nests them under the
// space column.
func listFoldersAsScopes(apiClient plugin.ApiClient, spaceId, parentId string) ([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpFolder], *ClickUpRemotePagination, errors.Error) {
	body, err := getEntities(apiClient, fmt.Sprintf("space/%s/folder", spaceId))
	if err != nil {
		return nil, nil, err
	}
	entries := make([]dsmodels.DsRemoteApiScopeListEntry[models.ClickUpFolder], 0, len(body.Folders))
	for _, folder := range body.Folders {
		folder := folder
		entries = append(entries, dsmodels.DsRemoteApiScopeListEntry[models.ClickUpFolder]{
			Type:     api.RAS_ENTRY_TYPE_SCOPE,
			ParentId: &parentId,
			Id:       folder.Id,
			Name:     folder.Name,
			FullName: folder.Name,
			Data: &models.ClickUpFolder{
				FolderId:  folder.Id,
				Name:      folder.Name,
				SpaceId:   spaceId,
				SpaceName: "",
			},
		})
	}
	return entries, nil, nil
}

// RemoteScopes lists the ClickUp folders available on the connection so the
// config UI can enumerate selectable scopes.
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// Proxy forwards arbitrary requests to the ClickUp API through the connection.
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
