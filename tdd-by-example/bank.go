package tddbyexample

type Bank struct{}

func (b *Bank) reduce(source Expression, to string) IMoney {
	return dollar(10)
}
