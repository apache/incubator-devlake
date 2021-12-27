package code

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type Commit struct {
	common.NoPKModel
	Sha            string `json:"sha" gorm:"primaryKey;comment:commit hash"`
	Additions      int    `json:"additions" gorm:"comment:Added lines of code"`
	Deletions      int    `json:"deletions" gorm:"comment:Deleted lines of code"`
	DevEq          int    `json:"deveq" gorm:"comment:Merico developer equivalent from analysis engine"`
	Message        string
	AuthorName     string
	AuthorEmail    string
	AuthoredDate   time.Time
	AuthorId       string `gorm:"index;type:varchar(255)"`
	CommitterName  string
	CommitterEmail string
	CommittedDate  time.Time
	CommiterId     string `gorm:"index;type:varchar(255)"`
}
