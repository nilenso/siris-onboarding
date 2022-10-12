package tddbyexample

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiplication(t *testing.T) {
	five := dollar(5)
	assert.True(t, dollar(10).equals(five.times(2)))
	assert.True(t, dollar(15).equals(five.times(3)))
}

func TestSimpleAddition(t *testing.T) {
	sum := dollar(5).plus(dollar(5))
	bank := new(Bank)
	reduced := bank.reduce(sum, "USD")
	assert.Equal(t, dollar(10), reduced)
}

func TestEquality(t *testing.T) {
	assert.True(t, dollar(5).equals(dollar(5)))
	assert.False(t, dollar(5).equals(dollar(6)))
	assert.False(t, dollar(5).equals(franc(5)))
}

func TestAmount(t *testing.T) {
	assert.Equal(t, 5, dollar(5).Amount())
}

func TestCurrency(t *testing.T) {
	assert.Equal(t, "USD", dollar(1).Currency())
}

func TestPlusReturnsSum(t *testing.T) {
	five := dollar(5)
	result := five.plus(five)
	sum := result.(*Sum)
	assert.Equal(t, five, sum.augend)
	assert.Equal(t, five, sum.addend)
}

func TestReduceSum(t *testing.T) {
	sum := NewSum(dollar(3), dollar(4))
	bank := new(Bank)
	result := bank.reduce(sum, "USD")
	assert.Equal(t, dollar(7), result)
}

func TestReduceMoney(t *testing.T) {
	bank := new(Bank)
	result := bank.reduce(dollar(1), "USD")
	assert.Equal(t, dollar(1), result)
}
