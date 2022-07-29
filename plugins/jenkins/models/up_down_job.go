package models

import "github.com/apache/incubator-devlake/models/common"

type JenkinsUpDownJob struct {
	ConnetionId   uint64 `gorm:"primaryKey"`
	UpstreamJob   string `gorm:"primaryKey;type:varchar(255)"`
	DownstreamJob string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

func (JenkinsUpDownJob) TableName() string {
	return "_tool_jenkins_up_down_jobs"
}
