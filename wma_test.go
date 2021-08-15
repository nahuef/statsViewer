package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWMA(t *testing.T) {

	wma := NewWMA(7)
	avg, count := wma.Average()

	assert.Equal(t, float64(0), avg)
	assert.Equal(t, 0, count)

	wma.Add(0)
	avg, count = wma.Average()
	assert.Equal(t, float64(0), avg)
	assert.Equal(t, 1, count)

	wma = NewWMA(1)
	wma.Add(7)
	wma.Add(6)
	wma.Add(2)
	avg, count = wma.Average()
	assert.Equal(t, float64(2), avg)
	assert.Equal(t, 1, count)

	wma = NewWMA(3)
	wma.Add(6)
	wma.Add(6)
	avg, count = wma.Average()
	assert.Equal(t, float64(6), avg)
	assert.Equal(t, 2, count)

	wma = NewWMA(7)
	wma.Add(7)
	wma.Add(6)
	wma.Add(4)
	avg, count = wma.Average()
	assert.Equal(t, float64(5.5), avg)
	assert.Equal(t, 3, count)

	wma = NewWMA(7)
	wma.Add(7, 6)
	wma.Add(6, 4)
	wma.Add(1, 3, 15)
	avg, count = wma.Average()
	assert.Equal(t, float64(6), avg)
	assert.Equal(t, 7, count)
}
