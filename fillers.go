//
// SPDX-FileCopyrightText: Copyright 2024 Frank Schwab
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileType: SOURCE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// You may not use this file except in compliance with the License.
//
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Author: Frank Schwab
//
// Version: 1.0.0
//
// Change history:
//    2024-04-19: V1.0.0: Created.
//

package blockpad

import (
	"crypto/rand"
	mrand "math/rand"
	"padding/internal/slicehelper"
)

// ******** This file contains the private padding fillers ********

// zeroFiller creates a filler with all zeroes.
// This filler panics if the clear data ends with a 0 byte in the last block.
func zeroFiller(lastByte byte, dataLen int, blockSize int) []byte {
	padLen, _, pad := makePadSlice(dataLen, blockSize)
	if lastByte == 0 && padLen != blockSize {
		panic(`last byte must not be 0`)
	}

	return pad
}

// pkcs7Filler creates a filler with all bytes containing the length of the filler.
func pkcs7Filler(lastByte byte, dataLen int, blockSize int) []byte {
	_, padLenByte, pad := makePadSlice(dataLen, blockSize)
	slicehelper.Fill(pad, padLenByte)
	return pad
}

// x923Filler contains a filler where the last byte contains the length and all other bytes are zero.
func x923Filler(lastByte byte, dataLen int, blockSize int) []byte {
	padLen, padLenByte, pad := makePadSlice(dataLen, blockSize)
	pad[padLen-1] = padLenByte
	return pad
}

// iso10126Filler contains a filler where the last byte contains the length and all other bytes have random values.
func iso10126Filler(lastByte byte, dataLen int, blockSize int) []byte {
	padLen, padLenByte, pad := makePadSlice(dataLen, blockSize)
	lastIndex := padLen - 1
	pad[lastIndex] = padLenByte
	_, _ = rand.Read(pad[:lastIndex])
	return pad
}

// rfc4303Filler contains a filler where the last byte contains the length and the other bytes are counted down from right to left.
func rfc4303Filler(lastByte byte, dataLen int, blockSize int) []byte {
	_, padLenByte, pad := makePadSlice(dataLen, blockSize)
	for actLen := padLenByte; actLen > 0; actLen-- {
		pad[actLen-1] = actLen
	}
	return pad
}

// iso78164Filler contains a filler where the first byte contains the value 0x80 and all other bytes are zero.
func iso78164Filler(lastByte byte, dataLen int, blockSize int) []byte {
	_, _, pad := makePadSlice(dataLen, blockSize)
	pad[0] = 0x80
	return pad
}

// arbitraryTailByteFiller contains a filler where all bytes contain the same random value which is not the value of the last data byte.
// This is the *only* padding that is *not* susceptible to a padding oracle!
func arbitraryTailByteFiller(lastByte byte, dataLen int, blockSize int) []byte {
	_, _, pad := makePadSlice(dataLen, blockSize)
	slicehelper.Fill(pad, anythingBut(lastByte))
	return pad
}

// -------- Helper functions --------

// makePadSlice makes a slice that has the right length to fill the data to a multiple of the block size.
// If the data length is already a multiple of the block size a full block is returned.
func makePadSlice(dataLen int, blockSize int) (int, byte, []byte) {
	padLen := blockSize - dataLen%blockSize
	pad := make([]byte, padLen)
	return padLen, byte(padLen), pad
}

// anythingBut returns a byte value that is not the same value as the argument.
func anythingBut(b byte) byte {
	var r byte
	for {
		r = byte(mrand.Int31n(256))
		if r != b {
			break
		}
	}
	return r
}
