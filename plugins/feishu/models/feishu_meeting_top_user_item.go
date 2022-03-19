package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type FeishuMeetingTopUserItem struct {
	common.Model           `json:"-"`
	StartTime       time.Time
	MeetingCount    string `json:"meeting_count"`
	MeetingDuration string `json:"meeting_duration"`
	Name            string `json:"name"`
	UserType        int64  `json:"user_type"`
	common.RawDataOrigin
}
