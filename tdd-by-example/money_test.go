package tddbyexample

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiplication(t *testing.T) {
	five := newDollar(5)
	five.times(2)
	assert.Equal(t, 10, five.amount)
}
