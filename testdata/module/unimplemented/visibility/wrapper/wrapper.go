package wrapper

import "factory/unimplemented/visibility/owner"

//gofactory:factory
func New() owner.T {
	return owner.T{}
}
