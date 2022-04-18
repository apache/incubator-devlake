package archived

import (
	"time"
)

type Board struct {
	DomainEntity
	Name        string `gorm:"type:varchar(255)"`
	Description string
	Url         string `gorm:"type:varchar(255)"`
	CreatedDate *time.Time
}

type BoardSprint struct {
	NoPKModel
	BoardId  string `gorm:"primaryKey;type:varchar(255)"`
	SprintId string `gorm:"primaryKey;type:varchar(255)"`
}

type BoardIssue struct {
	BoardId string `gorm:"primaryKey;type:varchar(255)"`
	IssueId string `gorm:"primaryKey;type:varchar(255)"`
	NoPKModel
}

type BoardRepo struct {
	BoardId string `gorm:"primaryKey;type:varchar(255)"`
	RepoId  string `gorm:"primaryKey;type:varchar(255)"`
}
