package models

import (
	"github.com/merico-dev/lake/models/common"
)

type TapdIteration struct {
	SourceId     uint64 `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID           string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name         string `json:"name"`
	WorkspaceID  string `json:"workspace_id"`
	Startdate    string `json:"startdate"`
	Enddate      string `json:"enddate"`
	Status       string `json:"status"`
	ReleaseID    string `json:"release_id"`
	Description  string `json:"description"`
	Creator      string `json:"creator"`
	Created      string `json:"created"`
	Modified     string `json:"modified"`
	Completed    string `json:"completed"`
	Releaseowner string `json:"releaseowner"`
	Launchdate   string `json:"launchdate"`
	Notice       string `json:"notice"`
	Releasename  string `json:"releasename"`
	common.NoPKModel
}
