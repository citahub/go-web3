package errors

const (
	invalidRequestCode = -32600
	methodNotFoundCode = -32601
	invalidParamsCode  = -32602
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func IsInvalidRequest(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == invalidRequestCode
}

func IsMethodNotFound(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == methodNotFoundCode
}

func IsInvalidParams(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == invalidParamsCode
}

func IsNull(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.Code == 0
}

func (e *Error) Error() string {
	return e.Message
}
