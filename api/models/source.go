package models

import (
	"gorm.io/datatypes"
)

type Source struct {
	Model
	Plugin  string
	Name    string
	Options datatypes.JSON
}
