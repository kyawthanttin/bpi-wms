package validation

type ErrorWithHttpStatus struct {
	statusCode int
	message    string
}

func NewErrorWithHttpStatus(statusCode int, message string) *ErrorWithHttpStatus {
	return &ErrorWithHttpStatus{statusCode, message}
}

func (e *ErrorWithHttpStatus) Status() int {
	return e.statusCode
}

func (e *ErrorWithHttpStatus) Error() string {
	return e.message
}
