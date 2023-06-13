package utils

func InStringArray(a []string, s string) bool {
	for _, v := range a {
		if v == s {
			return true
		}
	}
	return false
}

func InInt64Array(a []int64, i int64) bool {
	for _, v := range a {
		if v == i {
			return true
		}
	}
	return false
}

func ShortenString(s string) string {
	if len(s) < 2000 {
		return s
	}
	return s[:950] + "\n...\n" + s[950:]
}

func SumOfFloat64Array(a []float64) float64 {
	var out float64
	for _, v := range a {
		out += v
	}
	return out
}

func MaxOfArray(a []int) int {
	var max int
	for _, v := range a {
		if v > max {
			max = v
		}
	}
	return max
}

func MinOfArray(a []int) int {
	var min int
	for _, v := range a {
		if v < min {
			min = v
		}
	}
	return min
}

func MaxOfFloat64Array(a []float64) float64 {
	var max float64
	for _, v := range a {
		if v > max {
			max = v
		}
	}
	return max
}

func MinOfFloat64Array(a []float64) float64 {
	var min float64
	for _, v := range a {
		if v < min {
			min = v
		}
	}
	return min
}

func orderedIntInsert(a []int, i int) ([]int, int) {
	for index, v := range a {
		if i > v {
			return append(a[:index], append([]int{i}, a[index:]...)...), index
		}
	}
	return append(a, i), len(a)
}

func orderedFloatInsert(a []float64, i float64) ([]float64, int) {
	if i < a[len(a)-1] {
		return append(a, i), len(a)
	}
	for index, v := range a {
		if i > v {
			return append(a[:index], append([]float64{i}, a[index:]...)...), index
		}
	}
	return append(a, i), len(a)
}
