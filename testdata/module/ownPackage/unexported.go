package ownPackage

type loan struct{}

func NonProducerUnexportedIsSilent() {
	_ = loan{}
}

// NonProducerLocalTypeIsSilent: LocalResponse is capitalised but declared
// inside a function body, so Go does not export it — only a package-scope
// identifier is exported. No producer could ever name it in a signature,
// so it is unaffected by -ownPackage the same as an unexported type.
func NonProducerLocalTypeIsSilent() {
	type LocalResponse struct{}

	_ = LocalResponse{}
}
