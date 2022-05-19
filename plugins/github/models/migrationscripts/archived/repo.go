package archived

import (
	"time"

	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
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
	archived.NoPKModel
}

func (GithubRepo) TableName() string {
	return "_tool_github_repos"
}
