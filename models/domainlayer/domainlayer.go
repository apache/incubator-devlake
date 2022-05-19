package domainlayer

import (
	"github.com/apache/incubator-devlake/models/common"
)

type DomainEntity struct {
	Id string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	common.NoPKModel
}
