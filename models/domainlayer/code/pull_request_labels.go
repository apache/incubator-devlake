package code

import (
	"github.com/apache/incubator-devlake/models/common"
)

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type PullRequestLabel struct {
	PullRequestId string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	LabelName     string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}
