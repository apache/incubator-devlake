package ticket

import "github.com/merico-dev/lake/models/common"

// Please note that Issue Labels can also apply to Pull Requests.
// Pull Requests are considered Issues in GitHub.

type IssueLabel struct {
	IssueId   string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	LabelName string `gorm:"primaryKey"`
	common.NoPKModel
}
