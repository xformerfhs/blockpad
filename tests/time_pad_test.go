package tests

import (
	"fmt"
	"github.com/xformerfhs/blockpad"
	"testing"
)

func BenchmarkPadZeroP1(b *testing.B) {
	doBenchPad(b, blockpad.Zero, 1)
}

func BenchmarkPadZeroM1(b *testing.B) {
	doBenchPad(b, blockpad.Zero, testBlockSize-1)
}

func BenchmarkPadPKCS7P1(b *testing.B) {
	doBenchPad(b, blockpad.PKCS7, 1)
}

func BenchmarkPadPKCS7M1(b *testing.B) {
	doBenchPad(b, blockpad.PKCS7, testBlockSize-1)
}

func BenchmarkPadX923P1(b *testing.B) {
	doBenchPad(b, blockpad.X923, 1)
}

func BenchmarkPadX923M1(b *testing.B) {
	doBenchPad(b, blockpad.X923, testBlockSize-1)
}

func BenchmarkPadISO10126P1(b *testing.B) {
	doBenchPad(b, blockpad.ISO10126, 1)
}

func BenchmarkPadISO10126M1(b *testing.B) {
	doBenchPad(b, blockpad.ISO10126, testBlockSize-1)
}

func BenchmarkPadRFC4303P1(b *testing.B) {
	doBenchPad(b, blockpad.RFC4303, 1)
}

func BenchmarkPadRFC4303M1(b *testing.B) {
	doBenchPad(b, blockpad.RFC4303, testBlockSize-1)
}

func BenchmarkPadISO78164P1(b *testing.B) {
	doBenchPad(b, blockpad.ISO78164, 1)
}

func BenchmarkPadISO78164M1(b *testing.B) {
	doBenchPad(b, blockpad.ISO78164, testBlockSize-1)
}

func BenchmarkPadArbitraryTailByteP1(b *testing.B) {
	doBenchPad(b, blockpad.ArbitraryTailByte, 1)
}

func BenchmarkPadArbitraryTailByteM1(b *testing.B) {
	doBenchPad(b, blockpad.ArbitraryTailByte, testBlockSize-1)
}

// ******** Private function ********

// doBenchPad runs a Pad benchmark with the given parameters.
func doBenchPad(b *testing.B, padAlgorithm blockpad.PadAlgorithm, offset int) {
	b.StopTimer()
	testLen := dataLen + offset
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(padAlgorithm, testBlockSize)
	if err != nil {
		b.Fatalf(`Error creating padder: %v`, err)
	}
	b.StartTimer()

	var paddedData []byte
	for i := 0; i < b.N; i++ {
		paddedData = padder.Pad(data)
	}
	b.StopTimer()

	fmt.Println(len(paddedData))
}
