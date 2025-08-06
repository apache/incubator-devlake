package models

import (
	"gorm.io/gorm"
)

// Project represents a DevLake project
type Project struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Tags        []Tag  `json:"tags" gorm:"many2many:_devlake_project_tags;"`
}

// Tag represents a tag for a project
type Tag struct {
	gorm.Model
	Name string `json:"name"`
}