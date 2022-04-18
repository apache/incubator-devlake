package code

import "github.com/merico-dev/lake/models/common"

type RefsPrCherrypick struct {
	RepoName               string `gorm:"type:varchar(255)"`
	ParentPrKey            int
	CherrypickBaseBranches string `gorm:"type:varchar(255)"`
	CherrypickPrKeys       string `gorm:"type:varchar(255)"`
	ParentPrUrl            string `gorm:"type:varchar(255)"`
	ParentPrId             string `json:"parent_pr_id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	common.NoPKModel
}
