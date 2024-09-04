package benchmarks

import (
	"github.com/xformerfhs/blockpad"
	"runtime"
	"testing"
)

func BenchmarkPadLastBlockZeroLong(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.Zero, 1)
}

func BenchmarkPadLastBlockZeroShort(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.Zero, testBlockSize-1)
}

func BenchmarkPadLastBlockPKCS7Long(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.PKCS7, 1)
}

func BenchmarkPadLastBlockPKCS7Short(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.PKCS7, testBlockSize-1)
}

func BenchmarkPadLastBlockX923Long(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.X923, 1)
}

func BenchmarkPadLastBlockX923Short(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.X923, testBlockSize-1)
}

func BenchmarkPadLastBlockISO10126Long(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.ISO10126, 1)
}

func BenchmarkPadLastBlockISO10126Short(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.ISO10126, testBlockSize-1)
}

func BenchmarkPadLastBlockRFC4303Long(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.RFC4303, 1)
}

func BenchmarkPadLastBlockRFC4303Short(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.RFC4303, testBlockSize-1)
}

func BenchmarkPadLastBlockISO78164Long(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.ISO78164, 1)
}

func BenchmarkPadLastBlockISO78164Short(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.ISO78164, testBlockSize-1)
}

func BenchmarkPadLastBlockArbitraryTailByteLong(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.ArbitraryTailByte, 1)
}

func BenchmarkPadLastBlockArbitraryTailByteShort(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.ArbitraryTailByte, testBlockSize-1)
}

func BenchmarkPadLastBlockNotLastByteLong(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.NotLastByte, 1)
}

func BenchmarkPadLastBlockNotLastByteShort(b *testing.B) {
	b.StopTimer()
	doBenchPadLastBlock(b, blockpad.NotLastByte, testBlockSize-1)
}

// ******** Private function ********

// doBenchPadLastBlock runs a PadLastBlock benchmark with the given parameters.
func doBenchPadLastBlock(b *testing.B, padAlgorithm blockpad.PadAlgorithm, additionalDataLen int) {
	testLen := minimumDataLen + additionalDataLen
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(padAlgorithm, testBlockSize)
	if err != nil {
		b.Fatalf(`Error creating padder: %v`, err)
	}

	runtime.GC()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = padder.PadLastBlock(data)
	}
	b.StopTimer()
}
