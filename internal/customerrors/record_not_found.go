package customerrors

// RecordNotFound represents an error when a record is not found.
type RecordNotFound struct {
}

// Error implements the error interface for RecordNotFound.
func (e *RecordNotFound) Error() string {
	return "record not found"
}
