package utility

import (
	"math"
)

// Lerp accepts a min value and max value and a value between 0, 1. Will return the linearly interpolated value between the two first given values.
// i.e. (0, 255, 0.5) = 127
func Lerp(from float64, to float64, value float64) float64 {
	return from*(1-value) + to*value
}

// InverseLerp is similiar to Lerp, but accepts a value that lies between the first two values and returns a value between 0 and 1 depending on where the values lies
// i.e. (0, 255, 127) = 0.5
func InverseLerp(from float64, to float64, value float64) float64 {
	return (value - from) / (to - from)
}

// DegreesToRadian ...
func DegreesToRadian(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

// RadianToDegress ...
func RadianToDegress(radian float64) float64 {
	return radian * (180 / math.Pi)
}
