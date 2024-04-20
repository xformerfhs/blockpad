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
// package main
//
// # Example
// import (
//    "bytes"
//    "crypto/aes"
//    "crypto/cipher"
//    "github.com/xformerfhs/blockpad"
//    "log"
// )
//
// func main() {
//    data := []byte(`Beware the ides of march`)
//
//    // ATTENTION: Do not hard-code an encryption key! NEVER!
//    key := []byte{
//       0x7c, 0xa8, 0x69, 0xbc, 0x54, 0xf5, 0x87, 0x99,
//       0xf3, 0x89, 0x09, 0xab, 0x33, 0xfb, 0xdb, 0x5c,
//       0x84, 0x09, 0x4c, 0x05, 0x23, 0xc1, 0xb1, 0x07,
//       0xa5, 0xea, 0x5d, 0xf7, 0xf5, 0x42, 0x77, 0x42,
//    }
//
//    // ATTENTION: Never use a constant initialization vector! NEVER!
//    iv := []byte{
//       0x11, 0x42, 0xcd, 0x7d, 0xf9, 0x98, 0x46, 0x4d,
//       0xd8, 0x58, 0xdd, 0x4e, 0xc8, 0x3b, 0xfd, 0xe9,
//    }
//
//    // 1. Create block cipher.
//    aesCipher, err := aes.NewCipher(key)
//    if err != nil {
//       log.Fatalf(`Could not create AES cipher: %v`, err)
//    }
//
//    // 2. Create arbitrary tail byte padder.
//    var padder *blockpad.BlockPad
//    padder, err = blockpad.NewBlockPadding(blockpad.ArbitraryTailByte, aes.BlockSize)
//    if err != nil {
//       log.Fatalf(`Could not create padder: %v`, err)
//    }
//
//    // 3. Encrypt the data with a unique iv for every encryption.
//    var encryptedData []byte
//    encryptedData, err = doEncryption(aesCipher, iv, padder, data)
//    if err != nil {
//       log.Fatalf(`Encryption failed: %v`, err)
//    }
//
//    // 4. Decrypt the encrypted data.
//    var decryptedData []byte
//    decryptedData, err = doDecryption(aesCipher, iv, padder, encryptedData)
//    if err != nil {
//       log.Fatalf(`Decryption failed: %v`, err)
//    }
//
//    // 5. Check result.
//    if bytes.Compare(data, decryptedData) == 0 {
//       log.Print(`Success!`)
//    } else {
//       log.Fatalf(`Decrypted data '%02x' does not match clear data '%02x'`, decryptedData, data)
//    }
// }
//
// // doEncryption encrypts a slice of data.
// func doEncryption(blockCipher cipher.Block, iv []byte, padder *blockpad.BlockPad, clearData []byte) ([]byte, error) {
//    // 1. Create block mode from cipher.
//    encrypter := cipher.NewCBCEncrypter(blockCipher, iv)
//
//    // 2. Pad clear data.
//    paddedData := padder.Pad(clearData)
//
//    // 3. Encrypt padded data.
//    // After this, paddedData contains the encrypted padded data.
//    encrypter.CryptBlocks(paddedData, paddedData)
//
//    return paddedData, nil
// }
//
// // doDecryption decrypts a slice of data.
// func doDecryption(blockCipher cipher.Block, iv []byte, padder *blockpad.BlockPad, encryptedData []byte) ([]byte, error) {
//    // 1. Create block mode from cipher.
//    decrypter := cipher.NewCBCDecrypter(blockCipher, iv)
//
//    // 2. Decrypt padded data.
//    decryptedData := make([]byte, len(encryptedData))
//    decrypter.CryptBlocks(decryptedData, encryptedData)
//
//    // 3. Unpad padded data.
//    unpaddedData, err := padder.Unpad(decryptedData)
//    if err != nil {
//       return nil, err
//    }
//
//    return unpaddedData, nil
// }

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
