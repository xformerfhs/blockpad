package benchmarks

import (
	"github.com/xformerfhs/blockpad"
	"runtime"
	"testing"
)

func BenchmarkUnpadZeroLong(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.Zero, 1)
}

func BenchmarkUnpadZeroShort(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.Zero, testBlockSize-1)
}

func BenchmarkUnpadPKCS7Long(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.PKCS7, 1)
}

func BenchmarkUnpadPKCS7Short(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.PKCS7, testBlockSize-1)
}

func BenchmarkUnpadX923Long(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.X923, 1)
}
func BenchmarkUnpadX923Short(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.X923, testBlockSize-1)
}

func BenchmarkUnpadISO10126Long(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.ISO10126, 1)
}

func BenchmarkUnpadISO10126Short(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.ISO10126, testBlockSize-1)
}

func BenchmarkUnpadRFC4303Long(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.RFC4303, 1)
}

func BenchmarkUnpadRFC4303Short(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.RFC4303, testBlockSize-1)
}

func BenchmarkUnpadISO78164Long(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.ISO78164, 1)
}
func BenchmarkUnpadISO78164Short(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.ISO78164, testBlockSize-1)
}

func BenchmarkUnpadArbitraryTailByteLong(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.ArbitraryTailByte, 1)
}

func BenchmarkUnpadArbitraryTailByteShort(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.ArbitraryTailByte, testBlockSize-1)
}

func BenchmarkUnpadNotLastByteLong(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.NotLastByte, 1)
}

func BenchmarkUnpadNotLastByteShort(b *testing.B) {
	b.StopTimer()
	doUnpad(b, blockpad.ArbitraryTailByte, testBlockSize-1)
}

// ******** Private function ********

func doUnpad(b *testing.B, padAlgorithm blockpad.PadAlgorithm, unpaddedDataLen int) {
	testLen := unpaddedDataLen
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(padAlgorithm, testBlockSize)
	if err != nil {
		b.Fatalf(`Error creating padder: %v`, err)
	}

	paddedData := padder.Pad(data)

	runtime.GC()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = padder.Unpad(paddedData)
	}
	b.StopTimer()
}
