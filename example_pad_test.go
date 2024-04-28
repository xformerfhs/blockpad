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

package blockpad_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/xformerfhs/blockpad"
	"os"
)

func ExampleBlockPad_Pad() {
	data := []byte(`Beware the ides of march`)

	// ATTENTION: Do not hard-code an encryption key! NEVER!
	key := []byte{
		0x7c, 0xa8, 0x69, 0xbc, 0x54, 0xf5, 0x87, 0x99,
		0xf3, 0x89, 0x09, 0xab, 0x33, 0xfb, 0xdb, 0x5c,
		0x84, 0x09, 0x4c, 0x05, 0x23, 0xc1, 0xb1, 0x07,
		0xa5, 0xea, 0x5d, 0xf7, 0xf5, 0x42, 0x77, 0x42,
	}

	// ATTENTION: Never use a constant initialization vector! NEVER!
	iv := []byte{
		0x11, 0x42, 0xcd, 0x7d, 0xf9, 0x98, 0x46, 0x4d,
		0xd8, 0x58, 0xdd, 0x4e, 0xc8, 0x3b, 0xfd, 0xe9,
	}

	// 1. Create block cipher.
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf(`Could not create AES cipher: %v`, err)
		os.Exit(1)
	}

	// 2. Create arbitrary tail byte padder.
	var padder *blockpad.BlockPad
	padder, err = blockpad.NewBlockPadding(blockpad.ArbitraryTailByte, aes.BlockSize)
	if err != nil {
		fmt.Printf(`Could not create padder: %v`, err)
		os.Exit(1)
	}

	// 3. Encrypt the data with a unique iv for every encryption.
	var encryptedData []byte
	encryptedData, err = doEncryptionWithPad(aesCipher, iv, padder, data)
	if err != nil {
		fmt.Printf(`Encryption failed: %v`, err)
		os.Exit(1)
	}

	// 4. Decrypt the encrypted data.
	var decryptedData []byte
	decryptedData, err = doDecryptionWithPad(aesCipher, iv, padder, encryptedData)
	if err != nil {
		fmt.Printf(`Decryption failed: %v`, err)
		os.Exit(1)
	}

	// 5. Check result.
	if bytes.Equal(data, decryptedData) {
		fmt.Print(`Success!`)
	} else {
		fmt.Printf(`Decrypted data '%02x' does not match clear data '%02x'`, decryptedData, data)
		os.Exit(1)
	}

	// Output: Success!
}

// ******** Private functions ********

// doEncryptionWithPad encrypts a slice of data.
func doEncryptionWithPad(
	blockCipher cipher.Block,
	iv []byte,
	padder *blockpad.BlockPad,
	clearData []byte,
) ([]byte, error) {
	// 1. Create block mode from cipher.
	encrypter := cipher.NewCBCEncrypter(blockCipher, iv)

	// 2. Pad clear data.
	paddedData := padder.Pad(clearData)

	// 3. Encrypt padded data.
	// After this, paddedData contains the encrypted padded data.
	encrypter.CryptBlocks(paddedData, paddedData)

	return paddedData, nil
}

// doDecryptionWithPad decrypts a slice of data.
func doDecryptionWithPad(
	blockCipher cipher.Block,
	iv []byte,
	padder *blockpad.BlockPad,
	encryptedData []byte,
) ([]byte, error) {
	// 1. Create block mode from cipher.
	decrypter := cipher.NewCBCDecrypter(blockCipher, iv)

	// 2. Decrypt padded data.
	decryptedData := make([]byte, len(encryptedData))
	decrypter.CryptBlocks(decryptedData, encryptedData)

	// 3. Unpad padded data.
	unpaddedData, err := padder.Unpad(decryptedData)
	if err != nil {
		return nil, err
	}

	return unpaddedData, nil
}
