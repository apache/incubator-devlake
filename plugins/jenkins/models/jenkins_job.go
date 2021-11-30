package models

import (
	"github.com/merico-dev/lake/models/common"
)

// JenkinsJobProps current used jenkins job props
type JenkinsJobProps struct {
	Name  string
	Class string
	Color string
}

// JenkinsJob db entity for jenkins job
type JenkinsJob struct {
	common.Model
	JenkinsJobProps
}
