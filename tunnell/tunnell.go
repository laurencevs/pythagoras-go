package tunnell

import (
	"context"

	"github.com/laurencevs/pythagoras-go/math"
)

type TunnellResult struct {
	Valid bool
	AC    int
	BD    int
}

// n must be positive and square-free
func RunTunnellCheck(_ context.Context, n int64, out chan<- TunnellResult) {
	// TODO: use context to exit early
	out <- CheckTunnell(n)
}

// n must be positive and square-free
func CheckTunnell(n int64) TunnellResult {
	if n%2 == 0 {
		return checkTunnellEven(n)
	}
	return checkTunnellOdd(n)
}

func checkTunnellEven(n int64) TunnellResult {
	c := computeTunnellCount(n, 2, 8, 64)
	d := computeTunnellCount(n, 2, 8, 16)
	return TunnellResult{
		Valid: 2*c == d,
		AC:    c,
		BD:    d,
	}
}

func checkTunnellOdd(n int64) TunnellResult {
	a := computeTunnellCount(n, 1, 2, 32)
	b := computeTunnellCount(n, 1, 2, 8)
	return TunnellResult{
		Valid: 2*a == b,
		AC:    a,
		BD:    b,
	}
}

func computeTunnellCount(n, xc, yc, zc int64) int {
	total := 0
	for x := int64(0); xc*x*x <= n; x++ {
		xContrib := xc * x * x
		for y := int64(0); yc*y*y <= n-xContrib; y++ {
			yContrib := yc * y * y
			z, ok := math.IntSqrt(uint64((n - xContrib - yContrib) / zc))
			if !ok {
				continue
			}
			if xContrib+yContrib+zc*int64(z)*int64(z) != n {
				continue
			}
			numSolutions := 1
			if x != 0 {
				numSolutions *= 2
			}
			if y != 0 {
				numSolutions *= 2
			}
			if z != 0 {
				numSolutions *= 2
			}
			total += numSolutions
		}
	}
	return total
}
