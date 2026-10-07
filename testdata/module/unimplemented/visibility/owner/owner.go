package owner

// T has no factory recognised in this package; its only factory is
// wrapper.New, declared in a package indirect reaches only through mid
// (see ../indirect/indirect.go).
type T struct{}
