package ownPackage

// Order is a value object: PaidIsSilent returns a new Order rather than
// mutating the receiver.
type Order struct {
	Status Status
}

// PaidIsSilent: a method of Order returning Order is a producer the same
// way a top-level function returning Order would be, so this wither may
// build Order directly.
func (o Order) PaidIsSilent() Order {
	return Order{Status: o.Status}
}

// DescribeIsReported: a method of Order whose result does not include
// Order is not a producer.
func (o Order) DescribeIsReported() string {
	_ = Order{} // want `Use factory for ownPackage.Order`

	return "order"
}
