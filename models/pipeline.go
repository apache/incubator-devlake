package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"

	"gorm.io/datatypes"
)

type Pipeline struct {
	Name          string         `json:"name" gorm:"index"`
	BlueprintId   uint64         `json:"blueprintId"`
	Tasks         datatypes.JSON `json:"tasks"`
	TotalTasks    int            `json:"totalTasks"`
	FinishedTasks int            `json:"finishedTasks"`
	BeganAt       *time.Time     `json:"beganAt"`
	FinishedAt    *time.Time     `json:"finishedAt" gorm:"index"`
	Status        string         `json:"status"`
	Message       string         `json:"message"`
	SpentSeconds  int            `json:"spentSeconds"`
	common.Model
}

// We use a 2D array because the request body must be an array of a set of tasks
// to be executed concurrently, while each set is to be executed sequentially.
type NewPipeline struct {
	Name        string       `json:"name"`
	Tasks       [][]*NewTask `json:"tasks"`
	BlueprintId uint64
}
