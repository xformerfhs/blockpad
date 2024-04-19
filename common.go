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

import "math"

// ******** This file contains private data and utility functions ********

// ******* Private constants ********

// padImplementation holds the implementation information for the various padding algorithms.
var padImplementation = []implementationInfo{
	{name: `Zero`, filler: zeroFiller, remover: zeroRemover},
	{name: `PKCS#7`, filler: pkcs7Filler, remover: pkcs7Remover},
	{name: `X.923`, filler: x923Filler, remover: x923Remover},
	{name: `ISO 10126`, filler: iso10126Filler, remover: iso10126Remover},
	{name: `RFC 4303`, filler: rfc4303Filler, remover: rfc4303Remover},
	{name: `ISO 7816-4`, filler: iso78164Filler, remover: iso78164Remover},
	{name: `Arbitrary Tail Byte`, filler: arbitraryTailByteFiller, remover: arbitraryTailBytePaddingRemover},
}

// ******** Private functions ********

// getLenAndLastByteAndCheckBlockSize gets the length and the last byte of the data and
// checks if the block size is valid.
func getLenAndLastByteAndCheckBlockSize(blockSize int, data []byte) (int, byte, error) {
	err := checkBlockSize(blockSize)
	if err != nil {
		return 0, 0, err
	}

	dataLen, lastByte := getLenAndLastByte(data)
	return dataLen, lastByte, nil
}

// getLenAndLastByte gets the length and the last byte of the data.
func getLenAndLastByte(data []byte) (int, byte) {
	dataLen := len(data)
	var lastByte byte

	if dataLen != 0 {
		lastByte = data[dataLen-1]
	} else {
		// This can be any value, except 0, because this is an invalid value for zero padding.
		// It is not used for anything.
		lastByte = 0xff
	}

	return dataLen, lastByte
}

// checkPadAlgorithmAndBlockSize checks if the pad algorithm and the block size are valid.
func checkPadAlgorithmAndBlockSize(padAlgorithm PadAlgorithm, blockSize int) error {
	err := checkBlockSize(blockSize)
	if err != nil {
		return err
	}

	if padAlgorithm > maxAlgorithm {
		return ErrInvalidPadAlgorithm
	}

	return nil
}

// checkBlockSize checks if the block size is valid.
func checkBlockSize(blockSize int) error {
	if blockSize < 1 || blockSize > math.MaxUint8 {
		return ErrInvalidBlockSize
	}

	return nil
}