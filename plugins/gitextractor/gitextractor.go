package main

import (
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitextractor/tasks"
	"github.com/mitchellh/mapstructure"
)

var _ core.PluginMeta = (*GitExtractor)(nil)
var _ core.PluginTask = (*GitExtractor)(nil)

type GitExtractor struct{}

func (plugin GitExtractor) Description() string {
	return "extract infos from git repository"
}

// return all available subtasks, framework will run them for you in order
func (plugin GitExtractor) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectGitRepoMeta,
	}
}

// based on task context and user input options, return data that shared among all subtasks
func (plugin GitExtractor) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.GitExtractorOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	err = op.Valid()
	if err != nil {
		return nil, err
	}
	return op, nil
}

func (plugin GitExtractor) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/gitextractor"
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry GitExtractor //nolint
