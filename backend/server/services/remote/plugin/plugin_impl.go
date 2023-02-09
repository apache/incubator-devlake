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
	"fmt"

	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/server/services/remote/bridge"
	"github.com/apache/incubator-devlake/server/services/remote/models"
)

type (
	remotePluginImpl struct {
		name             string
		subtaskMetas     []plugin.SubTaskMeta
		pluginPath       string
		description      string
		invoker          bridge.Invoker
		connectionTabler *coreModels.DynamicTabler
		resources        map[string]map[string]plugin.ApiResourceHandler
	}
	RemotePluginTaskData struct {
		DbUrl        string                 `json:"db_url"`
		ConnectionId uint64                 `json:"connection_id"`
		Connection   interface{}            `json:"connection"`
		Options      map[string]interface{} `json:"options"`
	}
)

func newPlugin(info *models.PluginInfo, invoker bridge.Invoker) (*remotePluginImpl, errors.Error) {
	connectionTableName := fmt.Sprintf("_tool_%s_connections", info.Name)
	dynamicTabler, err := models.LoadTableModel(connectionTableName, info.ConnectionSchema)
	if err != nil {
		return nil, err
	}
	p := remotePluginImpl{
		name:             info.Name,
		invoker:          invoker,
		pluginPath:       info.PluginPath,
		description:      info.Description,
		connectionTabler: dynamicTabler,
		resources:        GetDefaultAPI(invoker, dynamicTabler, connectionHelper),
	}
	remoteBridge := bridge.NewBridge(invoker)
	for _, subtask := range info.SubtaskMetas {
		p.subtaskMetas = append(p.subtaskMetas, plugin.SubTaskMeta{
			Name:             subtask.Name,
			EntryPoint:       remoteBridge.RemoteSubtaskEntrypointHandler(subtask),
			Required:         subtask.Required,
			EnabledByDefault: subtask.EnabledByDefault,
			Description:      subtask.Description,
			DomainTypes:      subtask.DomainTypes,
		})
	}
	return &p, nil
}

func (p *remotePluginImpl) SubTaskMetas() []plugin.SubTaskMeta {
	return p.subtaskMetas
}

func (p *remotePluginImpl) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	dbUrl := taskCtx.GetConfig("db_url")
	connectionId := uint64(options["connectionId"].(float64))

	connectionHelper := api.NewConnectionHelper(
		taskCtx,
		nil,
	)

	connection := p.connectionTabler.New()
	err := connectionHelper.FirstById(connection, connectionId)
	if err != nil {
		return nil, errors.Convert(err)
	}

	return RemotePluginTaskData{
		DbUrl:        dbUrl,
		ConnectionId: connectionId,
		Connection:   connection.Unwrap(),
		Options:      options,
	}, nil
}

func (p *remotePluginImpl) Description() string {
	return p.description
}

func (p *remotePluginImpl) RootPkgPath() string {
	// RootPkgPath is only used to find to which plugin a given type belongs.
	// This in turn is only used by DomainIdGenerator.
	// Remote plugins define tool layer types in another language,
	// so the reflective implementation of NewDomainIdGenerator cannot work.
	return ""
}

func (p *remotePluginImpl) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return p.resources
}

func (p *remotePluginImpl) RunMigrations(forceMigrate bool) errors.Error {
	err := api.CallDB(basicRes.GetDal().AutoMigrate, p.connectionTabler.New())
	if err != nil {
		return err
	}

	err = p.invoker.Call("run-migrations", bridge.DefaultContext, forceMigrate).Get()
	return err
}

var _ models.RemotePlugin = (*remotePluginImpl)(nil)
var _ plugin.Scope = (*models.WrappedPipelineScope)(nil)
