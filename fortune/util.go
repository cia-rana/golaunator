package fortune

import "math"

func Clamp(a, b, x float64) float64 {
	return math.Min(math.Max(a, x), b)
}
