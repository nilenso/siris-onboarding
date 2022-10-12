package tddbyexample

type Expression interface {
	reduce(to string) IMoney
}

type Sum struct {
	augend IMoney
	addend IMoney
	Expression
}

func NewSum(augend IMoney, addend IMoney) *Sum {
	return &Sum{augend: augend, addend: addend}
}

func (s *Sum) reduce(to string) IMoney {
	amount := s.augend.Amount() + s.addend.Amount()
	return NewMoney(amount, to)
}
