package tddbyexample

type IMoney interface {
	equals(money IMoney) bool
	Amount() int
	Currency() string
	times(multiplier int) IMoney
	plus(money IMoney) Expression
	Expression
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

func (m *Money) plus(addend IMoney) Expression {
	return NewSum(m, addend)
}

func (m *Money) Amount() int {
	return m.amount
}

func (m *Money) Currency() string {
	return m.currency
}

func (m *Money) reduce(bank *Bank, to string) IMoney {
	rate := bank.rate(m.currency, to)
	return NewMoney(m.amount/rate, to)
}
