package tunnell

import "testing"

type knownResult struct {
	n     int
	valid bool
}

var knownResults = []knownResult{
	{1, false},
	{2, false},
	{3, false},
	// {4, false}, not square-free
	{5, true},
	{6, true},
	{7, true},
	// {8, false},
	// {9, false},
	{10, false},
	{11, false},
	// {12, false},
	{13, true},
	{14, true},
	{15, true},
	// {16, false},
	{17, false},
	// {18, false},
	{19, false},
	// {20, true},
	{21, true},

	{210, true},
	{1254, true},
}

func TestTunnellValid(t *testing.T) {
	for _, exp := range knownResults {
		res := CheckTunnell(exp.n)
		if res.Valid != exp.valid {
			t.Errorf("got %v, expected %v for n = %d", res.Valid, exp.valid, exp.n)
		}
	}
}
