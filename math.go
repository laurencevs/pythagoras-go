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

func intSqrt(n uint) (uint, bool) {
	if n <= 1 {
		return n, true
	}
	bl := bitLength(n) / 2
	var x, y uint
	x = uint(2) << bl
	for {
		y = (x + n/x) / 2
		if y < x {
			x = y
		} else {
			return x, x*x == n
		}
	}
}

func bitLength(n uint) uint {
	l := uint(1)
	for n > 1 {
		l++
		n >>= 1
	}
	return l
}
