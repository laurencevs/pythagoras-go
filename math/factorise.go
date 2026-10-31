package math

// TODO: these are horribly inefficient proofs-of-concept

func ListPrimesUpTo(n int64) map[int64]struct{} {
	if n < 2 {
		return map[int64]struct{}{}
	}
	primes := map[int64]struct{}{2: {}}
outer:
	for i := int64(3); i <= n; i += 2 {
		for p := range primes {
			if i%p == 0 {
				continue outer
			}
		}
		primes[i] = struct{}{}
	}
	return primes
}

func GetPrimeFactors(n int64) map[int64]struct{} {
	primes := ListPrimesUpTo(n)
	factors := make(map[int64]struct{})
	for p := range primes {
		if n%p == 0 {
			factors[p] = struct{}{}
		}
	}
	return factors
}
