package archived

import (
	"time"

	"gorm.io/datatypes"
)

type Pipeline struct {
	Model
	Name        string `json:"name" gorm:"index"`
	BlueprintId uint64
	Tasks       datatypes.JSON
	TotalTasks  int
	// Deprecated
	FinishedTasks int
	BeganAt       *time.Time
	FinishedAt    *time.Time `gorm:"index"`
	Status        string
	Message       string
	SpentSeconds  int
	Step          int
}

func (Pipeline) TableName() string {
	return "_devlake_pipelines"
}
