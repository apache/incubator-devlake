package models

type Deployment struct {
	ID      string `gorm:"primaryKey"`
	ScopeID string `gorm:"index"`
}

func (Deployment) TableName() string {
	return "dora_deployments"
}
