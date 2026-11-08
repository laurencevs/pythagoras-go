package search

import (
	"context"
	"fmt"
	"time"

	"github.com/laurencevs/pythagoras-go/curve"
	"github.com/laurencevs/pythagoras-go/internal/sync"
	"github.com/laurencevs/pythagoras-go/math"
	"github.com/laurencevs/pythagoras-go/pythagoras"
	"github.com/laurencevs/pythagoras-go/search/descent"
	"github.com/laurencevs/pythagoras-go/search/euclid"
	"github.com/laurencevs/pythagoras-go/search/tunnell"
)

func FindRationalTrianglesOfArea(n int, timeout time.Duration, maxCoeff int, verbose bool) []pythagoras.RationalTriple {
	nSquareFree, nSquaredFactor := math.GetSquarefree(uint(n))
	if verbose {
		fmt.Printf("Finding rational Pythagorean triangles with area %d (square-free part %d)\n", n, nSquareFree)
	}

	E := curve.WeierstrassCurve{A: int64(-nSquareFree * nSquareFree), B: 0}

	initialPoints, err := SearchForInitialPoints(int64(nSquareFree), timeout, verbose)
	if err != nil {
		fmt.Printf("Ending search: %v\n", err)
		return nil
	}

	points := curve.MultiplyPoints(E, initialPoints, maxCoeff)

	for _, p := range points {
		if !E.Contains(p) {
			panic(fmt.Sprintf("fatal: result point %v is not on curve %v", p, E))
		}
	}

	uniqueTriples := pythagoras.TriplesFromPoints(points, int(nSquareFree))
	for _, t := range uniqueTriples {
		t.Scale(int(nSquaredFactor))
	}

	return uniqueTriples
}

func SearchForInitialPoints(n int64, timeout time.Duration, verbose bool) ([]curve.RationalPoint, error) {
	returnFirstPoint := timeout == 0

	var searchCtx context.Context
	var endSearch context.CancelFunc

	if timeout > 0 {
		searchCtx, endSearch = context.WithTimeout(context.Background(), timeout)
	} else {
		searchCtx, endSearch = context.WithCancel(context.Background())
	}
	defer endSearch()

	euclidPointsCh := make(chan curve.RationalPoint)
	defer close(euclidPointsCh)
	descentPointsCh := make(chan curve.RationalPoint)
	defer close(descentPointsCh)
	tunnellResultCh := make(chan tunnell.TunnellResult)
	defer close(tunnellResultCh)

	wg := new(sync.WaitGroup)
	wg.Add(
		func() { tunnell.RunTunnellCheck(searchCtx, n, tunnellResultCh) },
		func() { euclid.SearchForInitialPoints(searchCtx, n, euclidPointsCh) },
		func() { descent.SearchForInitialPoints(searchCtx, n, descentPointsCh) },
	)
	wgDone := wg.Go()

	var initialPoints []curve.RationalPoint
	var err error
outer:
	for {
		select {
		case tunnellResult := <-tunnellResultCh:
			if tunnellResult.Valid {
				if verbose {
					fmt.Println("Tunnell's criterion passed, expect some triangles")
				}
			} else {
				endSearch()
				err = fmt.Errorf("Tunnell's criterion failed (2 * %d != %d)", tunnellResult.AC, tunnellResult.BD)
			}
		case point := <-euclidPointsCh:
			initialPoints = append(initialPoints, point)
			if returnFirstPoint {
				endSearch()
			}
			if verbose {
				fmt.Println("Brute force search found initial point", point)
			}
		case point := <-descentPointsCh:
			initialPoints = append(initialPoints, point)
			if returnFirstPoint {
				endSearch()
			}
			if verbose {
				fmt.Println("Descent found initial point", point)
			}
		case <-wgDone:
			break outer
		}
	}

	return initialPoints, err
}
