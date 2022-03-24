package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type GithubRepo struct {
	GithubId       int `gorm:"primaryKey"`
	Name           string
	HTMLUrl        string
	Description    string
	OwnerId        int        `json:"ownerId"`
	OwnerLogin     string     `json:"ownerLogin" gorm:"type:varchar(255)"`
	Language       string     `json:"language" gorm:"type:varchar(255)"`
	ParentGithubId int        `json:"parentId"`
	ParentHTMLUrl  string     `json:"parentHtmlUrl"`
	CreatedDate    time.Time  `json:"createdDate"`
	UpdatedDate    *time.Time `json:"updatedDate"`
	common.NoPKModel
}
