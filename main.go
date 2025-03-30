package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
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

	fmt.Printf("Finding rational Pythagorean triangles with area %d\n", n)

	var timeout <-chan time.Time
	if *timeoutDuration > 0 {
		timeout = time.After(*timeoutDuration)
	}
	initialPoints := initialPointSearch(n, timeout)

	E := weierstrassCurve{int64(-n * n), 0}

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
		t := P.ToTriple(n)
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
