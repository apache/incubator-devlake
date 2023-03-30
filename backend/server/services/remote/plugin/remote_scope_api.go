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

package plugin

import (
	"net/http"
	"strconv"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/server/services/remote/bridge"
)

type RemoteScopesOutput struct {
	Children []RemoteScopesTreeNode `json:"children"`
}

type RemoteScopesTreeNode struct {
	Type     string      `json:"type"`
	ParentId *string     `json:"parentId"`
	Id       string      `json:"id"`
	Name     string      `json:"name"`
	Data     interface{} `json:"data"`
}

func (pa *pluginAPI) GetRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}

	connection := pa.connType.New()
	err := pa.helper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}

	groupId := input.Query.Get("groupId")

	remoteScopes := make([]RemoteScopesTreeNode, 0)
	err = pa.invoker.Call("remote-scopes", bridge.DefaultContext, connection.Unwrap(), groupId).Get(&remoteScopes)
	if err != nil {
		return nil, err
	}

	output := RemoteScopesOutput{
		Children: remoteScopes,
	}

	return &plugin.ApiResourceOutput{Body: output, Status: http.StatusOK}, nil
}

func (pa *pluginAPI) SearchRemoteScopes(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return &plugin.ApiResourceOutput{Status: http.StatusNotImplemented}, nil
}
