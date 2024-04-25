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

// ******** This file contains the private padding removers ********

// zeroRemover removes zero padding (ISO 10118-1 and ISO 9797-1).
// It may return an ErrInvalidPadding error and is therefore susceptible to a padding oracle!
func zeroRemover(data []byte, dataLen int, blockSize int) ([]byte, error) {
	lastIndex := dataLen - 1
	if data[lastIndex] != 0 {
		return nil, ErrInvalidPadding
	}

	firstIndex := dataLen - blockSize

	firstPadIndex := firstIndex
	// Always scan *all* data of the last block to thwart timing attacks.
	for i := lastIndex - 1; i >= firstIndex; i-- {
		if data[i] != 0 && firstPadIndex == firstIndex {
			firstPadIndex = i + 1
		}
	}

	return data[:firstPadIndex], nil
}

// pkcs7Remover removes PKCS#7 padding (RFC 5652).
// It may return an ErrInvalidPadding error and is therefore susceptible to a padding oracle!
func pkcs7Remover(data []byte, dataLen int, blockSize int) ([]byte, error) {
	firstIndex, lastIndex, firstPadIndex, padLenByte, _, err := checkLengthByte(data, dataLen, blockSize)
	if err != nil {
		return nil, err
	}

	isValid := true
	// Always scan *all* data of the last block to thwart timing attacks.
	for i := lastIndex - 1; i >= firstIndex; i-- {
		isValid = isValid && (data[i] == padLenByte || i < firstPadIndex)
	}
	if !isValid {
		return nil, ErrInvalidPadding
	}

	return data[:firstPadIndex], nil
}

// x923Remover removes ANSI X.923 padding.
// It may return an ErrInvalidPadding error and is therefore susceptible to a padding oracle!
func x923Remover(data []byte, dataLen int, blockSize int) ([]byte, error) {
	firstIndex, lastIndex, firstPadIndex, _, _, err := checkLengthByte(data, dataLen, blockSize)
	if err != nil {
		return nil, err
	}

	isValid := true
	// Always scan *all* data of the last block to thwart timing attacks.
	for i := lastIndex - 1; i >= firstIndex; i-- {
		isValid = isValid && (data[i] == 0 || i < firstPadIndex)
	}
	if !isValid {
		return nil, ErrInvalidPadding
	}

	return data[:firstPadIndex], nil
}

// iso10126Remover removes ISO 10126 padding.
// It is the fastest to unpad and always takes constant time.
// It may return an ErrInvalidPadding error and is therefore susceptible to a padding oracle!
func iso10126Remover(data []byte, dataLen int, blockSize int) ([]byte, error) {
	_, _, _, _, padLen, err := checkLengthByte(data, dataLen, blockSize)
	if err != nil {
		return nil, err
	}

	if padLen == 0 {
		return nil, ErrInvalidPadding
	}

	return data[:dataLen-padLen], nil
}

// rfc4303Remover removes RFC 4303 padding (IPSec).
// It may return an ErrInvalidPadding error and is therefore susceptible to a padding oracle!
func rfc4303Remover(data []byte, dataLen int, blockSize int) ([]byte, error) {
	firstIndex, lastIndex, firstPadIndex, padLenByte, _, err := checkLengthByte(data, dataLen, blockSize)
	if err != nil {
		return nil, err
	}

	isValid := true
	padValue := padLenByte
	// Always scan *all* data of the last block to thwart timing attacks.
	for i := lastIndex - 1; i >= firstIndex; i-- {
		padValue--
		isValid = isValid && (data[i] == padValue || i < firstPadIndex)
	}
	if !isValid {
		return nil, ErrInvalidPadding
	}

	return data[:firstPadIndex], nil
}

// iso78164Remover removes ISO 7816-4 padding (Smart cards).
// It may return an ErrInvalidPadding error and is therefore susceptible to a padding oracle!
func iso78164Remover(data []byte, dataLen int, blockSize int) ([]byte, error) {
	lastIndex := dataLen - 1
	firstIndex := dataLen - blockSize
	isValid := true
	wasZero := true
	firstPadIndex := firstIndex

	// Always scan *all* data of the last block to thwart timing attacks.
	for i := lastIndex; i >= firstIndex; i-- {
		actData := data[i]

		if actData != 0 && firstPadIndex == firstIndex {
			if wasZero {
				if actData == 0x80 {
					firstPadIndex = i
				} else {
					isValid = false
				}

				wasZero = false
			}
		} else {
			// This is here to do something over all the bytes looped over.
			// It does not change the value of isValid.
			isValid = isValid && (wasZero || i < firstPadIndex)
		}
	}
	if !isValid {
		return nil, ErrInvalidPadding
	}

	return data[:firstPadIndex], nil
}

// arbitraryTailBytePaddingRemover removes arbitrary tail byte padding.
// It never returns an error and is therefore not susceptible to a padding oracle!
func arbitraryTailBytePaddingRemover(data []byte, dataLen int, blockSize int) ([]byte, error) {
	firstIndex := dataLen - blockSize
	lastIndex := dataLen - 1
	padByte := data[lastIndex]

	firstPadIndex := firstIndex
	// Always scan *all* data of the last block to thwart timing attacks.
	for i := lastIndex - 1; i >= firstIndex; i-- {
		if data[i] != padByte && firstPadIndex == firstIndex {
			firstPadIndex = i + 1
		}
	}

	return data[:firstPadIndex], nil
}

// -------- Helper functions --------

func checkLengthByte(data []byte, dataLen int, blockSize int) (int, int, int, byte, int, error) {
	lastIndex := dataLen - 1
	padLenByte := data[lastIndex]
	padLen := int(padLenByte)

	if (padLen > blockSize) || (padLen > dataLen) || padLenByte == 0 {
		return 0, 0, 0, 0, 0, ErrInvalidPadding
	}

	firstIndex := dataLen - blockSize
	firstPadIndex := dataLen - padLen
	return firstIndex, lastIndex, firstPadIndex, padLenByte, padLen, nil
}
