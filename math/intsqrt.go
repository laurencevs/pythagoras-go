package math

func IntSqrt[T ~uint | ~uint64](n T) (T, bool) {
	if n <= 1 {
		return n, true
	}
	bl := bitLength(n) / 2
	var x, y T
	x = T(2) << bl
	for {
		y = (x + n/x) / 2
		if x <= y {
			return x, x*x == n
		}
		x = y
	}
}

func IntFloorSqrt[T ~uint | ~uint64](n T) T {
	sqrt, _ := IntSqrt(n)
	return sqrt
}

func bitLength[T ~uint | ~uint64](n T) T {
	l := T(1)
	for n > 1 {
		l++
		n >>= 1
	}
	return l
}
