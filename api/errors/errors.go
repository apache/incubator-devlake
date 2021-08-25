package errors

type HttpError struct {
	Status  int
	Message string
}

func (e *HttpError) Code() int {
	return e.Status
}

func (e *HttpError) Error() string {
	return e.Message
}

func NewHttpError(status int, message string) *HttpError {
	return &HttpError{
		status,
		message,
	}
}
