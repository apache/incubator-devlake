package base

import (
	"github.com/merico-dev/lake/models"
)

type DomainEntity struct {
	OriginKey string `json:"originKey" gorm:"primaryKey;type:varchar(255)"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	models.NoPKModel
}
