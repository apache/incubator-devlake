package impl

import (
	"fmt"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/plugins/azure/api"
	"github.com/apache/incubator-devlake/plugins/azure/models"
	"github.com/apache/incubator-devlake/plugins/azure/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/azure/tasks"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// make sure interface is implemented
var _ core.PluginMeta = (*Azure)(nil)
var _ core.PluginInit = (*Azure)(nil)
var _ core.PluginTask = (*Azure)(nil)
var _ core.PluginApi = (*Azure)(nil)
var _ core.CloseablePluginTask = (*Azure)(nil)

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Azure //nolint

type Azure struct{}

func (plugin Azure) Description() string {
	return "collect some Azure data"
}

func (plugin Azure) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	api.Init(config, logger, db)
	return nil
}

func (plugin Azure) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectApiRepoMeta,
		tasks.ExtractApiRepoMeta,
		tasks.CollectApiBuildDefinitionMeta,
		tasks.ExtractApiBuildDefinitionMeta,
	}
}

func (plugin Azure) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	op, err := tasks.DecodeAndValidateTaskOptions(options)
	if err != nil {
		return nil, err
	}
	if op.ConnectionId == 0 {
		return nil, errors.Default.New("connectionId is invalid")
	}

	connection := &models.AzureConnection{}
	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return nil, err
	}

	apiClient, err := tasks.CreateApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}
	return &tasks.AzureTaskData{
		Options:    op,
		ApiClient:  apiClient,
		Connection: connection,
	}, nil
}

// PkgPath information lost when compiled as plugin(.so)
func (plugin Azure) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/azure"
}

func (plugin Azure) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"connections": {
			"POST": api.PostConnections,
			"GET":  api.ListConnections,
		},
		"connections/:connectionId": {
			"GET":    api.GetConnection,
			"PATCH":  api.PatchConnection,
			"DELETE": api.DeleteConnection,
		},
	}
}

func (plugin Azure) MigrationScripts() []migration.Script {
	return migrationscripts.All()
}

func (plugin Azure) Close(taskCtx core.TaskContext) error {
	data, ok := taskCtx.GetData().(*tasks.AzureTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	data.ApiClient.Release()
	return nil
}
