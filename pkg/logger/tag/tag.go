package tag

const (
	// RequestIDKey is msg_id for request identifier
	RequestIDKey = "request_id"
)

// Tag is key value pair with value in string
type Tag struct {
	Key   string
	Value string
}

// Err to print tag error
func Err(err error) Tag {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	return Tag{
		Key:   "error",
		Value: errMsg,
	}
}
