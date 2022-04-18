package code

import "github.com/merico-dev/lake/models/common"

type RepoCommit struct {
	RepoId    string `json:"repoId" gorm:"primaryKey;type:varchar(255)"`
	CommitSha string `json:"commitSha" gorm:"primaryKey;type:varchar(40)"`
	common.NoPKModel
}
