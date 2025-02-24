package models

type CICDScope struct {
	ID          string `gorm:"primaryKey;size:255;column:id"`
	Name        string `gorm:"size:255;column:name"`
	Description string `gorm:"type:longtext;column:description"`
	URL         string `gorm:"size:255;column:url"`
	CreatedDate string `gorm:"type:datetime(3);column:created_date"`
	UpdatedDate string `gorm:"type:datetime(3);column:updated_date"`
}

func (CICDScope) TableName() string {
	return "cicd_scopes"
}
