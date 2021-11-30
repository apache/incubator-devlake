package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

// JenkinsBuildProps current used jenkins build props
type JenkinsBuildProps struct {
	Duration          float64 // build time
	DisplayName       string  // "#7"
	EstimatedDuration float64
	Number            int64
	Result            string
	Timestamp         int64     // start time
	StartTime         time.Time // convered by timestamp
}

// JenkinsBuild db entity for jenkins build
type JenkinsBuild struct {
	common.Model
	JenkinsBuildProps
	JobID uint64
}
