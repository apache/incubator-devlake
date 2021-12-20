package test

import (
	"testing"

	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/test/factory"
	"github.com/stretchr/testify/assert"
)

func TestInsertTask(t *testing.T) {
	// TODO: Create pipeline, pass pipelineId to CreateCommit(pipeline.id)
	// Need to solve type problem for domain layer associations using origin key first.
	task, err := factory.CreateTask(1)
	assert.Nil(t, err)
	tx := models.Db.Create(&task)
	assert.Nil(t, tx.Error)
	expected := int64(1)
	assert.Equal(t, tx.RowsAffected, expected)
}
