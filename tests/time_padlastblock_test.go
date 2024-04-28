package tests

import (
	"fmt"
	"github.com/xformerfhs/blockpad"
	"testing"
)

func BenchmarkPadLastBlockZeroP1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.Zero, 1)
}

func BenchmarkPadLastBlockZeroM1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.Zero, testBlockSize-1)
}

func BenchmarkPadLastBlockPKCS7P1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.PKCS7, 1)
}

func BenchmarkPadLastBlockPKCS7M1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.PKCS7, testBlockSize-1)
}

func BenchmarkPadLastBlockX923P1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.X923, 1)
}

func BenchmarkPadLastBlockX923M1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.X923, testBlockSize-1)
}

func BenchmarkPadLastBlockISO10126P1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.ISO10126, 1)
}

func BenchmarkPadLastBlockISO10126M1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.ISO10126, testBlockSize-1)
}

func BenchmarkPadLastBlockRFC4303P1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.RFC4303, 1)
}

func BenchmarkPadLastBlockRFC4303M1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.RFC4303, testBlockSize-1)
}

func BenchmarkPadLastBlockISO78164P1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.ISO78164, 1)
}

func BenchmarkPadLastBlockISO78164M1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.ISO78164, testBlockSize-1)
}

func BenchmarkPadLastBlockArbitraryTailByteP1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.ArbitraryTailByte, 1)
}

func BenchmarkPadLastBlockArbitraryTailByteM1(b *testing.B) {
	doBenchPadLastBlock(b, blockpad.ArbitraryTailByte, testBlockSize-1)
}

// ******** Private function ********

// doBenchPadLastBlock runs a PadLastBlock benchmark with the given parameters.
func doBenchPadLastBlock(b *testing.B, padAlgorithm blockpad.PadAlgorithm, offset int) {
	b.StopTimer()
	testLen := dataLen + offset
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(padAlgorithm, testBlockSize)
	if err != nil {
		b.Fatalf(`Error creating padder: %v`, err)
	}
	b.StartTimer()

	var paddedFullBlock []byte
	var paddedLastBlock []byte
	for i := 0; i < b.N; i++ {
		paddedFullBlock, paddedLastBlock = padder.PadLastBlock(data)
	}
	b.StopTimer()

	fmt.Println(len(paddedFullBlock) + len(paddedLastBlock))
}
