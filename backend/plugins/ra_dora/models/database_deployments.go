package models

type DatabaseDeployments struct {
	ID           string `json:"id" gorm:"primaryKey;type:varchar(255);column:id"`
	ScopeID      string `json:"cicd_scope_id" gorm:"type:varchar(255);column:cicd_scope_id"`
	Name         string `json:"name" gorm:"type:varchar(255);column:name"`
	Result       string `json:"result" gorm:"type:varchar(100);column:result"`
	Status       string `json:"status" gorm:"type:varchar(100);column:status"`
	Environment  string `json:"environment" gorm:"type:varchar(255);column:environment"`
	CreatedDate  string `json:"created_date" gorm:"type:datetime(3);column:created_date"`
	StartedDate  string `json:"started_date" gorm:"type:datetime(3);column:started_date"`
	FinishedDate string `json:"finished_date" gorm:"type:datetime(3);column:finished_date"`
	DurationSec  int64  `json:"duration_sec" gorm:"type:bigint;column:duration_sec"`
}

func (DatabaseDeployments) TableName() string {
	return "cicd_deployments"
}
