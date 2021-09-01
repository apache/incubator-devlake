package models

import "github.com/merico-dev/lake/models"

// JenkinsBuildProps current used jenkins build props
type JenkinsBuildProps struct {
	Duration          float64 // build time
	DisplayName       string  // "#7"
	EstimatedDuration float64
	Number            int64
	Result            string
	Timestamp         int64 // start time
}

// JenkinsBuild db entity for jenkins build
type JenkinsBuild struct {
	models.Model
	JenkinsBuildProps
	JobID uint64
}
