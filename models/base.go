package models

import "time"

type Model struct {
	ID        uint64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NoPKModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}
