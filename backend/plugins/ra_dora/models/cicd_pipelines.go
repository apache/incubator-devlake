package models

type CICDPipeline struct {
	ID           string `gorm:"primaryKey;size:255;column:id"`
	Name         string `gorm:"size:255;column:name"`
	Result       string `gorm:"size:100;column:result"`
	Status       string `gorm:"size:100;column:status"`
	Type         string `gorm:"size:100;column:type"`
	DurationSec  int64  `gorm:"type:bigint unsigned;column:duration_sec"`
	CreatedDate  string `gorm:"type:datetime(3);column:created_date"`
	FinishedDate string `gorm:"type:datetime(3);column:finished_date"`
	Environment  string `gorm:"size:255;column:environment"`
	ScopeID      string `gorm:"type:longtext;column:cicd_scope_id"`
}

func (CICDPipeline) TableName() string {
	return "cicd_pipelines"
}
