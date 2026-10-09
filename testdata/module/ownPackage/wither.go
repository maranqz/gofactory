package ownPackage

type Order struct {
	Status Status
}

func (o Order) PaidIsSilent() Order {
	return Order{Status: o.Status}
}

func (o Order) DescribeIsReported() string {
	_ = Order{} // want `Use factory for ownPackage.Order`

	return "order"
}
