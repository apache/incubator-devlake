package common

import (
	"regexp"
	"time"
)

type Model struct {
	ID        uint64 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type NoPKModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	DUPLICATE_REGEX = regexp.MustCompile(`(?i)\bduplicate\b`)
)

func IsDuplicateError(err error) bool {
	return err != nil && DUPLICATE_REGEX.MatchString(err.Error())
}
