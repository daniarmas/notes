package customerrors

type DuplicateRecord struct {
	Field string
}

func (e *DuplicateRecord) Error() string {
	return e.Field + " already exists"
}
