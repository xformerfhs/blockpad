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

// Package blockpad implements functions for block cipher paddings.
//
// It can be used for block ciphers ([crypto/cipher/Block] or [crypto/cipher/BlockMode])
// in e.g. ECB, CBC or PCBC mode, as they require that the size of the data to be encrypted
// is a multiple of the block size.
//
// # Encryption example
//
//		func doEncryption(key []byte, iv []byte, clearData []byte) ([]byte, error) {
//		  aesCipher, err := aes.NewCipher(key)
//		  if err != nil {
//		    return nil, err
//		  }
//
//		  encrypter := cipher.NewCBCEncrypter(aesCipher, iv)
//
//	     var padder blockPad.BlockPad
//		  padder, err = NewBlockPadding(blockpad.ArbitraryTailByte, aes.BlockSize)
//		  if err != nil {
//		    return nil, err
//		  }
//
//	     var paddedData []byte
//		  paddedData, err = padder.Pad(clearData)
//		  if err != nil {
//		    return nil, err
//		  }
//
//		  // After this, paddedData contains the encrypted padded data.
//		  encrypter.CryptBlocks(paddedData, paddedData)
//
//		  return paddedData, nil
//		}
//
// # Decryption example
//
//		func doDecryption(key []byte, iv []byte, encryptedData []byte) ([]byte, error) {
//		  aesCipher, err := aes.NewCipher(key)
//		  if err != nil {
//		    return nil, err
//		  }
//
//		  decrypter := cipher.NewCBCDecrypter(aesCipher, iv)
//
//		  decryptedData := make([]byte, len(encryptedData))
//		  encrypter.CryptBlocks(decryptedData, paddedData)
//
//	     var padder blockPad.BlockPad
//		  padder, err = NewBlockPadding(blockpad.ArbitraryTailByte, aes.BlockSize)
//		  if err != nil {
//		    return nil, err
//		  }
//
//	     var unpaddedData []byte
//		  unpaddedData, err = padder.Unpad(decryptedData)
//		  if err != nil {
//		    return nil, err
//		  }
//
//		  return unpaddedData, nil
//		}
//
// Of course, if multiple data have to be encrypted or decrypted it is advisable to create
// the cipher and the padder only once and pass them as parameters to the encryption and
// decryption functions.
package blockpad

import "errors"

// ******** This file contains the public types, constants and errors ********

// ******** Public types ********

// BlockPad represents an implementation of a block cipher padding.
// It provides the capability to pad data before it is encrypted with a block cipher
// and to unpad padded data after has been decrypted with a block cipher.
//
// A BlockPad is safe for concurrent use by multiple goroutines, as it is used read-only.
type BlockPad struct {
	worker    implementationInfo
	blockSize int
}

// PadAlgorithm is the type that holds pad algorithms.
type PadAlgorithm byte

// ******** Public constants ********

// These are the valid pad algorithms.
const (
	// Zero implements zero padding (ISO 10118-1 and ISO 9797-1 method 1), i.e. zero bytes are appended.
	// Data to be padded *must not* end with a 0 byte! If it does, the Pad function will panic in this mode.
	// This padding should only be used with integrity protection as it is susceptible to a padding oracle attack.
	Zero PadAlgorithm = iota

	// PKCS7 implements PKCS#7 padding (RFC 5652).
	// This padding should only be used with integrity protection as it is susceptible to a padding oracle attack.
	PKCS7

	// X923 implements ANSI X.923 padding.
	// This padding should only be used with integrity protection as it is susceptible to a padding oracle attack.
	X923

	// ISO10126 implements ISO 10126 padding.
	// This padding should only be used with integrity protection as it is susceptible to a padding oracle attack.
	ISO10126

	// RFC4303 implements RFC 4303 padding (IPSec).
	// This padding should only be used with integrity protection as it is susceptible to a padding oracle attack.
	RFC4303

	// ISO78164 implements ISO 7816-4 padding (ISO 9797-1 method 2, smart cards).
	// This padding should only be used with integrity protection as it is susceptible to a padding oracle attack.
	ISO78164

	// ArbitraryTailByte implements arbitrary tail byte padding.
	// This padding is the only one that is *not* susceptible to a padding oracle attack.
	ArbitraryTailByte

	// maxAlgorithm is a helper constant and always contains the maximum defined padding type constant.
	// It must always be the last constant in this const block!
	maxAlgorithm = iota - 1
)

// These are the public errors.

var (
	// ErrInvalidBlockSize means that the provided block size is invalid.
	ErrInvalidBlockSize = errors.New(`invalid block size`)

	// ErrInvalidPadAlgorithm means that the provided block size is invalid.
	ErrInvalidPadAlgorithm = errors.New(`invalid pad algorithm`)

	// ErrInvalidPaddedDataLen means that the provided "padded" data is obviously not padded.
	ErrInvalidPaddedDataLen = errors.New(`padded data length is not a multiple of the block size`)

	// ErrInvalidPadding means that something is wrong with the padding.
	// It is deliberately not stated what exactly is wrong so that
	// an attacker does not obtain too much information.
	ErrInvalidPadding = errors.New(`invalid padding`)
)
