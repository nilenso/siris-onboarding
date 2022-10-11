package tddbyexample

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiplication(t *testing.T) {
	five := newDollar(5)
	assert.True(t, newDollar(10).equals(five.times(2)))
	assert.True(t, newDollar(15).equals(five.times(3)))
}

func TestEquality(t *testing.T) {
	assert.True(t, newDollar(5).equals(newDollar(5)))
	assert.False(t, newDollar(5).equals(newDollar(6)))
	assert.True(t, newFranc(5).equals(newFranc(5)))
	assert.False(t, newFranc(5).equals(newFranc(6)))
	assert.False(t, newDollar(5).equals(newFranc(5)))
}

func TestFrancMultiplication(t *testing.T) {
	five := newFranc(5)
	assert.True(t, newFranc(10).equals(five.times(2)))
	assert.True(t, newFranc(15).equals(five.times(3)))
}

func TestAmount(t *testing.T) {
	five := newDollar(5)
	assert.Equal(t, 5, five.Amount())
}

func TestCurrency(t *testing.T) {
	assert.Equal(t, "USD", newDollar(1).Currency())
	assert.Equal(t, "CHF", newFranc(1).Currency())
}
