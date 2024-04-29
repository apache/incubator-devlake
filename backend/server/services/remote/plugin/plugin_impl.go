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
	"reflect"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/server/services/remote/bridge"
	"github.com/apache/incubator-devlake/server/services/remote/models"
	"github.com/apache/incubator-devlake/server/services/remote/models/migrationscripts"
	"github.com/apache/incubator-devlake/server/services/remote/plugin/doc"
)

type (
	remotePluginImpl struct {
		name                 string
		subtaskMetas         []plugin.SubTaskMeta
		pluginPath           string
		description          string
		invoker              bridge.Invoker
		connectionModelInfo  *models.RemoteConnectionModelInfo
		scopeModelInfo       *models.RemoteScopeModelInfo
		scopeConfigModelInfo *models.RemoteScopeConfigModelInfo
		toolModelInfos       []dal.Tabler
		migrationScripts     []plugin.MigrationScript
		resources            map[string]map[string]plugin.ApiResourceHandler
		openApiSpec          string
		dsHelper             *api.DsAnyHelper
	}
	RemotePluginTaskData struct {
		DbUrl       string                 `json:"db_url"`
		Scope       interface{}            `json:"scope"`
		Connection  interface{}            `json:"connection"`
		ScopeConfig interface{}            `json:"scope_config"`
		Options     map[string]interface{} `json:"options"`
	}
)

func newPlugin(info *models.PluginInfo, invoker bridge.Invoker) (*remotePluginImpl, errors.Error) {
	// connectionTabler, err := info.ConnectionModelInfo.LoadDynamicTabler(common.Model{})
	connectionModelInfo, err := models.NewRemoteConnectionModelInfo[common.Model](info.ConnectionModelInfo)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("Couldn't load Connection type for plugin %s", info.Name))
	}
	// scopeTabler, err := info.ScopeModelInfo.LoadDynamicTabler(models.ScopeModel{})
	scopeModelInfo, err := models.NewRemoteScopeModelInfo[models.ScopeModel](info.ScopeModelInfo)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("Couldn't load Scope type for plugin %s", info.Name))
	}
	// scopeConfigTabler, err := info.ScopeConfigModelInfo.LoadDynamicTabler(models.ScopeConfigModel{})
	scopeConfigModelInfo, err := models.NewRemoteScopeConfigModelInfo[models.ScopeConfigModel](info.ScopeConfigModelInfo)
	if err != nil {
		return nil, errors.Default.Wrap(err, fmt.Sprintf("Couldn't load ScopeConfig type for plugin %s", info.Name))
	}
	// put the scope and connection models in the tool list to be consistent with Go plugins
	toolModelInfos := []dal.Tabler{
		connectionModelInfo,
		scopeModelInfo,
		scopeConfigModelInfo,
	}
	for _, toolModelInfo := range info.ToolModelInfos {
		mi, err := models.GenerateRemoteModelInfo[models.ToolModel](toolModelInfo)
		if err != nil {
			return nil, err
		}
		toolModelInfos = append(toolModelInfos, mi)
	}
	openApiSpec, err := doc.GenerateOpenApiSpec(info)
	if err != nil {
		panic(err)
	}
	scripts := make([]plugin.MigrationScript, 0)
	for _, script := range info.MigrationScripts {
		script := script
		scripts = append(scripts, &script)
	}
	p := remotePluginImpl{
		name:                 info.Name,
		invoker:              invoker,
		pluginPath:           info.PluginPath,
		description:          info.Description,
		connectionModelInfo:  connectionModelInfo,
		scopeModelInfo:       scopeModelInfo,
		scopeConfigModelInfo: scopeConfigModelInfo,
		toolModelInfos:       toolModelInfos,
		migrationScripts:     scripts,
		openApiSpec:          *openApiSpec,
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

func (p *remotePluginImpl) Init(basicRes context.BasicRes) errors.Error {
	p.dsHelper = api.NewDataSourceAnyHelper(
		basicRes,
		p.Name(),
		[]string{"name"},
		func(c any) any {
			reflect.ValueOf(c).Elem().FieldByName("token").SetString("")
			return c
		},
		p.connectionModelInfo,
		p.scopeModelInfo,
		p.scopeConfigModelInfo,
	)
	p.resources = GetDefaultAPI(p.invoker, p.dsHelper)
	return nil
}

func (p *remotePluginImpl) SubTaskMetas() []plugin.SubTaskMeta {
	return p.subtaskMetas
}

func (p *remotePluginImpl) GetTablesInfo() []dal.Tabler {
	return p.toolModelInfos
}

func (p *remotePluginImpl) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	dbUrl := taskCtx.GetConfig("db_url")
	connectionId := uint64(options["connectionId"].(float64))
	connection, err := p.dsHelper.ConnSrv.FindByPkAny(connectionId)
	if err != nil {
		return nil, err
	}

	scopeId, ok := options["scopeId"].(string)
	if !ok {
		return nil, errors.BadInput.New("missing scopeId")
	}
	scopeDetail, err := p.dsHelper.ScopeSrv.GetScopeDetailAny(false, connectionId, scopeId)
	if err != nil {
		return nil, err
	}

	return RemotePluginTaskData{
		DbUrl:       dbUrl,
		Scope:       scopeDetail.Scope,
		Connection:  connection,
		ScopeConfig: scopeDetail.ScopeConfig,
		Options:     options,
	}, nil
}

func (p *remotePluginImpl) Description() string {
	return p.description
}

func (p *remotePluginImpl) Name() string {
	return p.name
}

func (p *remotePluginImpl) RootPkgPath() string {
	// RootPkgPath is used by DomainIdGenerator to find the name of the plugin that defines a given type.
	// While remote plugins do not use the DomainIdGenerator, we still need to implement this function.
	// Indeed, DomainIdGenerator uses FindPluginNameBySubPkgPath that returns the name of the first plugin
	// whose RootPkgPath is a prefix of the type package path.
	// So we forge a fake package path that is not a prefix of any go plugin package path.
	return "github.com/apache/incubator-devlake/services/remote/fakepackages/" + p.name
}

func (p *remotePluginImpl) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return p.resources
}

func (p *remotePluginImpl) OpenApiSpec() string {
	return p.openApiSpec
}

func (p *remotePluginImpl) MigrationScripts() []plugin.MigrationScript {
	return append(p.migrationScripts, migrationscripts.All(p.name)...)
}

var _ models.RemotePlugin = (*remotePluginImpl)(nil)
