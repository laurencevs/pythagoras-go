package pythagoras

import (
	"fmt"
	"math/big"

	"github.com/laurencevs/pythagoras-go/curve"
)

type RationalTriple struct {
	A, B, C *big.Rat
}

func (t RationalTriple) String() string {
	return fmt.Sprintf("(%s, %s, %s)", t.A.RatString(), t.B.RatString(), t.C.RatString())
}

func (t RationalTriple) Scale(factor int) {
	factorRat := big.NewRat(int64(factor), 1)
	t.A.Mul(t.A, factorRat)
	t.B.Mul(t.B, factorRat)
	t.C.Mul(t.C, factorRat)
}

var twoRat = big.NewRat(2, 1)

func (t RationalTriple) ToPoint(n int64) curve.RationalPoint {
	p := curve.RationalPoint{X: new(big.Rat), Y: new(big.Rat)}
	p.X.Add(t.A, t.C).Mul(p.X, t.A).Quo(p.X, twoRat)
	p.Y.Set(p.X)
	p.Y.Mul(p.Y, t.A)
	return p
}

func TripleFromPoint(p curve.RationalPoint, n int) RationalTriple {
	var a, b, c big.Rat
	a.Mul(p.X, p.X)
	c.Set(&a)
	n2 := new(big.Rat).SetInt64(int64(n * n))
	a.Sub(&a, n2).Quo(&a, p.Y)
	c.Add(&c, n2).Quo(&c, p.Y)
	b.SetFrac64(int64(2*n), 1).Mul(&b, p.X).Quo(&b, p.Y)
	a.Abs(&a)
	b.Abs(&b)
	if a.Cmp(&b) == 1 {
		a, b = b, a
	}
	return RationalTriple{
		A: &a,
		B: &b,
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
type PrimitiveTriple struct {
	U, V int64
	Sqrt *big.Int
}

func (t PrimitiveTriple) ToPoint(n int64) curve.RationalPoint {
	var a, b, c big.Rat

	nDivSqrt := new(big.Rat)
	nDivSqrt.SetInt64(int64(n))
	nDivSqrt.Quo(nDivSqrt, new(big.Rat).SetInt(t.Sqrt))

	a.SetInt64(t.U*t.U-t.V*t.V).Mul(&a, nDivSqrt)
	b.SetInt64(2*t.U*t.V).Mul(&b, nDivSqrt)
	c.SetInt64(t.U*t.U+t.V*t.V).Mul(&c, nDivSqrt)

	r := RationalTriple{
		A: &a,
		B: &b,
		C: &c,
	}
	return r.ToPoint(n)
}
