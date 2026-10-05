package gofactory

// ignoredFact marks a *types.TypeName as taken out of protection by a valid
// //gofactory:ignore directive.
type ignoredFact struct{}

func (*ignoredFact) AFact() {}

func (*ignoredFact) String() string {
	return "gofactory:ignore"
}
