package math

import (
	stdmath "math"
)

func IntSqrt[T ~uint | ~uint64](n T) (T, bool) {
	if n <= 1 {
		return n, true
	}

	s := T(stdmath.Sqrt(float64(n)))
	if s*s > n || s == 0x100000000 {
		s--
	} else if (s+1)*(s+1) <= n && s != 0xffffffff {
		s++
	}

	return s, s*s == n
}

func IntFloorSqrt[T ~uint | ~uint64](n T) T {
	sqrt, _ := IntSqrt(n)
	return sqrt
}
