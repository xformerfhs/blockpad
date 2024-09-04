package benchmarks

import (
	"github.com/xformerfhs/blockpad"
	"runtime"
	"testing"
)

// ******** Pad benchmarks ********

func BenchmarkPadZeroLong(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.Zero, 1)
}

func BenchmarkPadZeroShort(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.Zero, testBlockSize-1)
}

func BenchmarkPadPKCS7Long(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.PKCS7, 1)
}

func BenchmarkPadPKCS7Short(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.PKCS7, testBlockSize-1)
}

func BenchmarkPadX923Long(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.X923, 1)
}

func BenchmarkPadX923Short(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.X923, testBlockSize-1)
}

func BenchmarkPadISO10126Long(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.ISO10126, 1)
}

func BenchmarkPadISO10126Short(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.ISO10126, testBlockSize-1)
}

func BenchmarkPadRFC4303Long(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.RFC4303, 1)
}

func BenchmarkPadRFC4303Short(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.RFC4303, testBlockSize-1)
}

func BenchmarkPadISO78164Long(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.ISO78164, 1)
}

func BenchmarkPadISO78164Short(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.ISO78164, testBlockSize-1)
}

func BenchmarkPadArbitraryTailByteLong(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.ArbitraryTailByte, 1)
}

func BenchmarkPadArbitraryTailByteShort(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.ArbitraryTailByte, testBlockSize-1)
}

func BenchmarkPadNotLastByteLong(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.NotLastByte, 1)
}

func BenchmarkPadNotLastByteShort(b *testing.B) {
	b.StopTimer()
	doBenchPad(b, blockpad.ArbitraryTailByte, testBlockSize-1)
}

// ******** Private function ********

// doBenchPad runs a Pad benchmark with the given parameters.
func doBenchPad(b *testing.B, padAlgorithm blockpad.PadAlgorithm, unpaddedDataLen int) {
	data := makeTestSlice(unpaddedDataLen)
	padder, err := blockpad.NewBlockPadding(padAlgorithm, testBlockSize)
	if err != nil {
		b.Fatalf(`Error creating padder: %v`, err)
	}

	runtime.GC()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = padder.Pad(data)
	}
	b.StopTimer()
}
