package models

import "github.com/apache/incubator-devlake/core/models/common"

type Deployments struct {
	common.NoPKModel

	//TODO
	ID      string `gorm:"primaryKey;type:varchar(255);column:id"`
	ScopeID string
}

func (Deployments) TableName() string {
	return "_tool_cicd_deployments"
}
