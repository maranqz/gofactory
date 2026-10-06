// Package infra lies in the outer fence (factory/nestedfence/*/a/**) but
// not the inner fence scoped to domain (factory/nestedfence/*/a/domain/**),
// so it cannot bypass domain's factory even though both packages sit under
// the same outer fence.
package infra

import "factory/nestedfence/infraToDomain/a/domain"

type Struct struct{}

func New() Struct {
	return Struct{}
}

func CreateDomain() {
	_ = domain.Struct{}  // want `Use factory for domain.Struct \(domain.New\)`
	_ = &domain.Struct{} // want `Use factory for domain.Struct \(domain.New\)`
}
