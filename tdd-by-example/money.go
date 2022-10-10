package tddbyexample

type Dollar struct {
	Money
}

type Franc struct {
	Money
}

type Money struct {
	amount int
}

func newDollar(amount int) *Dollar {
	return &Dollar{Money{amount: amount}}
}

func (d *Dollar) times(multiplier int) *Dollar {
	return newDollar(d.amount * multiplier)
}

func (d *Dollar) equals(object interface{}) bool {
	dollar := object.(*Dollar)
	return d.amount == dollar.amount
}

func newFranc(amount int) *Franc {
	return &Franc{Money{amount: amount}}
}

func (f *Franc) times(multiplier int) *Franc {
	return newFranc(f.amount * multiplier)
}

func (f *Franc) equals(object interface{}) bool {
	franc := object.(*Franc)
	return f.amount == franc.amount
}

