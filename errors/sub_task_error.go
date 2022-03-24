package errors

var _ error = (*SubTaskError)(nil)

type SubTaskError struct {
	SubTaskName string
	Message     string
}

func (e *SubTaskError) Error() string {
	return e.Message
}

func (e *SubTaskError) GetSubTaskName() string {
	return e.SubTaskName
}
