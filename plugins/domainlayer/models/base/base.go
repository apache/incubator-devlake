package base

import "gorm.io/gorm"

type DomainEntity struct {
	gorm.Model
	OriginKey string `json:"originKey" gorm:"type:varchar(255);uniqueIndex"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
}
