package test

import (
	"testing"

	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/test/factory"
	"github.com/stretchr/testify/assert"
)

func TestInsertNote(t *testing.T) {
	// TODO: Create pr, pass prId to CreateCommit(pr.id)
	// Need to solve type problem for domain layer associations using origin key first.
	note, err := factory.CreateNote(1)
	assert.Nil(t, err)
	tx := models.Db.Create(note)
	assert.Nil(t, tx.Error)
	expected := int64(1)
	assert.Equal(t, tx.RowsAffected, expected)
}
