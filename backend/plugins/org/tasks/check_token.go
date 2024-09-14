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

package tasks

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"golang.org/x/sync/errgroup"
)

var TaskCheckTokenMeta = plugin.SubTaskMeta{
	Name:             "checkTokens",
	EntryPoint:       checkProjectTokens,
	EnabledByDefault: true,
	Description:      "set project mapping",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

func checkProjectTokens(taskCtx plugin.SubTaskContext) errors.Error {
	logger := taskCtx.GetLogger()
	taskData := taskCtx.GetData().(*TaskData)
	connections := taskData.Options.ProjectConnections
	logger.Debug("connections %+v", connections)
	if len(connections) == 0 {
		return nil
	}

	g := new(errgroup.Group)
	for _, connection := range connections {
		conn := connection
		logger.Debug("check conn: %+v", conn)
		g.Go(func() error {
			if err := checkConnectionToken(conn.ConnectionId, conn.PluginName); err != nil {
				return err
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return errors.Convert(err)
	}
	return nil
}

func checkConnectionToken(connectionID uint64, pluginName string) errors.Error {
	pluginEntry, err := plugin.GetPlugin(pluginName)
	if err != nil {
		return err
	}
	if v, ok := pluginEntry.(plugin.PluginTestConnectionAPI); ok {
		if err := v.TestConnection(connectionID); err != nil {
			return err
		}
		return nil
	} else {
		msg := fmt.Sprintf("plugin: %s doesn't impl test connection api", pluginName)
		return errors.Default.New(msg)
	}
}
