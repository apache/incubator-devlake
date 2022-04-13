package archived

import (
	"github.com/merico-dev/lake/models/migrationscripts/archived"
)

type AECommit struct {
	HexSha      string `gorm:"primaryKey;type:varchar(255)"`
	AnalysisId  string `gorm:"type:varchar(255)"`
	AuthorEmail string `gorm:"type:varchar(255)"`
	DevEq       int
	AEProjectId int
	archived.NoPKModel
}

func (AECommit) TableName() string {
	return "_tool_ae_commits"
}
