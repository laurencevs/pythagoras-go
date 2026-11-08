package curve

import (
	"github.com/laurencevs/pythagoras-go/internal/set"
)

func MultiplyPoints(E WeierstrassCurve, initialPoints []RationalPoint, maxCoeff int) []RationalPoint {
	// Multiply the initial points and compute all linear combinations thereof,
	// ignoring duplicates.
	var points []RationalPoint
	seen := make(set.StringerSet[RationalPoint])
	for _, c := range initialPoints {
		if seen.Contains(c) || seen.Contains(E.Invert(c)) {
			continue
		}
		var newPoints []RationalPoint
		Q := Inf
		for i := 0; i < maxCoeff; i++ {
			Q = E.Add(Q, c)
			NQ := E.Invert(Q)
			newPoints = append(newPoints, Q, NQ)
			seen.Add(Q)
			seen.Add(NQ)
			for _, p := range points {
				pQ, pNQ := E.Add(p, Q), E.Add(p, NQ)
				newPoints = append(newPoints, pQ, pNQ)
				seen.Add(pQ)
				seen.Add(pNQ)
			}
		}
		points = append(points, newPoints...)
	}

	return points
}
