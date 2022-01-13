package test

import (
	"testing"

	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/test/factory"
	"github.com/stretchr/testify/assert"
)

func TestInsertChangelog(t *testing.T) {
	_, err := factory.CreateBoard()
	assert.Nil(t, err)
	issue, err := factory.CreateIssue()
	assert.Nil(t, err)
	changelog, err := factory.CreateChangelog(issue.Id)
	assert.Nil(t, err)
	tx := models.Db.Create(&changelog)
	assert.Nil(t, tx.Error)
	expected := int64(1)
	assert.Equal(t, tx.RowsAffected, expected)
}
