package test

import (
	"testing"

	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/test/factory"
	"github.com/stretchr/testify/assert"
)

func TestInsertIssue(t *testing.T) {
	issue, err := factory.CreateIssue()
	assert.Nil(t, err)
	tx := models.Db.Create(&issue)
	assert.Nil(t, tx.Error)
	expected := int64(1)
	assert.Equal(t, tx.RowsAffected, expected)
}
