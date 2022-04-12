package archived

import "github.com/merico-dev/lake/models/common"

type NotificationType string

// Notification records notifications sent by lake
type Notification struct {
	common.Model
	Type         NotificationType
	Endpoint     string
	Nonce        string
	ResponseCode int
	Response     string
	Data         string
}

func (Notification) TableName() string {
	return "_devlake_notifications"
}
