package models

import (
	"github.com/apache/incubator-devlake/models/common"
	"time"
)

type FeishuMeetingTopUserItem struct {
	common.Model    `json:"-"`
	StartTime       time.Time
	MeetingCount    string `json:"meeting_count" gorm:"type:varchar(255)"`
	MeetingDuration string `json:"meeting_duration" gorm:"type:varchar(255)"`
	Name            string `json:"name" gorm:"type:varchar(255)"`
	UserType        int64  `json:"user_type"`
	common.RawDataOrigin
}

func (FeishuMeetingTopUserItem) TableName() string {
	return "_tool_feishu_meeting_top_user_items"
}
