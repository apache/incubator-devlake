package archived

import (
	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

type TapdTaskLabel struct {
	TaskId    models.Uint64s `gorm:"primaryKey;autoIncrement:false"`
	LabelName string         `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

func (TapdTaskLabel) TableName() string {
	return "_tool_tapd_task_labels"
}
