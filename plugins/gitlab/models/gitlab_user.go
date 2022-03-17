package models

type GitlabUser struct {
	Email string `gorm:"primaryKey;type:varchar(255)"`
	Name  string `gorm:"type:varchar(255)"`
}
