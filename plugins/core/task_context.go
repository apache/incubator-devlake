package core

import (
	"context"

	"gorm.io/gorm"
)

// prepare for temporal

type SubTaskContext interface {
	GetName() string
	GetConfig(name string) string
	GetDb() *gorm.DB
	GetContext() context.Context
	GetLogger() Logger
	GetData() interface{}
	SetProgress(current int, total int)
	IncProgress(quantity int)
}

// Transient resources needed for task execution
type TaskContext interface {
	SubTaskContext
	SubTaskContext(subtask string) (SubTaskContext, error)
}
