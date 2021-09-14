package models

import (
	"gorm.io/datatypes"
)

type Task struct {
	Model
	Plugin   string         `json:"plugin"`
	Options  datatypes.JSON `json:"options"`
	Status   string         `json:"status"`
	Message  string         `json:"message"`
	Progress float32        `json:"progress"`
}
