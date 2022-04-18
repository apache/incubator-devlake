package code

import "github.com/merico-dev/lake/models/common"

type PullRequestCommit struct {
	CommitSha     string `gorm:"primaryKey;type:varchar(40)"`
	PullRequestId string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	common.NoPKModel
}
