package math

import (
	"testing"
)

func setEqual[T comparable](a, b map[T]struct{}) bool {
	if len(a) != len(b) {
		return false
	}
	for x := range a {
		if _, ok := b[x]; !ok {
			return false
		}
	}
	return true
}

func Test_ListPrimesUpTo_SmallN(t *testing.T) {
	upTo1 := ListPrimesUpTo(1)
	if !setEqual(upTo1, map[int64]struct{}{}) {
		t.Errorf("unexpected primes up to 1: %v", upTo1)
	}

	upTo2 := ListPrimesUpTo(2)
	if !setEqual(upTo2, map[int64]struct{}{2: {}}) {
		t.Errorf("unexpected primes up to 2: %v", upTo2)
	}

	upTo3 := ListPrimesUpTo(3)
	if !setEqual(upTo3, map[int64]struct{}{2: {}, 3: {}}) {
		t.Errorf("unexpected primes up to 3: %v", upTo3)
	}

	upTo4 := ListPrimesUpTo(4)
	if !setEqual(upTo4, map[int64]struct{}{2: {}, 3: {}}) {
		t.Errorf("unexpected primes up to 4: %v", upTo4)
	}

	upTo5 := ListPrimesUpTo(5)
	if !setEqual(upTo5, map[int64]struct{}{2: {}, 3: {}, 5: {}}) {
		t.Errorf("unexpected primes up to 5: %v", upTo5)
	}
}

func Test_ListPrimesUpTo(t *testing.T) {
	upTo78 := ListPrimesUpTo(78)

	expectedList := []int64{
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73,
	}
	expected := make(map[int64]struct{})
	for _, p := range expectedList {
		expected[p] = struct{}{}
	}

	if !setEqual(upTo78, expected) {
		t.Errorf("unexpected primes up to 78: %v", upTo78)
	}

	expected[79] = struct{}{}
	upTo79 := ListPrimesUpTo(79)
	if !setEqual(upTo79, expected) {
		t.Errorf("unexpected primes up to 79: %v", upTo79)
	}

	upTo80 := ListPrimesUpTo(80)
	if !setEqual(upTo80, expected) {
		t.Errorf("unexpected primes up to 80: %v", upTo80)
	}
}

func Test_GetPrimeFactors(t *testing.T) {
	f1 := GetPrimeFactors(1)
	if !setEqual(f1, map[int64]struct{}{}) {
		t.Errorf("unexpected factors of 1: %v", f1)
	}
	f2 := GetPrimeFactors(2)
	if !setEqual(f2, map[int64]struct{}{2: {}}) {
		t.Errorf("unexpected factors of 2: %v", f2)
	}
	f4 := GetPrimeFactors(4)
	if !setEqual(f4, map[int64]struct{}{2: {}}) {
		t.Errorf("unexpected factors of 4: %v", f4)
	}
	f6 := GetPrimeFactors(6)
	if !setEqual(f6, map[int64]struct{}{2: {}, 3: {}}) {
		t.Errorf("unexpected factors of 6: %v", f6)
	}
	f15 := GetPrimeFactors(15)
	if !setEqual(f15, map[int64]struct{}{3: {}, 5: {}}) {
		t.Errorf("unexpected factors of 15: %v", f15)
	}
	f70 := GetPrimeFactors(70)
	if !setEqual(f70, map[int64]struct{}{2: {}, 5: {}, 7: {}}) {
		t.Errorf("unexpected factors of 70: %v", f70)
	}
	f1254 := GetPrimeFactors(1254)
	if !setEqual(f1254, map[int64]struct{}{2: {}, 3: {}, 11: {}, 19: {}}) {
		t.Errorf("unexpected factors of 1254: %v", f1254)
	}
	f1254x2x11 := GetPrimeFactors(1254 * 2 * 11)
	if !setEqual(f1254x2x11, map[int64]struct{}{2: {}, 3: {}, 11: {}, 19: {}}) {
		t.Errorf("unexpected factors of 1254 * 2 * 11: %v", f1254x2x11)
	}
}
