package helper

import (
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Test_getPrimaryKeyValue(t *testing.T) {
	type args struct {
		iface interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"",
			args{&ticket.Sprint{
				DomainEntity: domainlayer.DomainEntity{Id: "abc"},
			},
			},
			"abc",
		},
		{
			"",
			args{ticket.Sprint{
				DomainEntity: domainlayer.DomainEntity{Id: "abc"},
			},
			},
			"abc",
		},
		{
			"",
			args{ticket.SprintIssue{
				SprintId: "abc",
				IssueId:  "123",
			},
			},
			"abc:123",
		},
		{
			"",
			args{&ticket.SprintIssue{
				SprintId: "abc",
				IssueId:  "123",
			},
			},
			"abc:123",
		},
		{
			"",
			args{ticket.Issue{}},
			"",
		},
		{
			"",
			args{&ticket.Issue{}},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, getPrimaryKeyValue(tt.args.iface), "getPrimaryKeyValue(%v)", tt.args.iface)
		})
	}
}

func Test_hasPrimaryKey(t *testing.T) {
	type args struct {
		iface interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"",
			args{ticket.SprintIssue{
				SprintId: "abc",
				IssueId:  "123",
			},
			},
			true,
		},
		{
			"",
			args{&ticket.SprintIssue{
				SprintId: "abc",
				IssueId:  "123",
			},
			},
			true,
		},
		{
			"",
			args{ticket.Issue{}},
			true,
		},
		{
			"",
			args{&ticket.Issue{}},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, hasPrimaryKey(reflect.TypeOf(tt.args.iface)), "hasPrimaryKey(%v)", tt.args.iface)
		})
	}
}

// go test -gcflags=all=-l
func TestBatchSave(t *testing.T) {
	db := &gorm.DB{}
	sqlTimes := 0

	gcl := gomonkey.ApplyMethod(reflect.TypeOf(&gorm.DB{}), "Clauses", func(db *gorm.DB, conds ...clause.Expression) (tx *gorm.DB) {
		sqlTimes++
		return db
	},
	)
	defer gcl.Reset()

	gcr := gomonkey.ApplyMethod(reflect.TypeOf(&gorm.DB{}), "Create", func(db *gorm.DB, value interface{}) *gorm.DB {
		assert.Equal(t, TestTableData, value.([]*TestTable)[0])
		return db
	},
	)
	defer gcr.Reset()

	TestBatchSize = 1
	rowType := reflect.TypeOf(TestTableData)
	batch, err := NewBatchSave(db, rowType, TestBatchSize)

	// test diff type
	assert.Equal(t, err, nil)
	err = batch.Add(&TestBatchSize)
	assert.NotEqual(t, err, nil)

	// test right type
	err = batch.Add(TestTableData)
	assert.Equal(t, err, nil)

	assert.Equal(t, sqlTimes, 1)
}
