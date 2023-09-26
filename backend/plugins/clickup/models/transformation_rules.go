package models

import "github.com/apache/incubator-devlake/core/models/common"

type ClickUpTransformationRule struct {
	common.Model `mapstructure:"-"`
	ConnectionId uint64 `mapstructure:"connectionId" json:"connectionId"`
	Name         string `mapstructure:"name" json:"name" gorm:"type:varchar(255);index:idx_name_jira,unique" validate:"required"`
}

func (t ClickUpTransformationRule) TableName() string {
	return "_tool_clickup_transformation_rules"
}
