package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/laurencevs/pythagoras-go/math"
)

type rationalTriple struct {
	A, B, C *big.Rat
}

func (t rationalTriple) String() string {
	return fmt.Sprintf("(%s, %s, %s)", t.A.RatString(), t.B.RatString(), t.C.RatString())
}

func (t rationalTriple) Scale(factor int) {
	factorRat := big.NewRat(int64(factor), 1)
	t.A.Mul(t.A, factorRat)
	t.B.Mul(t.B, factorRat)
	t.C.Mul(t.C, factorRat)
}

func (t rationalTriple) ToPoint(n int64) rationalPoint {
	p := rationalPoint{x: new(big.Rat), y: new(big.Rat)}
	p.x.Add(t.A, t.C).Mul(p.x, t.A).Quo(p.x, twoRat)
	p.y.Set(p.x)
	p.y.Mul(p.y, t.A)
	return p
}

func (p rationalPoint) ToTriple(n int) rationalTriple {
	var a, b, c big.Rat
	a.Mul(p.x, p.x)
	c.Set(&a)
	n2 := new(big.Rat).SetInt64(int64(n * n))
	a.Sub(&a, n2).Quo(&a, p.y)
	c.Add(&c, n2).Quo(&c, p.y)
	b.SetFrac64(int64(2*n), 1).Mul(&b, p.x).Quo(&b, p.y)
	return rationalTriple{
		A: a.Abs(&a),
		B: b.Abs(&b),
		C: c.Abs(&c),
	}
}

// Primitive Pythagorean triples (a, b, c) are represented by (u, v) where
//
//	a=u*u-v*v, b=2*u*v, c=u*u+v*v.
//
// sqrt is the square root (by construction, an integer) of
//
//	u*v*(u*u-v*v)*n
//
// so that
//
//	(a*n/sqrt, b*n/sqrt, c*n/sqrt)
//
// are the side lengths of a rational Pythagorean triangle of area n.
type primitiveTriple struct {
	u, v int64
	sqrt *big.Int
}

func (t primitiveTriple) ToPoint(n int64) rationalPoint {
	var a, b, c big.Rat

	nDivSqrt := new(big.Rat)
	nDivSqrt.SetInt64(int64(n))
	nDivSqrt.Quo(nDivSqrt, new(big.Rat).SetInt(t.sqrt))

	a.SetInt64(t.u*t.u-t.v*t.v).Mul(&a, nDivSqrt)
	b.SetInt64(2*t.u*t.v).Mul(&b, nDivSqrt)
	c.SetInt64(t.u*t.u+t.v*t.v).Mul(&c, nDivSqrt)

	r := rationalTriple{
		A: &a,
		B: &b,
		C: &c,
	}
	return r.ToPoint(n)
}

const maxInt = int64((^uint64(0)) >> 1)

var sqrtMaxInt, _ = math.IntSqrt(uint64(maxInt))

func initialPointSearch(ctx context.Context, n int64, out chan<- rationalPoint) {
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
			if (u%2 == 1 && v%2 == 1) || gcd(u, v) > 1 {
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
			case out <- primitiveTriple{u, v, s}.ToPoint(n):
			}
		}
	}
}
