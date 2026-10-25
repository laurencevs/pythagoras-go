package math

func GetSquarefree(n int) (int, int) {
	squaredPart := 1
	for x := n - 1; x > 1; x-- {
		if n%(x*x) == 0 {
			squaredPart = x
			break
		}
	}
	return n / (squaredPart * squaredPart), squaredPart
}
