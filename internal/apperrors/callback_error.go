package apperrors

type CallbackError struct {
	Err     error
	Message string
}

func (c CallbackError) Error() string {
	return c.Err.Error()
}

func NewCallbackError(err error, message string) *CallbackError {
	return &CallbackError{Err: err, Message: message}
}
