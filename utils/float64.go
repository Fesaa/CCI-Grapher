package utils

import (
	"math"
)

func FloorToN(f float64, n int) float64 {
	return math.Floor(f*math.Pow10(n)) / math.Pow10(n)
}
