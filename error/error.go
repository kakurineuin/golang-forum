package error

// CustomError 自訂錯誤物件。
type CustomError struct {
	HTTPStatusCode int
	Message        string
}

func (ce CustomError) Error() string {
	return ce.Message
}
