package crossdomain

import "github.com/apache/incubator-devlake/models/common"

type PullRequestIssue struct {
	PullRequestId     string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	IssueId           string `gorm:"primaryKey;type:varchar(255)"`
	PullRequestNumber int
	IssueNumber       int
	common.NoPKModel
}
