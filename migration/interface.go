package migration

import (
	"context"
	"gorm.io/gorm"
)

type Script interface {
	Up(ctx context.Context, db *gorm.DB) error
	Version() uint64
	Name() string
}
