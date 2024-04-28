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
