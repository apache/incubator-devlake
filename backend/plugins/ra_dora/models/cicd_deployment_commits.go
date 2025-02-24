package models

type CICDDeploymentCommit struct {
	ID           string `gorm:"primaryKey;size:255;column:id"`
	ScopeID      string `gorm:"size:255;column:cicd_scope_id"`
	DeploymentID string `gorm:"size:255;column:cicd_deployment_id"`
	Name         string `gorm:"size:255;column:name"`
	Result       string `gorm:"size:100;column:result"`
	Status       string `gorm:"size:100;column:status"`
	Environment  string `gorm:"size:255;column:environment"`
	CreatedDate  string `gorm:"type:datetime(3);column:created_date"`
	StartedDate  string `gorm:"type:datetime(3);column:started_date"`
	FinishedDate string `gorm:"type:datetime(3);column:finished_date"`
	DurationSec  int64  `gorm:"type:bigint;column:duration_sec"`
	CommitSHA    string `gorm:"size:40;column:commit_sha"`
	RefName      string `gorm:"size:255;column:ref_name"`
	RepoID       string `gorm:"size:255;column:repo_id"`
	RepoURL      string `gorm:"size:255;column:repo_url"`
}

func (CICDDeploymentCommit) TableName() string {
	return "cicd_deployment_commits"
}
