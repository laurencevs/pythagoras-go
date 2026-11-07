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
func RunTunnellCheck(ctx context.Context, n int64, out chan<- TunnellResult) {
	res, ok := CheckTunnell(ctx, n)
	if !ok {
		return
	}
	select {
	case out <- res:
	case <-ctx.Done():
	}
}

// n must be positive and square-free
func CheckTunnell(ctx context.Context, n int64) (TunnellResult, bool) {
	var ac, bd int
	if n%2 == 1 {
		a := computeTunnellCount(ctx, n, 1, 2, 32)
		b := computeTunnellCount(ctx, n, 1, 2, 8)
		ac, bd = a, b
	} else {
		c := computeTunnellCount(ctx, n, 2, 8, 64)
		d := computeTunnellCount(ctx, n, 2, 8, 16)
		ac, bd = c, d
	}

	if err := ctx.Err(); err != nil {
		return TunnellResult{}, false
	}

	return TunnellResult{
		Valid: 2*ac == bd,
		AC:    ac,
		BD:    bd,
	}, true
}

func computeTunnellCount(ctx context.Context, n, xc, yc, zc int64) int {
	total := 0
	for x := int64(0); xc*x*x <= n; x++ {
		select {
		case <-ctx.Done():
			return 0
		default:
		}
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
