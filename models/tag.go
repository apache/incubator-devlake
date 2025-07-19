package models

import (
	"github.com/apache/incubator-devlake/core/models/common"
)

// Tag represents a tag that can be assigned to projects
type Tag struct {
	common.DynamicMapBase
	Name        string `json:"name" gorm:"type:varchar(255);uniqueIndex"`
	Description string `json:"description" gorm:"type:varchar(255)"`
	Color       string `json:"color" gorm:"type:varchar(50)"`
}

// ProjectTag represents the many-to-many relationship between projects and tags
type ProjectTag struct {
	ProjectId string `json:"projectId" gorm:"primaryKey;type:varchar(255)"`
	TagId     string `json:"tagId" gorm:"primaryKey;type:varchar(255)"`
}

// TableName returns the table name for ProjectTag
func (ProjectTag) TableName() string {
	return "_devlake_project_tags"
}
