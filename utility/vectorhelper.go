package utility

import "math"

// GetDistance ...
func GetDistance(x, y, x2, y2 int32) int32 {
	// using Pythagorean theorem
	dx := x - x2
	dy := y - y2

	dr := dx*dx + dy*dy
	return int32(math.Round(math.Sqrt(math.Abs(float64(dr)))))
}

// GetDistanceInFloat ...
func GetDistanceInFloat(x, y, x2, y2 int32) float64 {
	// using Pythagorean theorem
	dx := x - x2
	dy := y - y2

	dr := dx*dx + dy*dy
	return math.Sqrt(math.Abs(float64(dr)))
}
