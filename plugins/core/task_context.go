package core

import (
	"context"

	"gorm.io/gorm"
)

// prepare for temporal

// Specifically designed for task execution
type TaskLogger interface {
	Logger
	// update progress for subtask, pass -1 for total if it was unavailable
	Progress(subtask string, current int, total int)
}

// Transient resources needed for task execution
type TaskContext interface {
	GetConfig(name string) string
	GetDb() *gorm.DB
	GetContext() context.Context
	GetData() interface{}
	GetLogger() TaskLogger
}
