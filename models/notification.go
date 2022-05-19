package models

import "github.com/apache/incubator-devlake/models/common"

type NotificationType string

const (
	NotificationPipelineStatusChanged NotificationType = "PipelineStatusChanged"
)

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
