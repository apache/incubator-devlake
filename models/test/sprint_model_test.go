package test

import (
	"testing"

	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/test/factory"
	"github.com/stretchr/testify/assert"
)

func TestInsertSprint(t *testing.T) {
	board, err := factory.CreateBoard()
	assert.Nil(t, err)
	sprint, err := factory.CreateSprint(board.DomainEntity.Id)
	assert.Nil(t, err)
	tx := models.Db.Create(&sprint)
	assert.Nil(t, tx.Error)
	expected := int64(1)
	assert.Equal(t, tx.RowsAffected, expected)
}
