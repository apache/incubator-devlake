package impl

import (
	"fmt"

	"github.com/apache/incubator-devlake/helpers/pluginhelper/subtaskmeta/sorter"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/models/migrationscripts"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/ra_dora/api"
	"github.com/apache/incubator-devlake/plugins/ra_dora/models"
	"github.com/apache/incubator-devlake/plugins/ra_dora/tasks"
)

var _ interface {
	plugin.PluginMeta
	plugin.PluginTask
	plugin.PluginModel
	plugin.PluginSource
	plugin.DataSourcePluginBlueprintV200
	plugin.CloseablePluginTask
} = (*RaDoraMetrics)(nil)

type RaDoraMetrics struct{}

func (r RaDoraMetrics) MakeDataSourcePipelinePlanV200(
	connectionId uint64,
	scopes []*coreModels.BlueprintScope,
	skipCollectors bool,
) (pp coreModels.PipelinePlan, sc []plugin.Scope, err errors.Error) {
	return nil, nil, nil
}

func init() {
	// check subtask meta loop when init subtask meta
	if _, err := sorter.NewDependencySorter(tasks.SubTaskMetaList).Sort(); err != nil {
		panic(err)
	}
}

func (r RaDoraMetrics) Init(br context.BasicRes) errors.Error {
	//api.Init(br, r)

	return nil
}

func (r RaDoraMetrics) Description() string {
	return "Collection Argo data for DORA metrics"
}

func (r RaDoraMetrics) Name() string {
	return "ra_dora"
}

func (r RaDoraMetrics) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/ra_dora"
}

func (r RaDoraMetrics) Connection() dal.Tabler {
	return &models.ArgoConnection{}
}

// TODO
func (r RaDoraMetrics) Scope() plugin.ToolLayerScope {
	return nil
}

func (r RaDoraMetrics) ScopeConfig() dal.Tabler {
	return nil
}

func (r RaDoraMetrics) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{
		&models.Deployment{},
	}
}

func (r RaDoraMetrics) SubTaskMetas() []plugin.SubTaskMeta {
	return []plugin.SubTaskMeta{
		tasks.CollectDeploymentsMeta,
		tasks.ExtractDeploymentsMeta,
		tasks.ConvertDeploymentsMeta,
	}
}

func (r RaDoraMetrics) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	connection := &models.ArgoConnection{}
	err := taskCtx.GetDal().First(connection, dal.Where("id = ?", options["connectionId"]))
	if err != nil {
		return nil, err
	}

	apiClient, err := api.NewApiClient(connection)
	if err != nil {
		return nil, err
	}

	taskData := &tasks.ArgoTaskData{
		ApiClient: apiClient.Client,
	}

	return taskData, nil
}

func (r RaDoraMetrics) MigrationScripts() []plugin.MigrationScript {
	return migrationscripts.All()
}

func (r RaDoraMetrics) TestConnection(id uint64) errors.Error {
	return nil
}

func (r RaDoraMetrics) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{}
}

func (r RaDoraMetrics) Close(taskCtx plugin.TaskContext) errors.Error {
	data, ok := taskCtx.GetData().(*tasks.ArgoTaskData)
	if !ok {
		return errors.Default.New(fmt.Sprintf("GetData failed when try to close %+v", taskCtx))
	}
	if data != nil && data.ApiClient != nil {
		// No need to call Release here as we are not managing connections manually
	}
	return nil
}
