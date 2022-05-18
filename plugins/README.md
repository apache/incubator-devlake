# So you want to Build a New Plugin...

...the good news is, it's easy!

## Preparation work
1. Create a directory named `yourplugin` under directory `plugins`
2. Under `yourplugin`, you need three more packages: `api`, `models` and `tasks`
    1. `api` interacts with `config-ui` for test/get/save connection of data source. Please check [How to create connection to be used by config-ui for a data source]() for detail.
    2. `models` stores all `data entities` and `data migration scripts`. Please check [How to create models and data migrations]() for detail.
    3. `tasks` contains all of our `sub tasks` for a plugin
3. Create a yourplugin.go in `yourplugin`
```golang
type YourPlugin struct{}

var _ core.PluginMeta = (*YourPlugin)(nil)
var _ core.PluginInit = (*YourPlugin)(nil)
var _ core.PluginTask = (*YourPlugin)(nil)
var _ core.PluginApi = (*YourPlugin)(nil)
var _ core.Migratable = (*YourPlugin)(nil)

func (plugin YourPlugin) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	return nil
}

func (plugin YourPlugin) Description() string {
	return "To collect and enrich data from YourPlugin"
}
// Register all subtasks
func (plugin YourPlugin) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectXXXX,
		tasks.ExtractXXXX,
		tasks.ConvertXXXX,
	}
}
// Prepare your apiClient which will be used to request remote api, 
// `apiClient` is defined in `client.go` under `tasks`
// `YourPluginTaskData` is defined in `task_data.go` under `tasks`
func (plugin YourPlugin) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.YourPluginOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}
  // Handle error.
  if err != nil {
    logger.Error(err)
  }

  // Export a variable named PluginEntry for Framework to search and load
  var PluginEntry YourPlugin //nolint
}
```

## Summary

To build a new plugin you will need a few things. You should choose an API that you'd like to see data from. Think about the metrics you would like to see first, and then look for data that can support those metrics.

## Create your sub tasks

1. [Create collector will collect data from remote api server and save into raw layer]()
2. [Create extractor will extract data from raw layer and save into tool layer]()
3. [Create convertor will convert data from tool layer and save into domain layer]()

## You're Done!

Congratulations! You have created your first plugin! 🎖 
