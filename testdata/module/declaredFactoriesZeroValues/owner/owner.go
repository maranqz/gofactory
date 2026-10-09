package owner

// T's only factory is other.New, a //gofactory:factory function in another
// package.
type T struct{ X int }
