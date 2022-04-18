package code

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type Commit struct {
	common.NoPKModel
	Sha            string `json:"sha" gorm:"primaryKey;type:varchar(40);comment:commit hash"`
	Additions      int    `json:"additions" gorm:"comment:Added lines of code"`
	Deletions      int    `json:"deletions" gorm:"comment:Deleted lines of code"`
	DevEq          int    `json:"deveq" gorm:"comment:Merico developer equivalent from analysis engine"`
	Message        string
	AuthorName     string `gorm:"type:varchar(255)"`
	AuthorEmail    string `gorm:"type:varchar(255)"`
	AuthoredDate   time.Time
	AuthorId       string `gorm:"type:varchar(255)"`
	CommitterName  string `gorm:"type:varchar(255)"`
	CommitterEmail string `gorm:"type:varchar(255)"`
	CommittedDate  time.Time
	CommitterId    string `gorm:"index;type:varchar(255)"`
}

type CommitFile struct {
	common.NoPKModel
	CommitSha string `gorm:"primaryKey;type:varchar(40)"`
	FilePath  string `gorm:"primaryKey;type:varchar(255)"`
	Additions int
	Deletions int
}
