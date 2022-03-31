package models

type TapdUser struct {
	WorkspaceId uint64 `gorm:"primaryKey"`
	Name        string `json:"user" gorm:"primaryKey"`
}
