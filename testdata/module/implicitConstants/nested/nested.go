package nested

type Status int

const (
	Active Status = iota + 1
	Inactive
)

// The owner package may store untyped constants in its own type.
const Closed Status = 3

func NewStatus() Status {
	return 1
}

var Default Status = 1
