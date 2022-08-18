package devops

import "github.com/apache/incubator-devlake/models/domainlayer"

type CiCDPipelineRepo struct {
	domainlayer.DomainEntity
	CommitSha string `gorm:"primaryKey;type:varchar(255)"`
	Branch    string `gorm:"type:varchar(255)"`
	RepoUrl   string `gorm:"type:varchar(255)"`
}

func (CiCDPipelineRepo) TableName() string {
	return "cicd_pipeline_repos"
}
