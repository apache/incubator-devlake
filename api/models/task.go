package models

import "gorm.io/datatypes"

type Task struct {
	Model
	Plugin   string
	SourceID int
	Options  datatypes.JSON
}
