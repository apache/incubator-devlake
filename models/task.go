package models

import (
	"gorm.io/datatypes"
)

type Task struct {
	Model
	Plugin  string
	Options datatypes.JSON
	Status  string
	Message string
}
