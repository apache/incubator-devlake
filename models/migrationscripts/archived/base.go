package archived

import (
	"time"
)

type DomainEntity struct {
	Id string `json:"id" gorm:"primaryKey;type:varchar(255);comment:This key is generated based on details from the original plugin"` // format: <Plugin>:<Entity>:<PK0>:<PK1>
	NoPKModel
}

type Model struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NoPKModel struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	RawDataOrigin
}

// embedded fields for tool layer tables
type RawDataOrigin struct {
	// can be used for flushing outdated records from table
	RawDataParams string `gorm:"column:_raw_data_params;type:varchar(255);index" json:"_raw_data_params"`
	RawDataTable  string `gorm:"column:_raw_data_table;type:varchar(255)" json:"_raw_data_table"`
	// can be used for debugging
	RawDataId uint64 `gorm:"column:_raw_data_id" json:"_raw_data_id"`
	// we can store record index into this field, which is helpful for debugging
	RawDataRemark string `gorm:"column:_raw_data_remark" json:"_raw_data_remark"`
}
