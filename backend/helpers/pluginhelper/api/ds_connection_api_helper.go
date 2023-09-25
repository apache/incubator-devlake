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
	"strconv"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/srvhelper"
	"github.com/apache/incubator-devlake/server/api/shared"
)

// DsConnectionApiHelper
type DsConnectionApiHelper[C plugin.ToolLayerConnection, S plugin.ToolLayerScope, SC plugin.ToolLayerScopeConfig] struct {
	*ModelApiHelper[C]
	*srvhelper.ConnectionSrvHelper[C, S, SC]
}

func NewDsConnectionApiHelper[C plugin.ToolLayerConnection, S plugin.ToolLayerScope, SC plugin.ToolLayerScopeConfig](
	basicRes context.BasicRes,
	connSrvHelper *srvhelper.ConnectionSrvHelper[C, S, SC],
) *DsConnectionApiHelper[C, S, SC] {
	return &DsConnectionApiHelper[C, S, SC]{
		ModelApiHelper:      NewModelApiHelper[C](basicRes, connSrvHelper.ModelSrvHelper, []string{"connectionId"}),
		ConnectionSrvHelper: connSrvHelper,
	}
}

func (self *DsConnectionApiHelper[C, S, SC]) Delete(input *plugin.ApiResourceInput) (out *plugin.ApiResourceOutput, err errors.Error) {
	var conn *C
	conn, err = self.FindByPk(input)
	if err != nil {
		return nil, err
	}
	refs, err := self.ConnectionSrvHelper.Delete(conn)
	if err != nil {
		return &plugin.ApiResourceOutput{Body: &shared.ApiBody{
			Success: false,
			Message: err.Error(),
			Data:    refs,
		}, Status: err.GetType().GetHttpCode()}, nil
	}
	return &plugin.ApiResourceOutput{
		Body: conn,
	}, nil
}

// DsScopeApiHelper
func extractConnectionId(input *plugin.ApiResourceInput) (uint64, errors.Error) {
	connectionId, ok := input.Params["connectionId"]
	if !ok {
		return 0, errors.BadInput.New("connectionId is required")
	}
	id, err := strconv.ParseUint(connectionId, 10, 64)
	if err != nil {
		return 0, errors.BadInput.Wrap(err, "connectionId must be a number")
	}
	return id, nil
}
