package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := float64(1)
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}

func SqrtWithOtherCondition(x float64, eps float64) (float64, int) {
	z := float64(1)
	prev := z
	itter := 0
	for {
		z -= (z*z - x) / (2 * z)
		itter += 1
		if math.Abs(prev-z) <= eps {
			break
		}

		prev = z
	}
	return z, itter
}

func main() {
	fmt.Printf("SIMPLE sqrt : res = %.30f; itter = %d\n", Sqrt(float64(2)), 10)
	res, itter := SqrtWithOtherCondition(float64(0), 1e-6)
	fmt.Printf("NOT_SIMPLE sqrt : res = %.30f; itter = %d\n", res, itter)
}
