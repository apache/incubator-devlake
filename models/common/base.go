package common

import (
	"regexp"
	"time"
)

type Model struct {
	ID        uint64    `gorm:"primaryKey"json:"id"`
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
	RawDataParams string `gorm:"column:_raw_data_params;type:varchar(255);index"`
	RawDataTable  string `gorm:"column:_raw_data_table"`
	// can be used for debugging
	RawDataId uint64 `gorm:"column:_raw_data_id"`
	// we can store record index into this field, which is helpful for debugging
	RawDataRemark string `gorm:"column:_raw_data_remark"`
}

var (
	DUPLICATE_REGEX = regexp.MustCompile(`(?i)\bduplicate\b`)
)

func IsDuplicateError(err error) bool {
	return err != nil && DUPLICATE_REGEX.MatchString(err.Error())
}
