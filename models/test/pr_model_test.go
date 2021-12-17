package test

import (
	"testing"

	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/test/factory"
	"github.com/stretchr/testify/assert"
)

func TestInsertPr(t *testing.T) {
	// TODO: Create repo, pass repoId to CreateCommit(repo.id)
	// Need to solve type problem for domain layer associations using origin key first.
	pr, err := factory.CreatePr(1)
	assert.Nil(t, err)
	tx := models.Db.Create(&pr)
	assert.Nil(t, tx.Error)
	expected := int64(1)
	assert.Equal(t, tx.RowsAffected, expected)
}
