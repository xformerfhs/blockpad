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

package blockpad

import (
	"crypto/rand"
	"github.com/xformerfhs/blockpad/internal/slicehelper"
	mrand "math/rand"
)

// ******** This file contains the private padding fillers ********

// zeroFiller creates a filler with all zeroes.
// This filler panics if the clear data ends with a 0 byte in the last block.
func zeroFiller(lastBlock []byte, blockSize int, lastData []byte, lastBlockDataLen int, padLen int) {
	if lastBlockDataLen > 0 && lastData[lastBlockDataLen-1] == 0 {
		panic(`last data byte must not be 0`)
	}
}

// pkcs7Filler creates a filler with all bytes containing the length of the filler.
func pkcs7Filler(lastBlock []byte, blockSize int, lastData []byte, lastBlockDataLen int, padLen int) {
	slicehelper.Fill(lastBlock, byte(padLen))
}

// x923Filler contains a filler where the last byte contains the length and all other bytes are zero.
func x923Filler(lastBlock []byte, blockSize int, lastData []byte, lastBlockDataLen int, padLen int) {
	lastBlock[blockSize-1] = byte(padLen)
}

// iso10126Filler contains a filler where the last byte contains the length and all other bytes have random values.
func iso10126Filler(lastBlock []byte, blockSize int, lastData []byte, lastBlockDataLen int, padLen int) {
	_, _ = rand.Read(lastBlock)
	lastBlock[blockSize-1] = byte(padLen)
}

// rfc4303Filler contains a filler where the last byte contains the length and the other bytes are counted down from right to left.
func rfc4303Filler(lastBlock []byte, blockSize int, lastData []byte, lastBlockDataLen int, padLen int) {
	padByte := byte(padLen)
	for i := blockSize - 1; i >= 0; i-- {
		lastBlock[i] = padByte
		padByte--
	}
}

// iso78164Filler contains a filler where the first byte contains the value 0x80 and all other bytes are zero.
func iso78164Filler(lastBlock []byte, blockSize int, lastData []byte, lastBlockDataLen int, padLen int) {
	lastBlock[lastBlockDataLen] = 0x80
}

// arbitraryTailByteFiller contains a filler where all bytes contain the same random value which is not the value of the last data byte.
// This is the *only* padding that is *not* susceptible to a padding oracle!
func arbitraryTailByteFiller(lastBlock []byte, blockSize int, lastData []byte, lastBlockDataLen int, padLen int) {
	fillByte := getFillByte(lastData, lastBlockDataLen)
	slicehelper.Fill(lastBlock, fillByte)
}

// -------- Helper functions --------

// getFillByte gets the byte that is used for padding with arbitrary tail byte padding.
func getFillByte(lastData []byte, lastBlockDataLen int) byte {
	var result byte

	if lastBlockDataLen != 0 {
		result = anythingBut(lastData[lastBlockDataLen-1])
	} else {
		result = byte(mrand.Int31n(256)) // Just use any byte value if the last block is padding-only.
	}

	return result
}

// anythingBut returns a byte value that is not the same value as the argument.
func anythingBut(notThisByte byte) byte {
	var result byte

	for {
		result = byte(mrand.Int31n(256))

		if result != notThisByte {
			break
		}
	}

	return result
}
