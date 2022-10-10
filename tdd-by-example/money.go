package tddbyexample

type Dollar struct {
	amount int
}

func newDollar(amount int) *Dollar {
	return &Dollar{amount: amount}
}

func (d *Dollar) times(multiplier int) *Dollar {
	return newDollar(d.amount * multiplier)
}

func (d *Dollar) equals(object interface{}) bool {
	return true
}
