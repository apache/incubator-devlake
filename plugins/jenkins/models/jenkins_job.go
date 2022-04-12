package models

import "github.com/merico-dev/lake/models/common"

// JenkinsJobProps current used jenkins job props
type JenkinsJobProps struct {
	Name  string `gorm:"primaryKey;type:varchar(255)"`
	Class string `gorm:"type:varchar(255)"`
	Color string `gorm:"type:varchar(255)"`
	Base  string `gorm:"type:varchar(255)"`
}

// JenkinsJob db entity for jenkins job
type JenkinsJob struct {
	JenkinsJobProps
	common.NoPKModel
}

func (JenkinsJob) TableName() string {
	return "_tool_jenkins_jobs"
}
