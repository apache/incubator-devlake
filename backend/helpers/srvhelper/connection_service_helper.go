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
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
)

// ConnectionSrvHelper
type ConnectionSrvHelper[C plugin.ToolLayerConnection, S plugin.ToolLayerScope, SC plugin.ToolLayerScopeConfig] struct {
	*ModelSrvHelper[C]
	pluginName string
}

// NewConnectionSrvHelper creates a ConnectionDalHelper for connection management
func NewConnectionSrvHelper[
	C plugin.ToolLayerConnection,
	S plugin.ToolLayerScope,
	SC plugin.ToolLayerScopeConfig,
](
	basicRes context.BasicRes,
	pluginName string,
) *ConnectionSrvHelper[C, S, SC] {
	return &ConnectionSrvHelper[C, S, SC]{
		ModelSrvHelper: NewModelSrvHelper[C](basicRes),
		pluginName:     pluginName,
	}
}

func (self *ConnectionSrvHelper[C, S, SC]) Delete(connection *C) (refs *DsRefs, err errors.Error) {
	err = self.ModelSrvHelper.NoRunningPipeline(func(tx dal.Transaction) errors.Error {
		// make sure no blueprint is using the connection
		connectionId := (*connection).ConnectionId()
		refs, err = toDsRefs(self.getAllBlueprinsByConnection(connectionId))
		if err != nil {
			return err
		}
		scopeCount := errors.Must1(self.db.Count(dal.From(new(S)), dal.Where("connection_id = ?", connectionId)))
		if scopeCount > 0 {
			return errors.Conflict.New("Please delete all data scope(s) before you delete this Data Connection.")
		}
		errors.Must(tx.Delete(connection))
		errors.Must(self.db.Delete(new(SC), dal.Where("connection_id = ?", connectionId)))
		return nil
	})
	return
}

func (self *ConnectionSrvHelper[C, S, SC]) getAllBlueprinsByConnection(connectionId uint64) []*models.Blueprint {
	blueprints := make([]*models.Blueprint, 0)
	errors.Must(self.db.All(
		&blueprints,
		dal.From("_devlake_blueprints bp"),
		dal.Join("JOIN _devlake_blueprint_connections cn ON cn.blueprint_id = bp.id"),
		dal.Where(
			"mode = ? AND cn.connection_id = ? AND cn.plugin_name = ?",
			"NORMAL",
			connectionId,
			self.pluginName,
		),
	))
	return blueprints
}
