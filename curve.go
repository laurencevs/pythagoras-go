package main

import (
	"fmt"
	"math/big"
)

type weierstrassCurve struct {
	a, b int64
}

type rationalPoint struct {
	x, y *big.Rat
}

func (p rationalPoint) String() string {
	if p.y == nil {
		return "(inf)"
	}
	return fmt.Sprintf("(%s, %s)", p.x.RatString(), p.y.RatString())
}

var inf = rationalPoint{}

func (c weierstrassCurve) Contains(p rationalPoint) bool {
	lhs := new(big.Rat)
	lhs.Mul(p.y, p.y)
	axb := new(big.Rat)
	axb.SetInt64(c.a).Mul(axb, p.x).Add(axb, big.NewRat(c.b, 1))
	rhs := new(big.Rat)
	rhs.Mul(p.x, p.x).Mul(rhs, p.x).Add(rhs, axb)
	return lhs.Cmp(rhs) == 0
}

func (c weierstrassCurve) Invert(p rationalPoint) rationalPoint {
	return rationalPoint{
		x: p.x,
		y: new(big.Rat).Neg(p.y),
	}
}

var (
	zeroInt = big.NewInt(0)
	twoInt  = big.NewInt(2)
	twoRat  = big.NewRat(2, 1)
)

func (c weierstrassCurve) Double(p rationalPoint) rationalPoint {
	if p.y == nil || p.y.Num().Cmp(zeroInt) == 0 {
		return inf
	}
	s := big.NewRat(3, 1)
	s.Mul(s, p.x).Mul(s, p.x).Add(s, big.NewRat(c.a, 1)).Quo(s, p.y)
	s.Denom().Mul(s.Denom(), twoInt)
	x := new(big.Rat).Mul(s, s)
	x.Sub(x, new(big.Rat).Mul(twoRat, p.x))
	y := new(big.Rat).Sub(p.y, new(big.Rat).Mul(s, new(big.Rat).Sub(p.x, x)))
	return c.Invert(rationalPoint{x, y})
}

func (c weierstrassCurve) Add(p, q rationalPoint) rationalPoint {
	if p.y == nil {
		return q
	}
	if q.y == nil {
		return p
	}
	if p.x.Cmp(q.x) == 0 {
		if p.y.Cmp(q.y) == 0 {
			return c.Double(p)
		}
		return inf
	}
	s := new(big.Rat).Sub(p.y, q.y)
	s.Quo(s, new(big.Rat).Sub(p.x, q.x))
	x := new(big.Rat).Mul(s, s)
	x.Sub(x, p.x).Sub(x, q.x)
	y := new(big.Rat).Sub(p.y, new(big.Rat).Mul(s, new(big.Rat).Sub(p.x, x)))
	return c.Invert(rationalPoint{x, y})
}
