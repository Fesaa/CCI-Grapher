package utils

func InStringArray(a []string, s string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}
	return false
}

func SumOfFloat64Array(a []float64) float64 {
	var out float64
	for _, v := range a {
		out += v
	}
	return out
}

func orderedIntInsert(a []int, i int) ([]int, int) {
	for index, v := range a {
		if i > v {
			return append(a[:index], append([]int{i}, a[index:]...)...), index
		}
	}
	return append(a, i), len(a)
}
