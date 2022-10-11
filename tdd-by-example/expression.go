package tddbyexample

type Expression interface{}

type Sum struct {
	augend IMoney
	addend IMoney
	Expression
}

func NewSum(augend IMoney, addend IMoney) *Sum {
	return &Sum{augend: augend, addend: addend}
}
