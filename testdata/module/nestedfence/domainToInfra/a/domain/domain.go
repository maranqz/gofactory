// Package domain lies in both the outer fence (factory/nestedfence/*/a/**)
// and the inner fence scoped to itself (factory/nestedfence/*/a/domain/**),
// so it can still bypass infra's factory: infra lies only in the outer
// fence, and domain satisfies that fence too.
package domain

import "factory/nestedfence/domainToInfra/a/infra"

type Struct struct{}

func New() Struct {
	return Struct{}
}

func CreateInfra() {
	_ = infra.Struct{}
	_ = &infra.Struct{}
}
