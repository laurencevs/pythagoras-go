package descent

import (
	"slices"
	"testing"

	"github.com/laurencevs/pythagoras-go/math"
)

func Test_ListGenerators(t *testing.T) {
	minusNSquaredBases := math.GetPrimeFactors(210)
	minusNSquaredBases[-1] = struct{}{}
	minusNSquaredGenerators := getPossibleGenerators(minusNSquaredBases)
	slices.Sort(minusNSquaredGenerators)

	expected := []int64{
		-210, -105, -70, -42, -35, -30, -21, -15,
		-14, -10, -7, -6, -5, -3, -2, -1,
		1, 2, 3, 5, 6, 7, 10, 14,
		15, 21, 30, 35, 42, 70, 105, 210,
	}

	if !slices.Equal(minusNSquaredGenerators, expected) {
		t.Errorf("got unexpected generator list %v, expected %v", minusNSquaredGenerators, expected)
	}
}
