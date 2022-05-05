package models

import "time"

type IssueComment struct {
	ConnectionId       uint64 `gorm:"primaryKey"`
	IssueId            uint64 `gorm:"primarykey"`
	ComentId           string `gorm:"primarykey"`
	Self               string `gorm:"type:varchar(255)"`
	Body               string
	CreatorAccountId   string     `gorm:"type:varchar(255)"`
	CreatorDisplayName string     `gorm:"type:varchar(255)"`
	Created            *time.Time `json:"created"`
	Updated            *time.Time `json:"updated"`
}
