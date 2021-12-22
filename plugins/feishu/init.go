package main

import (
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/feishu/models"
)

func (plugin Feishu) Init() {
	logger.Info("INFO >>> init Feishu plugin", true)
	err := lakeModels.Db.AutoMigrate(
		&models.FeishuMeetingTopUserItem{},
	)
	if err != nil {
		logger.Error("Error migrating feishu: ", err)
		panic(err)
	}
}
