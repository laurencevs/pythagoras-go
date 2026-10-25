package main

func gcd[T ~int | ~int64](a, b T) T {
	if abs(a) < abs(b) {
		a, b = b, a
	}
	if b == 0 {
		return a
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs[T ~int | ~int64](n T) T {
	if n < 0 {
		return -n
	}
	return n
}
