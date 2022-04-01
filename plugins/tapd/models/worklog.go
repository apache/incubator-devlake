package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdWorklog struct {
	SourceId    uint64     `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID          uint64     `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	WorkspaceId uint64     `json:"workspace_id"`
	EntityType  string     `json:"entity_type"`
	EntityID    uint64     `json:"entity_id"`
	Timespent   int        `json:"timespent"`
	Spentdate   *time.Time `json:"spentdate"`
	Owner       string     `json:"owner"`
	Created     *time.Time `json:"created"`
	Memo        string     `json:"memo"`
	common.NoPKModel
}
type TapdWorklogApiRes struct {
	ID          string `json:"id"`
	EntityType  string `json:"entity_type"`
	EntityID    string `json:"entity_id"`
	Timespent   string `json:"timespent"`
	Spentdate   string `json:"spentdate"`
	Owner       string `json:"owner"`
	Created     string `json:"created"`
	WorkspaceId string `json:"workspace_id"`
	Memo        string `json:"memo"`
}
