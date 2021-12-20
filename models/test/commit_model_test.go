package test

import (
	"testing"

	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/test/factory"
	"github.com/stretchr/testify/assert"
)

func TestInsertCommit(t *testing.T) {
	// TODO: Create repo, pass repoId to CreateCommit(repo.id)
	// Need to solve type problem for domain layer associations using origin key first.
	commit, err := factory.CreateCommit(1)
	assert.Nil(t, err)
	tx := models.Db.Create(&commit)
	assert.Nil(t, tx.Error)
	expected := int64(1)
	assert.Equal(t, tx.RowsAffected, expected)
}
