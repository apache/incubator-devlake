package models

type NotificationType string

const (
	NotificationTaskSuccess NotificationType = "TaskSuccess"
)

// Notification records notifications sent by lake
type Notification struct {
	Model
	Type         NotificationType
	Nonce        string
	ResponseCode int
	Response     string
	Data         string
}
