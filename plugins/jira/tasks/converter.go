package tasks

import (
	"context"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/gorm"
)

type JiraConverter interface {
	Convert(ctx context.Context, db *gorm.DB, logger core.Logger, args models.Args) error
}

type ConverterSubTask struct {
	ctx       core.ExecContext
	converter JiraConverter
}

func (t *ConverterSubTask) Execute() error {
	ctx := t.ctx.GetContext()
	db := t.ctx.GetDb()
	logger := t.ctx.GetLogger()
	args := t.ctx.GetData().(models.Args)
	return t.converter.Convert(ctx, db, logger, args)
}

func NewConverterSubTask(ctx core.ExecContext, converter JiraConverter) core.SubTask {
	return &ConverterSubTask{
		ctx:       ctx,
		converter: converter,
	}
}
