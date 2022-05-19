package archived

import (
	"github.com/apache/incubator-devlake/models/common"
)

type TapdConnection struct {
	common.Model
	Name             string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint         string `gorm:"type:varchar(255)"`
	BasicAuthEncoded string `gorm:"type:varchar(255)"`
	RateLimit        int    `comment:"api request rate limt per second"`
}

type TapdConnectionDetail struct {
	TapdConnection
}

func (TapdConnection) TableName() string {
	return "_tool_tapd_connections"
}
