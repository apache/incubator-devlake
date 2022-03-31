package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdWorkspace struct {
	SourceId    uint64     `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID          string     `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name        string     `json:"name"`
	PrettyName  string     `json:"pretty_name"`
	Category    string     `json:"category"`
	Status      string     `json:"status"`
	Description string     `json:"description"`
	BeginDate   string     `json:"begin_date"`
	EndDate     string     `json:"end_date"`
	ExternalOn  string     `json:"external_on"`
	Creator     string     `json:"creator"`
	Created     *time.Time `json:"created"`
	common.NoPKModel
}
