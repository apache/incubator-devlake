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

package srvhelper

import (
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

type GenericConnectionModelInfo[C plugin.ToolLayerConnection] struct {
	*GenericModelInfo[C]
}

func (*GenericConnectionModelInfo[C]) GetConnectionId(connection any) uint64 {
	return connection.(plugin.ToolLayerConnection).ConnectionId()
}

// ConnectionSrvHelper
type ConnectionSrvHelper[C plugin.ToolLayerConnection, S plugin.ToolLayerScope, SC plugin.ToolLayerScopeConfig] struct {
	*AnyConnectionSrvHelper
	*ModelSrvHelper[C]
}

// NewConnectionSrvHelper creates a ConnectionDalHelper for connection management
func NewConnectionSrvHelper[
	C plugin.ToolLayerConnection,
	S plugin.ToolLayerScope,
	SC plugin.ToolLayerScopeConfig,
](
	anyConnectionSrv *AnyConnectionSrvHelper,
) *ConnectionSrvHelper[C, S, SC] {
	return &ConnectionSrvHelper[C, S, SC]{
		AnyConnectionSrvHelper: anyConnectionSrv,
		ModelSrvHelper:         NewModelSrvHelper[C](anyConnectionSrv.basicRes, nil),
	}
}

func (connSrv *ConnectionSrvHelper[C, S, SC]) DeleteConnection(connection *C) (refs *DsRefs, err errors.Error) {
	return connSrv.DeleteConnectionAny(connection)
}
