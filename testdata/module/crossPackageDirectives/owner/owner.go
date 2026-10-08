// Package owner is the crossPackageDirectivesOff case's same-package side.
// Ignored carries the //gofactory:ignore that, with the setting false, no
// longer reaches the importer; GlobIgnored has no directive and is exempted
// only by the case's -ignoreTypes setting; Misplaced puts the directive
// where it does not belong, to show directive validation still runs with
// the setting false.
package owner

//gofactory:ignore
type Ignored struct{}

func NewIgnored() Ignored {
	return Ignored{}
}

type GlobIgnored struct{}

//gofactory:ignore // want `//gofactory:ignore must be in the doc comment of a single top-level type definition`
func Misplaced() {}
