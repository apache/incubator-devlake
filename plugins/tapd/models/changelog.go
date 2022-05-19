package models

import (
	"github.com/apache/incubator-devlake/models/common"
	"time"
)

type ChangelogTmp struct {
	Id              uint64
	IssueId         uint64
	AuthorId        string
	AuthorName      string
	FieldId         string
	FieldName       string
	From            string
	To              string
	IterationIdFrom uint64
	IterationIdTo   uint64
	CreatedDate     time.Time
	common.RawDataOrigin
}
