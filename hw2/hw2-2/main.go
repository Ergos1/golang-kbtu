package main

import (
	"fmt"
	"math"
)

const eps = 1e-6

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {

	if x == 0 {
		return 0, nil
	} else if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}

	z := float64(1)
	prev := z
	for {
		z -= (z*z - x) / (2 * z)
		if math.Abs(prev-z) <= eps {
			break
		}

		prev = z
	}
	return z, nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
