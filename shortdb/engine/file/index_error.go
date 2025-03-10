package file

// CreateExistIndexError is an error type returned when the index already exists
type CreateExistIndexError struct {
	Name string
}

func (e *CreateExistIndexError) Error() string {
	return "at CREATE INDEX: exist index " + e.Name
}
