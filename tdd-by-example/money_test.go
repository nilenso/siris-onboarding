package tddbyexample

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiplication(t *testing.T) {
	five := newDollar(5)
	assert.Equal(t, newDollar(10), five.times(2))
	assert.Equal(t, newDollar(15), five.times(3))
}

func TestEquality(t *testing.T) {
	assert.True(t, newDollar(5).equals(newDollar(5)))
	assert.False(t, newDollar(5).equals(newDollar(6)))
	assert.True(t, newFranc(5).equals(newFranc(5)))
	assert.False(t, newFranc(5).equals(newFranc(6)))
}

func TestFrancMultiplication(t *testing.T) {
	five := newFranc(5)
	assert.Equal(t, newFranc(10), five.times(2))
	assert.Equal(t, newFranc(15), five.times(3))
}
