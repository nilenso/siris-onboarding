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
	assert.Equal(t, dollar(10), sum)
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
