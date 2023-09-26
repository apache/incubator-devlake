package models

import "github.com/apache/incubator-devlake/core/models/common"

type ClickUpSpace struct {
	common.NoPKModel
	ConnectionId  uint64 `json:"connectionId" mapstructure:"connectionId" validate:"required" gorm:"primaryKey"`
	ScopeConfigId uint64 `json:"scopeConfigId,omitempty" mapstructure:"scopeConfigId"`
	Id            string `gorm:"primaryKey" json:"id"`
	Name          string `json:"name"`
}

func (ClickUpSpace) TableName() string {
	return "_tool_clickup_spaces"
}
