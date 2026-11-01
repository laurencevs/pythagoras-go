package euclid

import (
	"context"
	"testing"

	"github.com/laurencevs/pythagoras-go/curve"
)

func benchmarkBruteForce(b *testing.B, n int64) {
	for b.Loop() {
		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan curve.RationalPoint)
		go SearchForInitialPoints(ctx, n, ch)
		<-ch
		cancel()
	}
}

func BenchmarkBruteForce5(b *testing.B) {
	benchmarkBruteForce(b, 5)
}

func BenchmarkBruteForce6(b *testing.B) {
	benchmarkBruteForce(b, 6)
}

func BenchmarkBruteForce7(b *testing.B) {
	benchmarkBruteForce(b, 7)
}

func BenchmarkBruteForce13(b *testing.B) {
	benchmarkBruteForce(b, 13)
}

func BenchmarkBruteForce23(b *testing.B) {
	benchmarkBruteForce(b, 23)
}
