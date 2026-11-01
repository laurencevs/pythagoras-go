package math

import "github.com/laurencevs/pythagoras-go/internal/set"

// TODO: these are horribly inefficient proofs-of-concept

func ListPrimesUpTo(n int64) set.Set[int64] {
	if n < 2 {
		return set.New[int64]()
	}
	primes := set.New[int64](2)
outer:
	for i := int64(3); i <= n; i += 2 {
		for p := range primes {
			if i%p == 0 {
				continue outer
			}
		}
		primes.Add(i)
	}
	return primes
}

func GetPrimeFactors(n int64) set.Set[int64] {
	if n < 2 {
		return set.New[int64]()
	}
	factors := set.New[int64]()
	for n%2 == 0 {
		factors.Add(2)
		n /= 2
	}
	oddPrimes := set.New[int64]()
	p := int64(3)
outer:
	for n > 1 {
		for op := range oddPrimes {
			if p%op == 0 {
				p += 2
				continue outer
			}
		}
		oddPrimes.Add(p)

		for n%p == 0 {
			factors.Add(p)
			n /= p
		}

		p += 2
	}

	return factors
}
