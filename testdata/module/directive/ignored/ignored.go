// Package ignored declares types taken out of protection by
// //gofactory:ignore, so that importing packages see them as unprotected
// too.
package ignored

//gofactory:ignore
type Struct struct{} // want Struct:"gofactory:ignore"

//gofactory:ignore
type Status int // want Status:"gofactory:ignore"

func NewStruct() Struct {
	return Struct{}
}
