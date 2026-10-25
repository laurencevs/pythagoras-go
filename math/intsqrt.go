package math

func IntSqrt(n uint) (uint, bool) {
	if n <= 1 {
		return n, true
	}
	bl := bitLength(n) / 2
	var x, y uint
	x = uint(2) << bl
	for {
		y = (x + n/x) / 2
		if x <= y {
			return x, x*x == n
		}
		x = y
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
