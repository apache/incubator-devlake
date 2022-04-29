package models

import (
	"github.com/merico-dev/lake/models/common"
)

type TapdSource struct {
	common.Model
	Name             string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint         string `gorm:"type:varchar(255)"`
	BasicAuthEncoded string `gorm:"type:varchar(255)"`
	RateLimit        int    `comment:"api request rate limt per second"`
}

type TapdSourceDetail struct {
	TapdSource
}

func (TapdSource) TableName() string {
	return "_tool_tapd_sources"
}
