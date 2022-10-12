package tddbyexample

type Bank struct {
	rates map[pair]int
}

type pair struct {
	from string
	to   string
}

func NewBank() *Bank {
	return &Bank{rates: make(map[pair]int)}
}
func (b *Bank) reduce(source Expression, to string) IMoney {
	return source.reduce(b, to)
}

func (b *Bank) addRate(from string, to string, rate int) {
	b.rates[pair{
		from: from,
		to:   to,
	}] = rate
}

func (b *Bank) rate(from string, to string) int {
	if from == to {
		return 1
	}
	rate := b.rates[pair{from: from, to: to}]
	return rate
}
