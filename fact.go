package gofactory

import (
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// exportFact skips the export while Analyzer.FactTypes is empty: go vet's gob
// encoder registers only declared fact types and panics on any other.
// Importing a fact is a plain lookup in every driver and needs no such check.
func exportFact(pass *analysis.Pass, obj types.Object, fact analysis.Fact) {
	if len(pass.Analyzer.FactTypes) > 0 {
		pass.ExportObjectFact(obj, fact)
	}
}

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
