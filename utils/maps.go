package utils

import (
	"sort"
)

func TopNOfIntMap(m map[string]int, n int) []string {

	var out []string
	var values []int
	var index int

	for s, i := range m {
		values, index = orderedIntInsert(values, i)
		out = append(out[:index], append([]string{s}, out[index:]...)...)
	}

	var l int = n
	if len(out) < n {
		l = len(out)
	}

	out = out[:l]

	return out
}

func TopNOfFloat64Map(m map[string]float64, n int) []string {

	var out []string

	var values []float64
	var invMap = make(map[float64]string)

	for k, v := range m {
		values = append(values, v)
		invMap[v] = k
	}

	sort.Float64s(values[:])

	for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
		values[i], values[j] = values[j], values[i]
	}

	var l int = n
	if len(values) < n {
		l = len(values)
	}

	for _, v := range values[:l] {
		out = append(out, invMap[v])
	}

	return out
}
