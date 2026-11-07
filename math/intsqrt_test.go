package math

import (
	"math"
	"testing"
)

func Test_IntSqrt(t *testing.T) {
	type expected struct {
		n, sqrt uint64
	}
	var cases []expected
	var i uint64
	for i = range 1000 {
		cases = append(cases, expected{
			n:    i,
			sqrt: uint64(math.Floor(math.Sqrt(float64(i)))),
		})
	}
	for _, c := range cases {
		res, _ := IntSqrt(c.n)
		if res != c.sqrt {
			t.Errorf("incorrect square root %d for n = %d", res, c.n)
		}
	}
}

func bitLength[T ~uint | ~uint64](n T) T {
	l := T(1)
	for n > 1 {
		l++
		n >>= 1
	}
	return l
}

func legacyIntSqrt(n uint64) (uint64, bool) {
	bl := bitLength(n) / 2
	var x, y uint64
	x = uint64(2) << bl
	for {
		y = (x + n/x) / 2
		if x <= y {
			return x, x*x == n
		}
		x = y
	}
}

func Test_IntSqrtLarge(t *testing.T) {
	type expected struct {
		n, sqrt uint64
	}
	cases := []expected{
		{0xffffffffffffffff, 0xffffffff},
		{0xfffffffffffffffe, 0xffffffff},
		{0xfffffffffffffff0, 0xffffffff},
		{0xfffffffe00000002, 0xffffffff},
		{0xfffffffe00000001, 0xffffffff},
		{0xfffffffe00000000, 0xfffffffe},
		{0xfffffffdffffffff, 0xfffffffe},
	}
	for n := (uint64(1) << 32) - 1; n > (uint64(1)<<32)-10000; n-- {
		cases = append(cases, expected{
			n: n * n, sqrt: n,
		})
		cases = append(cases, expected{
			n: n*n - 1, sqrt: n - 1,
		})
		cases = append(cases, expected{
			n: n*n + 1, sqrt: n,
		})
	}

	for n2 := (uint64(1) << 54) + 1000; n2 > (uint64(1)<<54)-1000; n2-- {
		s, _ := legacyIntSqrt(n2)
		cases = append(cases, expected{
			n: n2, sqrt: s,
		})
	}
	for n2 := (uint64(1) << 53) + 1000; n2 > (uint64(1)<<53)-1000; n2-- {
		s, _ := legacyIntSqrt(n2)
		cases = append(cases, expected{
			n: n2, sqrt: s,
		})
	}
	for n2 := (uint64(1) << 52) + 1000; n2 > (uint64(1)<<52)-1000; n2-- {
		s, _ := legacyIntSqrt(n2)
		cases = append(cases, expected{
			n: n2, sqrt: s,
		})
	}

	for _, c := range cases {
		res, _ := IntSqrt(c.n)
		if res != c.sqrt {
			t.Errorf("incorrect square root %d for n = %d", res, c.n)
		}
	}
}
