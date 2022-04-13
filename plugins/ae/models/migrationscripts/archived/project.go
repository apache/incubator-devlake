package archived

import (
	"time"

	"github.com/merico-dev/lake/models/migrationscripts/archived"
)

type AEProject struct {
	Id           int    `gorm:"primaryKey;type:varchar(255)"`
	GitUrl       string `gorm:"type:varchar(255);comment:url of the repo in github"`
	Priority     int
	AECreateTime *time.Time
	AEUpdateTime *time.Time
	archived.NoPKModel
}

func (AEProject) TableName() string {
	return "_tool_ae_projects"
}
