package mid

import (
	"factory/unimplemented/visibility/owner"
	"factory/unimplemented/visibility/wrapper"
)

func Mid() owner.T {
	return wrapper.New()
}
