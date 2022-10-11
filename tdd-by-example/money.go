package tddbyexample

type Dollar struct {
	Money
}

type Franc struct {
	Money
}

type IMoney interface {
	equals(money IMoney) bool
	times(multiplier int) IMoney
	Amount() int
	Currency() string
}

type Money struct {
	amount   int
	currency string
}

func newDollar(amount int) IMoney {
	return &Dollar{Money{amount: amount, currency: "USD"}}
}

func newFranc(amount int) IMoney {
	return &Franc{Money{amount: amount, currency: "CHF"}}
}

func newMoney(amount int, currency string) IMoney {
	return &Money{amount: amount, currency: currency}
}

func (m *Money) equals(money IMoney) bool {
	return m.amount == money.Amount() && m.currency == money.Currency()
}

func (m *Money) times(multiplier int) IMoney {
	return newMoney(m.amount*multiplier, m.currency)
}

func (m *Money) Amount() int {
	return m.amount
}

func (m *Money) Currency() string {
	return m.currency
}
