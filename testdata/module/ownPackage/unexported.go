package ownPackage

// loan is unexported: -ownPackage restricts only exported protected types,
// so loan stays as free to build as it always was, producer or not.
type loan struct{}

func NonProducerUnexportedIsSilent() {
	_ = loan{}
}
