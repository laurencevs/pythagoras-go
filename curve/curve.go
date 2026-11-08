package curve

import (
	"fmt"
	"math/big"
)

type WeierstrassCurve struct {
	A, B int64
}

type RationalPoint struct {
	X, Y *big.Rat
}

func (p RationalPoint) String() string {
	if p.IsInf() {
		return "(inf)"
	}
	return fmt.Sprintf("(%s, %s)", p.X.RatString(), p.Y.RatString())
}

func (p RationalPoint) IsInf() bool {
	return p.Y == nil
}

func (p RationalPoint) IsTwoTorsion() bool {
	return p.Y == nil || p.Y.Cmp(zeroRat) == 0
}

var Inf = RationalPoint{}

func (c WeierstrassCurve) Contains(p RationalPoint) bool {
	if p.IsInf() {
		return true
	}
	lhs := new(big.Rat)
	lhs.Mul(p.Y, p.Y)
	axb := new(big.Rat)
	axb.SetInt64(c.A).Mul(axb, p.X).Add(axb, big.NewRat(c.B, 1))
	rhs := new(big.Rat)
	rhs.Mul(p.X, p.X).Mul(rhs, p.X).Add(rhs, axb)
	return lhs.Cmp(rhs) == 0
}

func (c WeierstrassCurve) Invert(p RationalPoint) RationalPoint {
	if p.IsInf() {
		return Inf
	}
	return RationalPoint{
		X: p.X,
		Y: new(big.Rat).Neg(p.Y),
	}
}

var (
	zeroRat = big.NewRat(0, 1)
	twoRat  = big.NewRat(2, 1)
)

func (c WeierstrassCurve) Double(p RationalPoint) RationalPoint {
	if p.IsTwoTorsion() {
		return Inf
	}
	s := big.NewRat(3, 1)
	s.Mul(s, p.X).Mul(s, p.X).Add(s, big.NewRat(c.A, 1)).Quo(s, p.Y).Quo(s, twoRat)
	x := new(big.Rat).Mul(s, s)
	x.Sub(x, new(big.Rat).Mul(twoRat, p.X))
	y := new(big.Rat).Sub(p.Y, new(big.Rat).Mul(s, new(big.Rat).Sub(p.X, x)))
	y.Neg(y)
	return RationalPoint{x, y}
}

func (c WeierstrassCurve) Add(p, q RationalPoint) RationalPoint {
	if p.IsInf() {
		return q
	}
	if q.IsInf() {
		return p
	}
	if p.X.Cmp(q.X) == 0 {
		if p.Y.Cmp(q.Y) == 0 {
			return c.Double(p)
		}
		return Inf
	}
	s := new(big.Rat).Sub(p.Y, q.Y)
	s.Quo(s, new(big.Rat).Sub(p.X, q.X))
	x := new(big.Rat).Mul(s, s)
	x.Sub(x, p.X).Sub(x, q.X)
	y := new(big.Rat).Sub(p.Y, new(big.Rat).Mul(s, new(big.Rat).Sub(p.X, x)))
	y.Neg(y)
	return RationalPoint{x, y}
}
