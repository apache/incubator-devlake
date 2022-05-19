package code

import "github.com/apache/incubator-devlake/models/common"

type RepoCommit struct {
	RepoId    string `json:"repoId" gorm:"primaryKey;type:varchar(255)"`
	CommitSha string `json:"commitSha" gorm:"primaryKey;type:varchar(40)"`
	common.NoPKModel
}
