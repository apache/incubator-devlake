package code

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type PullRequest struct {
	domainlayer.DomainEntity
	RepoId         string `gorm:"index"`
	Status         string `gorm:"comment:open/closed or other"`
	Title          string
	Description    string
	Url            string `gorm:"type:char(255)"`
	AuthorName     string `gorm:"type:char(100)"`
	AuthorId       int
	ParentPrId     string `gorm:"index;type:varchar(100)"`
	Key            int
	CreatedDate    time.Time
	MergedDate     *time.Time
	ClosedAt       *time.Time
	Type           string
	Component      string
	MergeCommitSha string `gorm:"type:char(40)"`
	HeadRef        string
	BaseRef        string
	BaseCommitSha  string
	HeadCommitSha  string
}
