package archived

import (
	"time"
)

type Commit struct {
	NoPKModel
	Sha            string `json:"sha" gorm:"primaryKey;type:varchar(40);comment:commit hash"`
	Additions      int    `gorm:"comment:Added lines of code"`
	Deletions      int    `gorm:"comment:Deleted lines of code"`
	DevEq          int    `gorm:"comment:Merico developer equivalent from analysis engine"`
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
	NoPKModel
	CommitSha string `gorm:"primaryKey;type:varchar(40)"`
	FilePath  string `gorm:"primaryKey;type:varchar(255)"`
	Additions int
	Deletions int
}

type CommitParent struct {
	CommitSha       string `json:"commitSha" gorm:"primaryKey;type:varchar(40);comment:commit hash"`
	ParentCommitSha string `json:"parentCommitSha" gorm:"primaryKey;type:varchar(40);comment:parent commit hash"`
}
