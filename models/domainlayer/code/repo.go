package code

import (
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
)

type Repo struct {
	domainlayer.DomainEntity
	Name        string     `json:"name"`
	Url         string     `json:"url"`
	Description string     `json:"description"`
	OwnerId     string     `json:"ownerId" gorm:"type:varchar(255)"`
	Language    string     `json:"language" gorm:"type:varchar(255)"`
	ForkedFrom  string     `json:"forkedFrom"`
	CreatedDate time.Time  `json:"createdDate"`
	UpdatedDate *time.Time `json:"updatedDate"`
	Deleted     bool       `json:"deleted"`
}

type RepoLanguage struct {
	RepoId   string `json:"repoId" gorm:"index;type:varchar(255)"`
	Language string `json:"language" gorm:"type:varchar(255)"`
	Bytes    int
}
