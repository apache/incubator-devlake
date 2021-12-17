package test

import (
	"testing"

	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/test/factory"
	"github.com/stretchr/testify/assert"
)

func TestInsertBuild(t *testing.T) {
	job, err := factory.CreateJob()
	assert.Nil(t, err)
	build, err := factory.CreateBuild(job.DomainEntity.Id)
	assert.Nil(t, err)
	tx := models.Db.Create(&build)
	assert.Nil(t, tx.Error)
	expected := int64(1)
	assert.Equal(t, tx.RowsAffected, expected)
}
