package main

// TODO: Make this configurable
const DefaultWMAWindow = 7

type WeightedMovingAverage struct {
	N int
	// Values in each group receive the same weight
	groups [][]float64
}

func NewWMA(n int) *WeightedMovingAverage {
	return &WeightedMovingAverage{
		N: n,
	}
}

func (wma *WeightedMovingAverage) Add(values ...float64) {
	n := wma.N
	if len(wma.groups) >= n {
		wma.groups = wma.groups[1:n]
	}

	wma.groups = append(wma.groups, values)
}

func (wma *WeightedMovingAverage) Average() (float64, int) {
	groups := wma.groups
	if groups == nil {
		return 0, 0
	}

	weight := wma.N - len(groups) + 1
	weightTotal := 0
	sum := float64(0)
	count := 0
	for _, group := range groups {
		for _, value := range group {
			sum += value * float64(weight)
			weightTotal += weight
			count++
		}
		weight++
	}

	avg := sum / float64(weightTotal)
	return avg, count
}
