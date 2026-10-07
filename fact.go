package gofactory

// ignoredFact marks a *types.TypeName as taken out of protection by a valid
// //gofactory:ignore directive.
type ignoredFact struct{}

func (*ignoredFact) AFact() {}

func (*ignoredFact) String() string {
	return "gofactory:ignore"
}

// factoryFact marks a *types.Func as a declared factory, exported by a
// valid //gofactory:factory directive or a -factories glob match. It
// carries no data: the protected types it is a factory of are the
// protected types among its own results, recomputed from its signature
// wherever needed instead of stored here.
type factoryFact struct{}

func (*factoryFact) AFact() {}

func (*factoryFact) String() string {
	return "gofactory:factory"
}
