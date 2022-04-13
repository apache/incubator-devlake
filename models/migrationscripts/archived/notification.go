package archived

type NotificationType string

// Notification records notifications sent by lake
type Notification struct {
	Model
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
