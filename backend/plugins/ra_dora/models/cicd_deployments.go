package models

type CICDDeployment struct {
	ID           string `gorm:"primaryKey;size:255;column:id"`
	ScopeID      string `gorm:"size:255;column:cicd_scope_id"`
	Name         string `gorm:"size:255;column:name"`
	Result       string `gorm:"size:100;column:result"`
	Status       string `gorm:"size:100;column:status"`
	Environment  string `gorm:"size:255;column:environment"`
	CreatedDate  string `gorm:"type:datetime(3);column:created_date"`
	StartedDate  string `gorm:"type:datetime(3);column:started_date"`
	FinishedDate string `gorm:"type:datetime(3);column:finished_date"`
	DurationSec  int64  `gorm:"type:bigint;column:duration_sec"`
}

func (CICDDeployment) TableName() string {
	return "cicd_deployments"
}
