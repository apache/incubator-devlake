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
	"github.com/apache/incubator-devlake/core/plugin"
)

// ScopeConfigSrvHelper
type ScopeConfigSrvHelper[C plugin.ToolLayerConnection, S plugin.ToolLayerScope, SC plugin.ToolLayerScopeConfig] struct {
	*ModelSrvHelper[SC]
}

func NewScopeConfigSrvHelper[
	C plugin.ToolLayerConnection,
	S plugin.ToolLayerScope,
	SC plugin.ToolLayerScopeConfig,
](basicRes context.BasicRes) *ScopeConfigSrvHelper[C, S, SC] {
	return &ScopeConfigSrvHelper[C, S, SC]{
		ModelSrvHelper: NewModelSrvHelper[SC](basicRes),
	}
}

func (self *ScopeConfigSrvHelper[C, S, SC]) GetAllByConnectionId(connectionId uint64) ([]*SC, errors.Error) {
	var scopeConfigs []*SC
	err := self.db.All(&scopeConfigs,
		dal.Where("connection_id = ?", connectionId),
		dal.Orderby("id DESC"),
	)
	return scopeConfigs, err
}

func (self *ScopeConfigSrvHelper[C, S, SC]) Delete(scopeConfig *SC) (refs []*S, err errors.Error) {
	err = self.ModelSrvHelper.NoRunningPipeline(func(tx dal.Transaction) errors.Error {
		// make sure no scope is using the scopeConfig
		sc := (*scopeConfig)
		errors.Must(tx.All(
			&refs,
			dal.Where("connection_id = ? AND scope_config_id = ?", sc.ScopeConfigConnectionId(), sc.ScopeConfigId()),
		))
		if len(refs) > 0 {
			return errors.Conflict.New("Please delete all data scope(s) before you delete this ScopeConfig.")
		}
		errors.Must(tx.Delete(scopeConfig))
		return nil
	})
	return
}
