package tddbyexample

type IMoney interface {
	equals(money IMoney) bool
	Amount() int
	Currency() string
	times(multiplier int) IMoney
	plus(addend IMoney) IMoney
}

type Money struct {
	amount   int
	currency string
}

func dollar(amount int) IMoney {
	return NewMoney(amount, "USD")
}

func franc(amount int) IMoney {
	return NewMoney(amount, "CHF")
}

func NewMoney(amount int, currency string) IMoney {
	return &Money{amount: amount, currency: currency}
}

func (m *Money) equals(money IMoney) bool {
	return m.amount == money.Amount() && m.currency == money.Currency()
}

func (m *Money) times(multiplier int) IMoney {
	return NewMoney(m.amount*multiplier, m.currency)
}

func (m *Money) plus(addend IMoney) IMoney {
	return NewMoney(m.amount+addend.Amount(), m.currency)
}

func (m *Money) Amount() int {
	return m.amount
}

func (m *Money) Currency() string {
	return m.currency
}
