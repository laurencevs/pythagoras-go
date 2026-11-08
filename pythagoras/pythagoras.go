package pythagoras

import (
	"fmt"
	"math/big"

	"github.com/laurencevs/pythagoras-go/curve"
	"github.com/laurencevs/pythagoras-go/internal/set"
)

type RationalTriple struct {
	A, B, C *big.Rat
}

func (t RationalTriple) String() string {
	return fmt.Sprintf("(%s, %s, %s)", t.A.RatString(), t.B.RatString(), t.C.RatString())
}

// Scale the triple's values in place
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

// Convert points to Pythagorean triples with de-duplication.
// Multiple points can correspond to the same triple.
func TriplesFromPoints(points []curve.RationalPoint, n int) []RationalTriple {
	triplesSeen := make(set.StringerSet[RationalTriple])
	var uniqueTriples []RationalTriple
	for _, P := range points {
		if P.IsTwoTorsion() {
			continue
		}
		t := TripleFromPoint(P, n)
		if !triplesSeen.Contains(t) {
			uniqueTriples = append(uniqueTriples, t)
		}
		triplesSeen.Add(t)
	}
	return uniqueTriples
}

func FilterByMaxDigits(triples []RationalTriple, maxDigits int) []RationalTriple {
	var filteredTriples = make([]RationalTriple, 0, len(triples))
	for _, t := range triples {
		if maxDigits > 0 && len(t.C.Denom().String()) > maxDigits {
			continue
		}
		if maxDigits == -1 && (!t.C.Denom().IsInt64() || t.C.Denom().Int64() != 1) {
			continue
		}
		filteredTriples = append(filteredTriples, t)
	}
	return filteredTriples
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
