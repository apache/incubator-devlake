package code

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type PullRequest struct {
	domainlayer.DomainEntity
	BaseRepoId     string `gorm:"index"`
	HeadRepoId     string `gorm:"index"`
	Status         string `gorm:"type:varchar(100);comment:open/closed or other"`
	Number         int
	Title          string
	Description    string
	Url            string `gorm:"type:char(255)"`
	AuthorName     string `gorm:"type:char(100)"`
	AuthorId       int
	ParentPrId     string `gorm:"index;type:varchar(100)"`
	Key            int
	CreatedDate    time.Time
	MergedDate     *time.Time
	ClosedDate     *time.Time
	Type           string `gorm:"type:char(100)"`
	Component      string `gorm:"type:char(100)"`
	MergeCommitSha string `gorm:"type:char(40)"`
	HeadRef        string `gorm:"type:char(255)"`
	BaseRef        string `gorm:"type:char(255)"`
	BaseCommitSha  string `gorm:"type:char(40)"`
	HeadCommitSha  string `gorm:"type:char(40)"`
}
