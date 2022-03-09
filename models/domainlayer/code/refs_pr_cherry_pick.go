package code

import "github.com/merico-dev/lake/models/common"

type RefsPrCherrypick struct {
	RepoName               string
	ParentPrKey            int
	CherrypickBaseBranches string
	CherrypickPrKeys       string
	ParentPrUrl            string
	ParentPrId             string `json:"parent_pr_id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	common.NoPKModel
}
