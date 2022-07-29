package models

import "github.com/apache/incubator-devlake/models/common"

type JenkinsBuildTriggeredBuilds struct {
	ConnectionId       uint64 `gorm:"primaryKey"`
	BuildName          string `gorm:"primaryKey;type:varchar(255)"`
	TriggeredBuildName string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

func (JenkinsBuildTriggeredBuilds) TableName() string {
	return "_tool_jenkins_build_triggered_builds"
}
