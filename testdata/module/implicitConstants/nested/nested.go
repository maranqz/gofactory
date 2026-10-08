package nested

type Status int

const (
	Active Status = iota + 1
	Inactive
)

const Closed Status = 3

// The owner package may store an untyped constant in its own type.
func NewStatus() Status {
	return 1
}

var Default Status = 1
