package models

import "github.com/apache/incubator-devlake/core/models/common"

type Deployments struct {
	common.NoPKModel

	ID           string `gorm:"primaryKey;type:varchar(255);column:id"`
	ScopeID      string `gorm:"type:varchar(255);column:cicd_scope_id"`
	Name         string `gorm:"type:varchar(255);column:name"`
	Result       string `gorm:"type:varchar(100);column:result"`
	Status       string `gorm:"type:varchar(100);column:status"`
	Environment  string `gorm:"type:varchar(255);column:environment"`
	CreatedDate  string `gorm:"type:datetime(3);column:created_date"`
	StartedDate  string `gorm:"type:datetime(3);column:started_date"`
	FinishedDate string `gorm:"type:datetime(3);column:finished_date"`
	DurationSec  int64  `gorm:"type:bigint;column:duration_sec"`
}

func (Deployments) TableName() string {
	return "_tool_cicd_deployments"
}
