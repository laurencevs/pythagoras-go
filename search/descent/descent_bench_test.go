package descent

import (
	"context"
	"testing"

	"github.com/laurencevs/pythagoras-go/curve"
)

func benchmarkDescent(b *testing.B, n int64) {
	for b.Loop() {
		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan curve.RationalPoint)
		go SearchForInitialPoints(ctx, n, ch)
		<-ch
		cancel()
	}
}

func BenchmarkDescent5(b *testing.B) {
	benchmarkDescent(b, 5)
}

func BenchmarkDescent6(b *testing.B) {
	benchmarkDescent(b, 6)
}

func BenchmarkDescent7(b *testing.B) {
	benchmarkDescent(b, 7)
}

func BenchmarkDescent13(b *testing.B) {
	benchmarkDescent(b, 13)
}

func BenchmarkDescent23(b *testing.B) {
	benchmarkDescent(b, 23)
}

func BenchmarkDescent37(b *testing.B) {
	benchmarkDescent(b, 37)
}

func BenchmarkDescent101(b *testing.B) {
	benchmarkDescent(b, 101)
}

func BenchmarkDescent199(b *testing.B) {
	benchmarkDescent(b, 199)
}

func benchmarkDescentWithCleanup(b *testing.B, n int64) {
	for b.Loop() {
		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan curve.RationalPoint)
		go func() {
			<-ch
			cancel()
		}()
		SearchForInitialPoints(ctx, n, ch)
	}
}

func BenchmarkDescent5Clean(b *testing.B) {
	benchmarkDescentWithCleanup(b, 5)
}

func BenchmarkDescent6Clean(b *testing.B) {
	benchmarkDescentWithCleanup(b, 6)
}

func BenchmarkDescent7Clean(b *testing.B) {
	benchmarkDescentWithCleanup(b, 7)
}
