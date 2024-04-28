package tests

import (
	"fmt"
	"github.com/xformerfhs/blockpad"
	"testing"
)

func BenchmarkUnpadZeroP1(b *testing.B) {
	doUnpad(b, blockpad.Zero, 1)
}

func BenchmarkUnpadZeroM1(b *testing.B) {
	doUnpad(b, blockpad.Zero, testBlockSize-1)
}

func BenchmarkUnpadPKCS7P1(b *testing.B) {
	doUnpad(b, blockpad.PKCS7, 1)
}

func BenchmarkUnpadPKCS7M1(b *testing.B) {
	doUnpad(b, blockpad.PKCS7, testBlockSize-1)
}

func BenchmarkUnpadX923P1(b *testing.B) {
	doUnpad(b, blockpad.X923, 1)
}
func BenchmarkUnpadX923M1(b *testing.B) {
	doUnpad(b, blockpad.X923, testBlockSize-1)
}

func BenchmarkUnpadISO10126P1(b *testing.B) {
	doUnpad(b, blockpad.ISO10126, 1)
}
func BenchmarkUnpadISO10126M1(b *testing.B) {
	doUnpad(b, blockpad.ISO10126, testBlockSize-1)
}

func BenchmarkUnpadRFC4303P1(b *testing.B) {
	doUnpad(b, blockpad.RFC4303, 1)
}

func BenchmarkUnpadRFC4303M1(b *testing.B) {
	doUnpad(b, blockpad.RFC4303, testBlockSize-1)
}

func BenchmarkUnpadISO78164P1(b *testing.B) {
	doUnpad(b, blockpad.ISO78164, 1)
}
func BenchmarkUnpadISO78164M1(b *testing.B) {
	doUnpad(b, blockpad.ISO78164, testBlockSize-1)
}

func BenchmarkUnpadArbitraryTailByteP1(b *testing.B) {
	doUnpad(b, blockpad.ArbitraryTailByte, 1)
}

func BenchmarkUnpadArbitraryTailByteM1(b *testing.B) {
	doUnpad(b, blockpad.ArbitraryTailByte, testBlockSize-1)
}

// ******** Private function ********

func doUnpad(b *testing.B, padAlgorithm blockpad.PadAlgorithm, offset int) {
	b.StopTimer()
	testLen := dataLen + offset
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(padAlgorithm, testBlockSize)
	if err != nil {
		b.Fatalf(`Error creating padder: %v`, err)
	}
	paddedData := padder.Pad(data)
	b.StartTimer()

	var unpaddedData []byte
	for i := 0; i < b.N; i++ {
		unpaddedData, _ = padder.Unpad(paddedData)
	}
	b.StopTimer()

	fmt.Println(len(unpaddedData))
}
