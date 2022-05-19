package models

import (
	"github.com/apache/incubator-devlake/models/common"
)

type TestConnectionRequest struct {
	Endpoint string `json:"endpoint" validate:"required,url"`
	Auth     string `json:"auth" validate:"required"`
	Proxy    string `json:"proxy"`
}

type WorkspaceResponse struct {
	Id    uint64
	Title string
	Value string
}

type TapdConnection struct {
	common.Model
	Name             string `gorm:"type:varchar(100);uniqueIndex" json:"name" validate:"required"`
	Endpoint         string `gorm:"type:varchar(255)" json:"endpoint"`
	BasicAuthEncoded string `gorm:"type:varchar(255)" json:"basicAuthEncoded"`
	RateLimit        int    `comment:"api request rate limt per hour" json:"rateLimit"`
}

type TapdConnectionDetail struct {
	TapdConnection
}

func (TapdConnection) TableName() string {
	return "_tool_tapd_connections"
}
