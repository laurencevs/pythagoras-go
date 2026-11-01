package math

func GetSquarefree(n uint) (uint, uint) {
	squaredPart := uint(1)
	sqrt := IntFloorSqrt(uint(n))
	for x := sqrt; x > 1; x-- {
		if n%(x*x) == 0 {
			squaredPart = x
			break
		}
	}
	return n / (squaredPart * squaredPart), squaredPart
}
