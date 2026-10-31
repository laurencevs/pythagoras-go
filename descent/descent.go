package descent

import (
	"context"
	"maps"
	"math/big"
	"slices"
	"sync"

	"github.com/laurencevs/pythagoras-go/curve"
	"github.com/laurencevs/pythagoras-go/math"
)

func SearchForInitialPoints(ctx context.Context, n int64, out chan<- curve.RationalPoint) {
	minusNSquaredBases := math.GetPrimeFactors(n)
	minusNSquaredBases[-1] = struct{}{}
	candidateGenerators := getPossibleGenerators(minusNSquaredBases)

	fourNSquaredBases := maps.Clone(minusNSquaredBases)
	delete(fourNSquaredBases, -1)
	fourNSquaredBases[2] = struct{}{}
	isogCandidateGenerators := getPossibleGenerators(fourNSquaredBases)

	// for candidateGenerators:      solve w^2 = g u^4 - (n^2 / g) v^4
	// for isogCandidateGenerators:  solve w^2 = g u^4 + (4 n^2 / g) v^4

	var wg sync.WaitGroup
	wgDone := make(chan struct{})

	resCh := make(chan curve.RationalPoint)
	wg.Add(len(candidateGenerators) + len(isogCandidateGenerators))

	for _, g := range candidateGenerators {
		go func() {
			res := solveQuarticBruteForce(ctx, g, -n*n/g)
			q := res.ToRationalPoint()
			if !q.IsTwoTorsion() {
				resCh <- q
			}
			wg.Done()
		}()
	}
	for _, g := range isogCandidateGenerators {
		go func() {
			res := solveQuarticBruteForce(ctx, g, 4*n*n/g)
			q := res.ToRationalPoint()
			if !q.IsTwoTorsion() {
				resCh <- dualIsogeny(q, n)
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(wgDone)
	}()
	for {
		select {
		case <-wgDone:
			return
		case res := <-resCh:
			out <- res
		}
	}
}

type descentResult struct {
	U, V, W int64
	G       int64
}

func (dr descentResult) ToRationalPoint() curve.RationalPoint {
	if dr.V == 0 {
		return curve.Inf
	}
	x := big.NewRat(dr.G*dr.U*dr.U, dr.V*dr.V)
	y := big.NewRat(abs(dr.G)*dr.U*dr.W, dr.V*dr.V*dr.V)
	return curve.RationalPoint{X: x, Y: y}
}

var (
	fourRat  = big.NewRat(4, 1)
	eightRat = big.NewRat(8, 1)
)

// The dual of the 2-isogeny used for descent
func dualIsogeny(p curve.RationalPoint, n int64) curve.RationalPoint {
	x := new(big.Rat).Quo(p.Y, p.X)
	x.Mul(x, x).Quo(x, fourRat)
	y := new(big.Rat).Mul(p.X, p.X)
	y.Sub(y, big.NewRat(4*n*n, 1)).Mul(y, p.Y)
	y.Quo(y, p.X).Quo(y, p.X).Quo(y, eightRat)
	return curve.RationalPoint{X: x, Y: y}
}

const int64Max = int64(^uint64(0) >> 1)

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// search for solutions to w^2 = a u^4 + b v^4, and post them to out
func solveQuarticBruteForce(ctx context.Context, a, b int64) descentResult {
	// avoid overflow
	lim := int64(math.IntFloorSqrt(math.IntFloorSqrt(uint64(int64Max / (2 * max(abs(a), abs(b)))))))

	for u := int64(0); u < lim; u++ {
		select {
		case <-ctx.Done():
			return descentResult{}
		default:
		}

		vStart := int64(0)
		if u == 0 {
			vStart = 1
		}
		for v := int64(vStart); v < lim; v++ {
			w2 := a*u*u*u*u + b*v*v*v*v
			if w2 < 0 {
				continue
			}
			if w2 == 0 {
				continue
			}
			w, ok := math.IntSqrt(uint64(w2))
			if !ok {
				continue
			}

			return descentResult{
				U: u, V: v, W: int64(w), G: a,
			}
		}
	}
	return descentResult{}
}

func getPossibleGenerators(factors map[int64]struct{}) []int64 {
	res := make([]int64, 0, 1<<len(factors))
	res = append(res, 1)

	for f := range factors {
		resCopy := slices.Clone(res)
		for _, g := range resCopy {
			res = append(res, f*g)
		}
	}

	return res
}
