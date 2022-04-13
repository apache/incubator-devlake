package archived

import (
	"gorm.io/datatypes"
)

type Blueprint struct {
	Name       string
	Tasks      datatypes.JSON
	Enable     bool
	CronConfig string
	Model
}

func (Blueprint) TableName() string {
	return "_devlake_blueprints"
}
