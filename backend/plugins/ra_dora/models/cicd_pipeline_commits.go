package models

type CICDPipelineCommit struct {
	PipelineID string `gorm:"size:255;column:pipeline_id"`
	CommitSHA  string `gorm:"size:40;column:commit_sha"`
	Branch     string `gorm:"size:255;column:branch"`
	Repo       string `gorm:"size:255;column:repo"`
	RepoID     string `gorm:"size:255;column:repo_id"`
}

func (CICDPipelineCommit) TableName() string {
	return "cicd_pipeline_commits"
}
