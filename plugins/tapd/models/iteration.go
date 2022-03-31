package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type TapdIteration struct {
	SourceId     uint64     `gorm:"primaryKey;type:INT(10) UNSIGNED NOT NULL"`
	ID           uint64     `gorm:"primaryKey;type:BIGINT(10) UNSIGNED NOT NULL" json:"id"`
	Name         string     `json:"name"`
	WorkspaceID  uint64     `json:"workspace_id"`
	Startdate    *time.Time `json:"startdate"`
	Enddate      *time.Time `json:"enddate"`
	Status       string     `json:"status"`
	ReleaseID    string     `json:"release_id"`
	Description  string     `json:"description"`
	Creator      string     `json:"creator"`
	Created      *time.Time `json:"created"`
	Modified     *time.Time `json:"modified"`
	Completed    *time.Time `json:"completed"`
	Releaseowner string     `json:"releaseowner"`
	Launchdate   *time.Time `json:"launchdate"`
	Notice       string     `json:"notice"`
	Releasename  string     `json:"releasename"`
	common.NoPKModel
}
type TapdIterationRes struct {
	ID           string `json:"id"`
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
}
