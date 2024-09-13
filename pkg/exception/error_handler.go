package exception

type Error struct {
	Status  string
	Message string
	Errors  interface{}
}

func ErrorF(status string, message string, err interface{}) Error {

	errorData := Error{
		Status:  status,
		Message: message,
		Errors:  err,
	}

	return errorData
}
