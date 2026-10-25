package math

import "testing"

func Test_GetSquarefree(t *testing.T) {
	samples := []struct{ n, squaredPart, squareFreePart int }{
		{1, 1, 1},
		{2, 1, 2},
		{3, 1, 3},
		{4, 2, 1},
		{5, 1, 5},
		{6, 1, 6},
		{7, 1, 7},
		{8, 2, 2},
		{9, 3, 1},
		{10, 1, 10},
		{11, 1, 11},
		{12, 2, 3},
	}
	for _, s := range samples {
		squareFreePart, squaredPart := GetSquarefree(s.n)
		if squareFreePart != s.squareFreePart {
			t.Errorf("got square-free part %d, expected %d for n = %d", squareFreePart, s.squareFreePart, s.n)
		}
		if squaredPart != s.squaredPart {
			t.Errorf("got squared part %d, expected %d for n = %d", squareFreePart, s.squareFreePart, s.n)
		}
	}
}
