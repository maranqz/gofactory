// Package owner declares a type taken out of protection by
// //gofactory:ignore and bypasses it locally, so the crossPackageDirectives
// test case can check that -crossPackageDirectives=false, which disables
// cross-package propagation, still leaves a same-package directive working.
package owner

//gofactory:ignore
type Ignored struct{}

func NewIgnored() Ignored {
	return Ignored{}
}

func Local() {
	_ = Ignored{}
	_ = &Ignored{}
	_ = new(Ignored)
}
