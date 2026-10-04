package nested

// Struct ends up with no recognised factory: every New-named candidate
// below is disqualified by a different recognition rule.
type Struct struct{}

// Paid is a wither: a method of Struct itself is never a factory, however
// its name and signature look.
func (s Struct) Paid() Struct {
	return s
}

// NewCopy and NewClone are withers too, but named like a factory: a method
// of Struct itself is never its factory even when the name matches ^New,
// so deleting the "methods of target are never factories" rule would make
// one of them appear in the suffix at factories.go:11.
func (s Struct) NewCopy() Struct {
	return s
}

func (s *Struct) NewClone() *Struct {
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

// StructDefPtr is a defined pointer type to Struct: taking one as a
// parameter is the same as taking a *Struct.
type StructDefPtr *Struct

// NewFromStructDefPtr takes a StructDefPtr parameter, so it is never a
// factory, just like NewFromStructPtr above.
func NewFromStructDefPtr(_ StructDefPtr) Struct {
	return Struct{}
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

// WithErr has a recognised factory that also returns an error: the rule is
// "target or *target among its results", and an extra error result is just
// another result that isn't target, so it doesn't disqualify the function.
type WithErr struct{}

func NewWithErr() (*WithErr, error) {
	return &WithErr{}, nil
}

// Sorted has two recognised factory methods declared out of source order,
// to pin the alphabetical tie-break: package-level names are already
// listed alphabetically by the type checker, so only methods, which come
// out in source order, can tell the sort apart from no sort at all.
type Sorted struct{}

type Maker struct{}

func (Maker) NewSortedB() Sorted {
	return Sorted{}
}

func (Maker) NewSortedA() Sorted {
	return Sorted{}
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
