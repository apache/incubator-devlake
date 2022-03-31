package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdWorkspace struct {
	ID          uint64     `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name"`
	PrettyName  uint64     `json:"pretty_name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	ParentID    uint64     `json:"parent_id"`
	Secrecy     string     `json:"secrecy"`
	Created     string     `json:"created"`
	CreatorID   uint64     `json:"creator_id"`
	Creator     string     `json:"creator"`
	BeginDate   *time.Time `json:"begin_date"`
	EndDate     *time.Time `json:"end_date"`
	MemberCount uint64     `json:"member_count"`
	common.NoPKModel
}
