package math

import (
	"testing"

	"github.com/laurencevs/pythagoras-go/internal/set"
)

func Test_ListPrimesUpTo_SmallN(t *testing.T) {
	upTo1 := ListPrimesUpTo(1)
	if !upTo1.Equals(set.New[int64]()) {
		t.Errorf("unexpected primes up to 1: %v", upTo1)
	}

	upTo2 := ListPrimesUpTo(2)
	if !upTo2.Equals(set.New[int64](2)) {
		t.Errorf("unexpected primes up to 2: %v", upTo2)
	}

	upTo3 := ListPrimesUpTo(3)
	if !upTo3.Equals(set.New[int64](2, 3)) {
		t.Errorf("unexpected primes up to 3: %v", upTo3)
	}

	upTo4 := ListPrimesUpTo(4)
	if !upTo4.Equals(set.New[int64](2, 3)) {
		t.Errorf("unexpected primes up to 4: %v", upTo4)
	}

	upTo5 := ListPrimesUpTo(5)
	if !upTo5.Equals(set.New[int64](2, 3, 5)) {
		t.Errorf("unexpected primes up to 5: %v", upTo5)
	}
}

func Test_ListPrimesUpTo(t *testing.T) {
	upTo78 := ListPrimesUpTo(78)

	expected := set.New[int64](
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73,
	)

	if !expected.Equals(upTo78) {
		t.Errorf("unexpected primes up to 78: %v", upTo78)
	}

	expected.Add(79)
	upTo79 := ListPrimesUpTo(79)
	if !upTo79.Equals(expected) {
		t.Errorf("unexpected primes up to 79: %v", upTo79)
	}

	upTo80 := ListPrimesUpTo(80)
	if !upTo80.Equals(expected) {
		t.Errorf("unexpected primes up to 80: %v", upTo80)
	}
}

func Test_GetPrimeFactors(t *testing.T) {
	f1 := GetPrimeFactors(1)
	if !f1.Equals(set.New[int64]()) {
		t.Errorf("unexpected factors of 1: %v", f1)
	}
	f2 := GetPrimeFactors(2)
	if !f2.Equals(set.New[int64](2)) {
		t.Errorf("unexpected factors of 2: %v", f2)
	}
	f4 := GetPrimeFactors(4)
	if !f4.Equals(set.New[int64](2)) {
		t.Errorf("unexpected factors of 4: %v", f4)
	}
	f6 := GetPrimeFactors(6)
	if !f6.Equals(set.New[int64](2, 3)) {
		t.Errorf("unexpected factors of 6: %v", f6)
	}
	f15 := GetPrimeFactors(15)
	if !f15.Equals(set.New[int64](3, 5)) {
		t.Errorf("unexpected factors of 15: %v", f15)
	}
	f70 := GetPrimeFactors(70)
	if !f70.Equals(set.New[int64](2, 5, 7)) {
		t.Errorf("unexpected factors of 70: %v", f70)
	}
	f1254 := GetPrimeFactors(1254)
	if !f1254.Equals(set.New[int64](2, 3, 11, 19)) {
		t.Errorf("unexpected factors of 1254: %v", f1254)
	}
	f1254x2x11 := GetPrimeFactors(1254 * 2 * 11)
	if !f1254x2x11.Equals(set.New[int64](2, 3, 11, 19)) {
		t.Errorf("unexpected factors of 1254 * 2 * 11: %v", f1254x2x11)
	}
}
