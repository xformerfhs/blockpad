package blockpad_test

import (
	"crypto/rand"
	"fmt"
	"github.com/xformerfhs/blockpad"
	"testing"
)

// dataLen is the minimum data length for tests with random data lengths.
const dataLen = 128

// testBlockSize is the block size used in the tests.
const testBlockSize = 16

func BenchmarkPadP1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.PKCS7, testBlockSize)
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

func BenchmarkPadM1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + testBlockSize - 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.PKCS7, testBlockSize)
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

func BenchmarkUnpadZeroP1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.Zero, testBlockSize)
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
func BenchmarkUnpadZeroM1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + testBlockSize - 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.Zero, testBlockSize)
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

func BenchmarkUnpadPKCS7P1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.PKCS7, testBlockSize)
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
func BenchmarkUnpadPKCS7M1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + testBlockSize - 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.PKCS7, testBlockSize)
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

func BenchmarkUnpadX923P1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.X923, testBlockSize)
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
func BenchmarkUnpadX923M1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + testBlockSize - 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.X923, testBlockSize)
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

func BenchmarkUnpadISO10126P1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.ISO10126, testBlockSize)
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
func BenchmarkUnpadISO10126M1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + testBlockSize - 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.ISO10126, testBlockSize)
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

func BenchmarkUnpadRFC4303P1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.RFC4303, testBlockSize)
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
func BenchmarkUnpadRFC4303M1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + testBlockSize - 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.RFC4303, testBlockSize)
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

func BenchmarkUnpadISO78164P1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.ISO78164, testBlockSize)
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
func BenchmarkUnpadISO78164M1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + testBlockSize - 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.ISO78164, testBlockSize)
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

func BenchmarkUnpadArbitraryTailByteP1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.ArbitraryTailByte, testBlockSize)
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
func BenchmarkUnpadArbitraryTailByteM1(b *testing.B) {
	b.StopTimer()
	testLen := dataLen + testBlockSize - 1
	data := makeTestSlice(testLen)
	padder, err := blockpad.NewBlockPadding(blockpad.ArbitraryTailByte, testBlockSize)
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

// ******** Private functions ********

// makeTestSlice creates a test slice with the specified length.
func makeTestSlice(l int) []byte {
	data := make([]byte, l)
	_, _ = rand.Read(data)

	return data
}
