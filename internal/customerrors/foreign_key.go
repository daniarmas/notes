package customerrors

type ForeignKeyConstraint struct {
	Field       string
	ParentTable string
}

func (e *ForeignKeyConstraint) Error() string {
	return "foreign key constraint violation on field: " + e.Field + " referencing table: " + e.ParentTable
}
