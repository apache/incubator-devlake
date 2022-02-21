package models

import "github.com/merico-dev/lake/models/common"

// JenkinsJobProps current used jenkins job props
type JenkinsJobProps struct {
	Name  string `gorm:"primaryKey;type:varchar(255)"`
	Class string
	Color string
	Base  string
}

// JenkinsJob db entity for jenkins job
type JenkinsJob struct {
	JenkinsJobProps
	common.NoPKModel
}
