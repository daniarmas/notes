package customerrors

// Unknown represents an error when something unknown happend.
type Unknown struct {
}

func (e *Unknown) Error() string {
	return "unknown error"
}
