package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/laurencevs/pythagoras-go/math"
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
	nSquareFree, nSquaredFactor := math.GetSquarefree(n)

	fmt.Printf("Finding rational Pythagorean triangles with area %d (square-free part %d)\n", n, nSquareFree)

	initialPoints := runInitialSearch(int64(nSquareFree), *timeoutDuration)
	if len(initialPoints) == 0 {
		return
	}

	E := weierstrassCurve{int64(-nSquareFree * nSquareFree), 0}

	for _, p := range initialPoints {
		if !E.Contains(p) {
			panic(fmt.Sprintf("fatal: initial point %v is not on curve %v", p, E))
		}
	}

	// Multiply the initial points and compute all linear combinations thereof,
	// recording seen values for de-duplication.
	var points []rationalPoint
	seen := make(valuesSeen[rationalPoint])
	for _, c := range initialPoints {
		if seen.contains(c) || seen.contains(E.Invert(c)) {
			continue
		}
		var newPoints []rationalPoint
		Q := inf
		for i := 0; i < *maxCoeff; i++ {
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

	// Convert points to Pythagorean triples, again with de-duplication.
	// Multiple points can correspond to the same triple.
	triplesSeen := make(valuesSeen[rationalTriple])
	var uniqueTriples []rationalTriple
	for _, P := range points {
		t := P.ToTriple(nSquareFree)
		t.Scale(nSquaredFactor)
		if !triplesSeen.contains(t) {
			uniqueTriples = append(uniqueTriples, t)
		}
		triplesSeen.add(t)
	}

	if *maxDigits != 0 {
		var filteredTriples = make([]rationalTriple, 0, len(uniqueTriples))
		for _, t := range uniqueTriples {
			if len(t.C.Denom().String()) > (*maxDigits) {
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
		err = os.WriteFile(*outputFile, b.Bytes(), 0666)
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

func runInitialSearch(n int64, timeout time.Duration) []rationalPoint {
	returnFirstPoint := timeout == 0

	var searchCtx context.Context
	var endSearch context.CancelFunc

	if timeout > 0 {
		searchCtx, endSearch = context.WithTimeout(context.Background(), timeout)
	} else {
		searchCtx, endSearch = context.WithCancel(context.Background())
	}
	defer endSearch()

	initialPointsCh := make(chan rationalPoint)
	defer close(initialPointsCh)
	tunnellResultCh := make(chan tunnell.TunnellResult)
	defer close(tunnellResultCh)

	var wg sync.WaitGroup
	wg.Add(2)
	wgDone := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgDone)
	}()
	go func() {
		initialPointSearch(searchCtx, n, initialPointsCh)
		wg.Done()
	}()
	go func() {
		tunnell.RunTunnellCheck(context.TODO(), n, tunnellResultCh)
		wg.Done()
	}()

	var initialPoints []rationalPoint
outer:
	for {
		select {
		case point := <-initialPointsCh:
			initialPoints = append(initialPoints, point)
			if returnFirstPoint {
				endSearch()
			}
		case tunnellResult := <-tunnellResultCh:
			if tunnellResult.Valid {
				fmt.Println("Tunnell's criterion passed, expect some triangles")
			} else {
				fmt.Printf("Tunnell's criterion failed (counts %d, %d), ending search\n", tunnellResult.AC, tunnellResult.BD)
				endSearch()
			}
		case <-wgDone:
			break outer
		}
	}

	return initialPoints
}
