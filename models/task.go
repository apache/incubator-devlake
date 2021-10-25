package models

import (
	"gorm.io/datatypes"
)

type Task struct {
	Model
	Plugin   string         `json:"plugin" gorm:"index"`
	Options  datatypes.JSON `json:"options"`
	Status   string         `json:"status"`
	Message  string         `json:"message"`
	Progress float32        `json:"progress"`
	SourceId int64          `json:"source_id" gorm:"index"`
}
