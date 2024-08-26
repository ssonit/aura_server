package common

type CustomError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Code       string `json:"code"`
}

func (e *CustomError) Error() string {
	return e.Message
}

func (e *CustomError) GetStatusCode() int {
	return e.StatusCode
}

func NewFullCustomError(statusCode int, message, code string) *CustomError {
	return &CustomError{
		StatusCode: statusCode,
		Message:    message,
		Code:       code,
	}
}
