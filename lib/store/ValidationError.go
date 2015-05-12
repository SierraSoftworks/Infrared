package store

type ValidationError struct {
	message string
}

func (e ValidationError) Error() string {
	return e.message
}
