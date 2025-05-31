package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type addTags struct{}

func (*addTags) Name() string {
	return "Add tag tables for project tagging"
}

func (*addTags) Up(baseRes context.BasicRes) errors.Error {
	db := baseRes.GetDal()
	
	err := db.AutoMigrate(&Tag{}, &ProjectTag{})
	if err != nil {
		return errors.Convert(err)
	}
	
	return nil
}

// Tag model for migration
type Tag struct {
	ID          string `gorm:"primaryKey;type:varchar(255)"`
	Name        string `gorm:"type:varchar(255);uniqueIndex"`
	Description string `gorm:"type:varchar(255)"`
	Color       string `gorm:"type:varchar(50)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Tag) TableName() string {
	return "_devlake_tags"
}

// ProjectTag model for migration
type ProjectTag struct {
	ProjectId string `gorm:"primaryKey;type:varchar(255)"`
	TagId     string `gorm:"primaryKey;type:varchar(255)"`
}

func (ProjectTag) TableName() string {
	return "_devlake_project_tags"
}

func init() {
	migrationhelper.RegisterMigration(&addTags{})
}
