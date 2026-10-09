package implicitConstants

import (
	"time"

	"factory/implicitConstants/nested"
)

func Comparison(st nested.Status) bool {
	return st == 3
}

func Case(st nested.Status) {
	switch st {
	case 3:
	}
}

func Arithmetic(st nested.Status) {
	st = st + 1
	st += 1
	st++
	st = nested.Active + 1
	st = nested.Active << 1
	st = max(nested.Active, 2)
	st = -nested.Active

	_ = st
}

const Const nested.Status = 3

func TypedConstant() nested.Status {
	var st nested.Status = nested.Active

	store(nested.Closed)

	st = Const

	_ = st

	return nested.Inactive
}

func Stdlib() time.Duration {
	var d time.Duration = 5

	return d * 2
}

func Interface() any {
	var v any = 3

	return v
}

func Builtins(st nested.Status) {
	_ = min(st, 3)
	_ = max(st, 3)
}
