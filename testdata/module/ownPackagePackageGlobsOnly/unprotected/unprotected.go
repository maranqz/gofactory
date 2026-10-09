package unprotected

// Request is outside every -packageGlobs fence, so -packageGlobsOnly
// leaves it unprotected; -ownPackage must not restrict a type that isn't
// protected at all.
type Request struct{}

func NonProducerUnfencedTypeIsSilent() {
	_ = Request{}
}
