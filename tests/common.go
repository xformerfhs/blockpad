package tests

import "crypto/rand"

// dataLen is the minimum data length for tests with random data lengths.
const dataLen = 128

// testBlockSize is the block size used in the tests.
const testBlockSize = 16

// makeTestSlice creates a test slice with the specified length.
func makeTestSlice(l int) []byte {
	data := make([]byte, l)
	_, _ = rand.Read(data)

	data[l-1] = 0x5a // For Zero padding.

	return data
}
