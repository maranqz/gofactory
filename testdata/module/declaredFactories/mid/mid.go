package mid

import (
	"factory/declaredFactories/flagged"
	"factory/declaredFactories/other"
	"factory/declaredFactories/owner"
)

func Restore(r other.Repo) owner.ViaMethod {
	return r.Restore()
}

func Flagged() owner.Flagged {
	return flagged.MakeFlagged()
}
