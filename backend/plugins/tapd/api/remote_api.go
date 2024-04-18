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
	"net/url"
	"sort"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	dsmodels "github.com/apache/incubator-devlake/helpers/pluginhelper/api/models"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

func listTapdRemoteScopes(
	connection *models.TapdConnection,
	apiClient plugin.ApiClient,
	groupId string,
	page srvhelper.NoPagintation,
) (
	children []dsmodels.DsRemoteApiScopeListEntry[models.TapdWorkspace],
	nextPage *srvhelper.NoPagintation,
	err errors.Error,
) {
	// construct the query and request
	query := url.Values{}
	query.Set("company_id", fmt.Sprintf("%v", connection.CompanyId))
	res, err := apiClient.Get("/workspaces/projects", query, nil)
	if err != nil {
		return
	}
	// parse the response
	var resBody models.WorkspacesResponse
	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return nil, nil, err
	}
	if resBody.Status != 1 {
		return nil, nil, errors.BadInput.Wrap(err, "failed to get workspaces")
	}
	// tapd returns the whole freaking tree as a list, well...let's convert it to a tree
	nodes := map[string]*Node{}
	// convert the list to nodes
	for i, workspace := range resBody.Data {
		entry := &Entry{
			Type:     api.RAS_ENTRY_TYPE_SCOPE, // default to scope
			Id:       fmt.Sprintf(`%d`, workspace.TapdWorkspace.Id),
			Name:     workspace.TapdWorkspace.Name,
			ParentId: toStringPointer(workspace.TapdWorkspace.ParentId),
			Data:     &resBody.Data[i].TapdWorkspace,
		}
		nodes[entry.Id] = &Node{
			entry: entry,
		}
	}
	// construct the tree
	var root *Node
	var current *Node
	for _, node := range nodes {
		// find parent and make sure it is a parent if any
		parent := nodes[*node.entry.ParentId]
		if parent != nil {
			// make sure the parent is a group
			parent.entry.Type = api.RAS_ENTRY_TYPE_GROUP
			parent.entry.Data = nil
			// add this node to the parent
			parent.children = append(parent.children, node)
		} else {
			// or this is root and it must be a group
			root = node
			root.entry.Type = api.RAS_ENTRY_TYPE_GROUP
			root.entry.Data = nil
		}
		if groupId == node.entry.Id {
			current = node
		}
	}
	generateFullNames(root, "")
	// failback to the root
	if current == nil {
		// the parentId in the first page is always nil
		for _, child := range root.children {
			child.entry.ParentId = nil
		}
		current = root
	}
	// sort children
	sort.Sort(current.children)
	// append to the final result
	for _, node := range current.children {
		children = append(children, *node.entry)
	}
	return
}

type Entry = dsmodels.DsRemoteApiScopeListEntry[models.TapdWorkspace]
type Node struct {
	entry    *Entry
	children Children
}
type Children []*Node

func (a Children) Len() int { return len(a) }
func (a Children) Less(i, j int) bool {
	if a[i].entry.Type != a[j].entry.Type {
		return a[i].entry.Type < a[j].entry.Type
	}
	return a[i].entry.Name < a[j].entry.Name
}
func (a Children) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func generateFullNames(node *Node, prefix string) {
	for _, child := range node.children {
		child.entry.FullName = prefix + child.entry.Name
		if child.entry.Type == api.RAS_ENTRY_TYPE_GROUP {
			generateFullNames(child, child.entry.FullName+" / ")
		}
	}
}

func toStringPointer(v any) *string {
	s := fmt.Sprintf("%v", v)
	return &s
}

// RemoteScopes list all available scope for users
// @Summary list all available scope for users
// @Description list all available scope for users
// @Tags plugins/tapd
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param groupId query string false "group ID"
// @Param pageToken query string false "page Token"
// @Success 200  {object} dsmodels.DsRemoteApiScopeList[models.TapdWorkspace]
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/tapd/connections/{connectionId}/remote-scopes [GET]
func RemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raScopeList.Get(input)
}

// @Summary Remote server API proxy
// @Description Forward API requests to the specified remote server
// @Param connectionId path int true "connection ID"
// @Param path path string true "path to a API endpoint"
// @Tags plugins/github
// @Router /plugins/github/connections/{connectionId}/proxy/{path} [GET]
func Proxy(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return raProxy.Proxy(input)
}
