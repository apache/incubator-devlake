package models

import (
	"time"
)

type FeishuMeetingTopUserItem struct {
	Model           `json:"-"`
	StartTime       time.Time
	MeetingCount    string `json:"meeting_count"`
	MeetingDuration string `json:"meeting_duration"`
	Name            string `json:"name"`
	UserType        int64  `json:"user_type"`
}
