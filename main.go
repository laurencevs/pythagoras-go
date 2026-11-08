package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"

	"github.com/laurencevs/pythagoras-go/pythagoras"
	"github.com/laurencevs/pythagoras-go/search"
)

var timeoutDuration = flag.Duration("timeout", 0, "The timeout length for the initial point search. If set to 0, the program will keep searching until it finds one point.")
var maxCoeff = flag.Int("max-coeff", 5, "The maximum coefficient by which to multiply each point.")
var maxDigits = flag.Int("max-digits", 0, "The maximum number of digits in the denominator of the hypotenuse. Triangles with larger hypotenuse denominators will not be outputted.")
var outputFile = flag.String("output", "", "The filename for output. By default, the program will print to stdout.")
var verbose = flag.Bool("verbose", false, "Whether to include verbose output.")

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	n, err := strconv.Atoi(flag.Args()[0])
	if err != nil {
		panic(err)
	}

	triples := search.FindRationalTrianglesOfArea(n, *timeoutDuration, *maxCoeff, *verbose)

	if *maxDigits != 0 {
		triples = pythagoras.FilterByMaxDigits(triples, *maxDigits)
	}

	if len(triples) == 0 {
		return
	}

	var b bytes.Buffer
	for _, t := range triples {
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
