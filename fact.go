package gofactory

// ignoredFact marks a *types.TypeName as taken out of protection by a valid
// //gofactory:ignore directive. Exporting and importing it as an
// analysis.Fact is what makes the directive take effect in every package
// and module that imports the type, not just the one that declares it.
type ignoredFact struct{}

func (*ignoredFact) AFact() {}

func (*ignoredFact) String() string {
	return "gofactory:ignore"
}
