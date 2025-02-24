package impl

import (
	"github.com/apache/incubator-devlake/helpers/pluginhelper/subtaskmeta/sorter"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
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

func (r RaDoraMetrics) Description() string {
	return "Description"
}

func (r RaDoraMetrics) Name() string {
	return "ra_dora"
}

func (r RaDoraMetrics) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/ra_dora"
}

func (r RaDoraMetrics) Connection() dal.Tabler {
	return nil
}

func (r RaDoraMetrics) Scope() plugin.ToolLayerScope {
	return nil
}

func (r RaDoraMetrics) ScopeConfig() dal.Tabler {
	return nil
}

func (r RaDoraMetrics) GetTablesInfo() []dal.Tabler {
	return []dal.Tabler{}
}

func (r RaDoraMetrics) SubTaskMetas() []plugin.SubTaskMeta {
	list, err := sorter.NewDependencySorter(tasks.SubTaskMetaList).Sort()
	if err != nil {
		panic(err)
	}
	return list
}

func (r RaDoraMetrics) PrepareTaskData(taskCtx plugin.TaskContext, options map[string]interface{}) (interface{}, errors.Error) {
	return nil, nil
}

func (r RaDoraMetrics) MigrationScripts() []plugin.MigrationScript {
	return nil
}

func (r RaDoraMetrics) TestConnection(id uint64) errors.Error {
	return nil
}

func (r RaDoraMetrics) ApiResources() map[string]map[string]plugin.ApiResourceHandler {
	return map[string]map[string]plugin.ApiResourceHandler{}
}

func (r RaDoraMetrics) Close(taskCtx plugin.TaskContext) errors.Error {
	return nil
}
