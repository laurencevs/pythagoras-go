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
	primes := ListPrimesUpTo(n)
	factors := make(set.Set[int64])
	for p := range primes {
		if n%p == 0 {
			factors.Add(p)
		}
	}
	return factors
}
