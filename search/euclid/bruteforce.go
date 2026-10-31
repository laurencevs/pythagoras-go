package euclid

import (
	"context"
	"fmt"
	"math/big"

	"github.com/laurencevs/pythagoras-go/curve"
	"github.com/laurencevs/pythagoras-go/math"
	"github.com/laurencevs/pythagoras-go/pythagoras"
)

const maxInt = int64((^uint64(0)) >> 1)

var sqrtMaxInt, _ = math.IntSqrt(uint64(maxInt))

func SearchForInitialPoints(ctx context.Context, n int64, out chan<- curve.RationalPoint) {
	limUint := math.IntFloorSqrt(math.IntFloorSqrt(uint64(maxInt / n)))
	lim := int64(limUint)

	if limUint*limUint*uint64(n) > ((^uint64(0))>>1)/(limUint*limUint) {
		panic(fmt.Sprintf("fatal: limUint is too big: %d", limUint))
	}

	for u := int64(1); u < int64(sqrtMaxInt/2); u++ {
		select {
		case <-ctx.Done():
			return
		default:
		}

		for v := int64(1); v < u; v++ {
			if (u%2 == 1 && v%2 == 1) || math.GCD(u, v) > 1 {
				continue
			}

			var isValidTriangle bool
			var s *big.Int
			if u < lim {
				var sInt uint
				sInt, isValidTriangle = math.IntSqrt(uint((u*u - v*v) * u * v * n))
				if isValidTriangle {
					s = big.NewInt(int64(sInt))
				}
			} else {
				s2 := big.NewInt(u*u - v*v)
				s2.Mul(s2, big.NewInt(u*v))
				s2.Mul(s2, big.NewInt(n))
				s = new(big.Int).Sqrt(s2)
				isValidTriangle = new(big.Int).Mul(s, s).Cmp(s2) == 0
			}

			if !isValidTriangle {
				continue
			}

			select {
			case <-ctx.Done():
				return
			case out <- pythagoras.PrimitiveTriple{U: u, V: v, Sqrt: s}.ToPoint(n):
			}
		}
	}
}
