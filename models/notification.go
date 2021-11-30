package models

import "github.com/merico-dev/lake/models/common"

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
