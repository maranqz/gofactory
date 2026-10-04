package nested

// Struct ends up with no recognised factory: every New-named candidate
// below is disqualified by a different recognition rule.
type Struct struct{}

// Paid is a wither: a method of Struct itself is never a factory, however
// its name and signature look.
func (s Struct) Paid() Struct {
	return s
}

// NewFromStruct takes a Struct parameter, so it is never a factory.
func NewFromStruct(_ Struct) Struct {
	return Struct{}
}

// NewFromStructPtr takes a *Struct parameter, so it is never a factory.
func NewFromStructPtr(_ *Struct) *Struct {
	return &Struct{}
}

// newUnexported is unexported, so it is never a factory.
func newUnexported() Struct {
	return Struct{}
}

// Make does not match the default ^New pattern.
func Make() Struct {
	return Struct{}
}

// Widget has a single recognised factory, a method of another type:
// Builder, not Widget itself.
type Widget struct{}

type Builder struct{}

// NewWidget is a recognised factory method of Builder.
func (Builder) NewWidget() Widget {
	return Widget{}
}

// WithErr has a recognised factory that also returns an error; the error
// result is ignored, both for "no target parameter" and for "returns
// target", the same way a declared factory's results ignore error.
type WithErr struct{}

func NewWithErr() (*WithErr, error) {
	return &WithErr{}, nil
}

// Secret has one recognised factory, but it is a method of hidden, an
// unexported type: outside this package nobody can write nested.hidden, so
// the factory is never accessible from another package's diagnostic.
type Secret struct{}

type hidden struct{}

// NewSecret is a recognised factory of Secret, accessible only from within
// this package.
func (hidden) NewSecret() Secret {
	return Secret{}
}

// Multi has four recognised factories, to exercise the three-factory cap
// and the deterministic (alphabetical) tie-break among equally-ranked
// names.
type Multi struct{}

func NewA() Multi {
	return Multi{}
}

func NewB() Multi {
	return Multi{}
}

func NewC() Multi {
	return Multi{}
}

func NewD() *Multi {
	return &Multi{}
}
