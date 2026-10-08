package mid

import (
	"factory/declaredFactories/other"
	"factory/declaredFactories/owner"
)

func Restore(r other.Repo) owner.ViaMethod {
	return r.Restore()
}
