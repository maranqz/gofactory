package wrapper

import "factory/unimplemented/visibility/owner"

//gofactory:factory
func New() owner.T { // want New:"gofactory:factory"
	return owner.T{}
}
