package main

import (
	"fmt"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/tapd/api"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"github.com/merico-dev/lake/plugins/tapd/tasks"
	"github.com/merico-dev/lake/runner"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"time"
)

var _ core.PluginMeta = (*Tapd)(nil)
var _ core.PluginInit = (*Tapd)(nil)
var _ core.PluginTask = (*Tapd)(nil)
var _ core.PluginApi = (*Tapd)(nil)

type Tapd struct{}

func (plugin Tapd) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	api.Init(config, logger, db)
	return db.AutoMigrate(
		&models.TapdSource{},
		&models.TapdWorkspace{},
		&models.TapdUser{},
		&models.TapdIteration{},
		&models.TapdStory{},
		&models.TapdBug{},
		&models.TapdTask{},
	)
}

func (plugin Tapd) Description() string {
	return "To collect and enrich data from Tapd"
}

func (plugin Tapd) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		//tasks.CollectWorkspaceMeta,
		//tasks.ExtractWorkspaceMeta,
		//tasks.ConvertWorkspaceMeta,
		tasks.CollectUserMeta,
		tasks.ExtractUserMeta,
		tasks.ConvertUserMeta,
		//tasks.CollectIterationMeta,
		//tasks.ExtractIterationMeta,
		//tasks.ConvertIterationMeta,
		tasks.CollectStoryMeta,
		tasks.ExtractStoryMeta,
		tasks.ConvertStoryMeta,
		tasks.CollectBugMeta,
		tasks.ExtractBugMeta,
		tasks.ConvertBugMeta,
		tasks.CollectTaskMeta,
		tasks.ExtractTaskMeta,
		tasks.ConvertTaskMeta,
	}
}

func (plugin Tapd) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	db := taskCtx.GetDb()
	var op tasks.TapdOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
	if op.SourceId == 0 {
		return nil, fmt.Errorf("SourceId is required for Tapd execution")
	}
	source := &models.TapdSource{}
	err = db.Find(source, op.SourceId).Error
	if err != nil {
		return nil, err
	}
	var since time.Time
	if op.Since != "" {
		since, err = time.Parse("2006-01-02T15:04:05Z", op.Since)
		if err != nil {
			return nil, fmt.Errorf("invalid value for `since`: %w", err)
		}
	}
	tapdApiClient, err := tasks.NewTapdApiClient(taskCtx, source)
	if err != nil {
		return nil, fmt.Errorf("failed to create tapd api client: %w", err)
	}
	taskData := &tasks.TapdTaskData{
		Options:   &op,
		ApiClient: tapdApiClient,
		Source:    source,
		Since:     &since,
	}
	tasks.UserIdGen = didgen.NewDomainIdGenerator(&models.TapdUser{})
	tasks.WorkspaceIdGen = didgen.NewDomainIdGenerator(&models.TapdWorkspace{})
	return taskData, nil
}

func (plugin Tapd) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/tapd"
}

func (plugin Tapd) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"POST": api.TestConnection,
		},
		"sources": {
			"POST": api.PostSources,
			"GET":  api.ListSources,
		},
		"sources/:sourceId": {
			"PUT":    api.PutSource,
			"DELETE": api.DeleteSource,
			"GET":    api.GetSource,
		},
	}
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Tapd //nolint

// standalone mode for debugging
func main() {
	cmd := &cobra.Command{Use: "tapd"}
	sourceId := cmd.Flags().Uint64P("source", "s", 0, "tapd source id")
	companyId := cmd.Flags().Uint64P("company", "c", 0, "tapd company id")
	workspaceId := cmd.Flags().Uint64P("workspace", "w", 0, "tapd workspace id")
	err := cmd.MarkFlagRequired("source")
	if err != nil {
		panic(err)
	}
	//err = cmd.MarkFlagRequired("company")
	//if err != nil {
	//	panic(err)
	//}
	err = cmd.MarkFlagRequired("workspace")
	if err != nil {
		panic(err)
	}
	cmd.Run = func(c *cobra.Command, args []string) {
		runner.DirectRun(c, args, PluginEntry, map[string]interface{}{
			"sourceId":    *sourceId,
			"companyId":   *companyId,
			"workspaceId": *workspaceId,
		})
	}
	runner.RunCmd(cmd)
}
