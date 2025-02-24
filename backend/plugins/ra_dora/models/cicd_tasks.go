package models

type CICDTask struct {
	ID           string `gorm:"primaryKey;size:255;column:id"`
	Name         string `gorm:"size:255;column:name"`
	PipelineID   string `gorm:"size:255;column:pipeline_id"`
	Result       string `gorm:"size:100;column:result"`
	Status       string `gorm:"size:100;column:status"`
	Type         string `gorm:"size:100;column:type"`
	DurationSec  int64  `gorm:"type:bigint unsigned;column:duration_sec"`
	StartedDate  string `gorm:"type:datetime(3);column:started_date"`
	FinishedDate string `gorm:"type:datetime(3);column:finished_date"`
	Environment  string `gorm:"size:255;column:environment"`
	ScopeID      string `gorm:"type:longtext;column:cicd_scope_id"`
}

func (CICDTask) TableName() string {
	return "cicd_tasks"
}
