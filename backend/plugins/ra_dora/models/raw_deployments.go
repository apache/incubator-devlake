package models

type RawDeployments struct {
	RawData string `gorm:"type:jsonb"`
}

func (RawDeployments) TableName() string {
	return "ra_dora_deployments"
}
