package models

import (
	"gorm.io/datatypes"
)

type Source struct {
	Model
	Type    string
	Options datatypes.JSON
}
