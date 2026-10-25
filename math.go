package main

func gcd(a, b int) int {
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

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
