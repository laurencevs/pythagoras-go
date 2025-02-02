package main

import (
	"fmt"
	"math/big"
	"time"
)

type rationalTriple struct {
	A, B, C *big.Rat
}

func (t rationalTriple) String() string {
	return fmt.Sprintf("(%s, %s, %s)", t.A.RatString(), t.B.RatString(), t.C.RatString())
}

func (t rationalTriple) ToPoint(n int) rationalPoint {
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
	u, v int
	sqrt uint
}

func (t primitiveTriple) ToPoint(n int) rationalPoint {
	r := rationalTriple{
		A: big.NewRat(int64((t.u*t.u-t.v*t.v)*n), int64(t.sqrt)),
		B: big.NewRat(int64(2*t.u*t.v*n), int64(t.sqrt)),
		C: big.NewRat(int64((t.u*t.u+t.v*t.v)*n), int64(t.sqrt)),
	}
	return r.ToPoint(n)
}

func initialPointSearch(n int, timeout <-chan time.Time) []rationalPoint {
	var triples []rationalPoint
	for u := 1; ; u++ {
		select {
		case <-timeout:
			return triples
		default:
			for v := 1; v < u; v++ {
				if (u%2 == 1 && v%2 == 1) || gcd(u, v) > 1 {
					continue
				}
				if s, ok := intSqrt(uint((u*u - v*v) * u * v * n)); ok {
					triples = append(triples, primitiveTriple{u, v, s}.ToPoint(n))
					if timeout == nil {
						return triples
					}
				}
			}
		}
	}
}
