package test

import (
	"testing"

	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/test/factory"
	"github.com/stretchr/testify/assert"
)

func TestInsertWorklog(t *testing.T) {
	board, err := factory.CreateBoard()
	assert.Nil(t, err)
	issue, err := factory.CreateIssue()
	assert.Nil(t, err)
	worklog, err := factory.CreateWorklog(board.Id, issue.Id)
	assert.Nil(t, err)
	tx := models.Db.Create(&worklog)
	assert.Nil(t, tx.Error)
	expected := int64(1)
	assert.Equal(t, tx.RowsAffected, expected)
}
