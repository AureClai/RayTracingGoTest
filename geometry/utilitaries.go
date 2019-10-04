package geometry

import "math/rand"

func Fmin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func Fmax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func drand48() float64 {
	return rand.Float64()
}
