package impl

import (
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/apache/incubator-devlake/migration"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/feishu/api"
	"github.com/apache/incubator-devlake/plugins/feishu/models"
	"github.com/apache/incubator-devlake/plugins/feishu/models/migrationscripts"
	"github.com/apache/incubator-devlake/plugins/feishu/tasks"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var _ core.PluginMeta = (*Feishu)(nil)
var _ core.PluginInit = (*Feishu)(nil)
var _ core.PluginTask = (*Feishu)(nil)
var _ core.PluginApi = (*Feishu)(nil)
var _ core.Migratable = (*Feishu)(nil)

type Feishu struct{}

func (plugin Feishu) Init(config *viper.Viper, logger core.Logger, db *gorm.DB) error {
	api.Init(config, logger, db)

	// FIXME after config-ui support feishu plugin
	// save env to db where name=feishu
	connection := &models.FeishuConnection{}
	err := db.Find(connection, map[string]string{"name": "Feishu"}).Error
	if err != nil {
		return err
	}
	if connection.ID != 0 {
		encodeKey := config.GetString(core.EncodeKeyEnvStr)
		connection.Endpoint = config.GetString(`FEISHU_ENDPOINT`)
		connection.AppId = config.GetString(`FEISHU_APPID`)
		connection.SecretKey = config.GetString(`FEISHU_APPSCRECT`)
		if connection.Endpoint != `` && connection.AppId != `` && connection.SecretKey != `` && encodeKey != `` {
			err = helper.UpdateEncryptFields(connection, func(plaintext string) (string, error) {
				return core.Encrypt(encodeKey, plaintext)
			})
			if err != nil {
				return err
			}
			// update from .env and save to db
			err = db.Updates(connection).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (plugin Feishu) Description() string {
	return "To collect and enrich data from Feishu"
}

func (plugin Feishu) SubTaskMetas() []core.SubTaskMeta {
	return []core.SubTaskMeta{
		tasks.CollectMeetingTopUserItemMeta,
		tasks.ExtractMeetingTopUserItemMeta,
	}
}

func (plugin Feishu) PrepareTaskData(taskCtx core.TaskContext, options map[string]interface{}) (interface{}, error) {
	var op tasks.FeishuOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return nil, err
	}

	connectionHelper := helper.NewConnectionHelper(
		taskCtx,
		nil,
	)
	connection := &models.FeishuConnection{}
	err = connectionHelper.FirstById(connection, op.ConnectionId)
	if err != nil {
		return err, nil
	}

	apiClient, err := tasks.NewFeishuApiClient(taskCtx, connection)
	if err != nil {
		return nil, err
	}
	return &tasks.FeishuTaskData{
		Options:   &op,
		ApiClient: apiClient,
	}, nil
}

func (plugin Feishu) RootPkgPath() string {
	return "github.com/apache/incubator-devlake/plugins/feishu"
}

func (plugin Feishu) MigrationScripts() []migration.Script {
	return migrationscripts.All()
}

func (plugin Feishu) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{}
}
