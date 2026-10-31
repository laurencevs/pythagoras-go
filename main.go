package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/laurencevs/pythagoras-go/curve"
	"github.com/laurencevs/pythagoras-go/internal/sync"
	"github.com/laurencevs/pythagoras-go/math"
	"github.com/laurencevs/pythagoras-go/pythagoras"
	"github.com/laurencevs/pythagoras-go/search/descent"
	"github.com/laurencevs/pythagoras-go/search/euclid"
	"github.com/laurencevs/pythagoras-go/tunnell"
)

var timeoutDuration = flag.Duration("timeout", 0, "The timeout length for the initial point search. If set to 0, the program will keep searching until it finds one point.")
var maxCoeff = flag.Int("max-coeff", 5, "The maximum coefficient by which to multiply each point.")
var maxDigits = flag.Int("max-digits", 0, "The maximum number of digits in the denominator of the hypotenuse. Triangles with larger hypotenuse denominators will not be outputted.")
var outputFile = flag.String("output", "", "The filename for output. By default, the program will print to stdout.")

func main() {
	flag.Parse()

	n, err := strconv.Atoi(flag.Args()[0])
	if err != nil {
		panic(err)
	}

	uniqueTriples := findRationalTrianglesOfArea(n, *timeoutDuration, *maxCoeff)
	if len(uniqueTriples) == 0 {
		return
	}

	if *maxDigits != 0 {
		var filteredTriples = make([]pythagoras.RationalTriple, 0, len(uniqueTriples))
		for _, t := range uniqueTriples {
			if (*maxDigits) > 0 && len(t.C.Denom().String()) > (*maxDigits) {
				continue
			}
			if (*maxDigits) == -1 && (!t.C.Denom().IsInt64() || t.C.Denom().Int64() != 1) {
				continue
			}
			filteredTriples = append(filteredTriples, t)
		}
		uniqueTriples = filteredTriples
	}

	var b bytes.Buffer
	for _, t := range uniqueTriples {
		b.WriteString(t.String())
		b.WriteByte('\n')
	}

	if *outputFile != "" {
		err := os.WriteFile(*outputFile, b.Bytes(), 0666)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Print(b.String())
	}
}

type valuesSeen[T fmt.Stringer] map[string]struct{}

func (ps valuesSeen[T]) add(p T) {
	ps[p.String()] = struct{}{}
}

func (ps valuesSeen[T]) contains(p T) bool {
	_, ok := ps[p.String()]
	return ok
}

func findRationalTrianglesOfArea(n int, timeout time.Duration, maxCoeff int) []pythagoras.RationalTriple {
	nSquareFree, nSquaredFactor := math.GetSquarefree(n)
	fmt.Printf("Finding rational Pythagorean triangles with area %d (square-free part %d)\n", n, nSquareFree)

	E := curve.WeierstrassCurve{A: int64(-nSquareFree * nSquareFree), B: 0}

	initialPoints := runInitialSearch(int64(nSquareFree), timeout)
	if len(initialPoints) == 0 {
		return nil
	}

	for _, p := range initialPoints {
		if !E.Contains(p) {
			panic(fmt.Sprintf("fatal: initial point %v is not on curve %v", p, E))
		}
	}

	points := multiplyAndDeduplicatePoints(E, initialPoints, maxCoeff)

	for _, p := range points {
		if p.IsTwoTorsion() {
			continue
		}
		if !E.Contains(p) {
			panic(fmt.Sprintf("fatal: result point %v is not on curve %v", p, E))
		}
	}

	// Convert points to Pythagorean triples, again with de-duplication.
	// Multiple points can correspond to the same triple.
	triplesSeen := make(valuesSeen[pythagoras.RationalTriple])
	var uniqueTriples []pythagoras.RationalTriple
	for _, P := range points {
		if P.IsTwoTorsion() {
			continue
		}
		t := pythagoras.TripleFromPoint(P, nSquareFree)
		t.Scale(nSquaredFactor)
		if !triplesSeen.contains(t) {
			uniqueTriples = append(uniqueTriples, t)
		}
		triplesSeen.add(t)
	}

	return uniqueTriples
}

func multiplyAndDeduplicatePoints(E curve.WeierstrassCurve, initialPoints []curve.RationalPoint, maxCoeff int) []curve.RationalPoint {
	// Multiply the initial points and compute all linear combinations thereof,
	// recording seen values for de-duplication.
	var points []curve.RationalPoint
	seen := make(valuesSeen[curve.RationalPoint])
	for _, c := range initialPoints {
		if seen.contains(c) || seen.contains(E.Invert(c)) {
			continue
		}
		var newPoints []curve.RationalPoint
		Q := curve.Inf
		for i := 0; i < maxCoeff; i++ {
			Q = E.Add(Q, c)
			NQ := E.Invert(Q)
			newPoints = append(newPoints, Q, NQ)
			seen.add(Q)
			seen.add(NQ)
			for _, p := range points {
				pQ, pNQ := E.Add(p, Q), E.Add(p, NQ)
				newPoints = append(newPoints, pQ, pNQ)
				seen.add(pQ)
				seen.add(pNQ)
			}
		}
		points = append(points, newPoints...)
	}

	return points
}

func runInitialSearch(n int64, timeout time.Duration) []curve.RationalPoint {
	returnFirstPoint := timeout == 0

	var searchCtx context.Context
	var endSearch context.CancelFunc

	if timeout > 0 {
		searchCtx, endSearch = context.WithTimeout(context.Background(), timeout)
	} else {
		searchCtx, endSearch = context.WithCancel(context.Background())
	}
	defer endSearch()

	initialPointsCh := make(chan curve.RationalPoint)
	defer close(initialPointsCh)
	descentPointsCh := make(chan curve.RationalPoint)
	defer close(descentPointsCh)
	tunnellResultCh := make(chan tunnell.TunnellResult)
	defer close(tunnellResultCh)

	wg := new(sync.WaitGroup)
	wg.Add(
		func() { tunnell.RunTunnellCheck(searchCtx, n, tunnellResultCh) },
		func() { euclid.SearchForInitialPoints(searchCtx, n, initialPointsCh) },
		func() { descent.SearchForInitialPoints(searchCtx, n, descentPointsCh) },
	)
	wgDone := wg.Go()

	var initialPoints []curve.RationalPoint
outer:
	for {
		select {
		case point := <-initialPointsCh:
			initialPoints = append(initialPoints, point)
			if returnFirstPoint {
				endSearch()
				fmt.Println("Brute force search found initial point")
			}
		case point := <-descentPointsCh:
			initialPoints = append(initialPoints, point)
			if returnFirstPoint {
				endSearch()
				fmt.Println("Descent found initial point")
			}
		case tunnellResult := <-tunnellResultCh:
			if tunnellResult.Valid {
				fmt.Println("Tunnell's criterion passed, expect some triangles")
			} else {
				endSearch()
				fmt.Printf("Tunnell's criterion failed (counts %d, %d), ending search\n", tunnellResult.AC, tunnellResult.BD)
			}
		case <-wgDone:
			break outer
		}
	}

	return initialPoints
}
