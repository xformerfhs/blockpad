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

func ExampleBlockPad_PadLastBlock() {
	data := []byte(`Cryptography is fun`)

	// ATTENTION: Do not hard-code an encryption key! NEVER!
	key := []byte{
		0x0d, 0xf4, 0x9f, 0x19, 0xe5, 0xf3, 0x91, 0x5c,
		0x55, 0xc0, 0x32, 0xc8, 0x51, 0x98, 0x3c, 0xaf,
		0x05, 0xf7, 0x17, 0xef, 0x7d, 0xc8, 0x5d, 0xcb,
		0xf5, 0xb9, 0xed, 0x74, 0x86, 0xec, 0xed, 0x7b,
	}

	// ATTENTION: Never use a constant initialization vector! NEVER!
	iv := []byte{
		0xea, 0xce, 0xff, 0x10, 0x4f, 0x6a, 0x65, 0xae,
		0x7f, 0x78, 0xe7, 0x43, 0x2c, 0x02, 0x66, 0xf0,
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
	var encryptedFullData []byte
	var encryptedLastBlock []byte
	encryptedFullData, encryptedLastBlock, err = doEncryptionWithPadLastBlock(aesCipher, iv, padder, data)
	if err != nil {
		fmt.Printf(`Encryption failed: %v`, err)
		os.Exit(1)
	}

	// 4. Decrypt the encrypted data.
	var decryptedFullData []byte
	var decryptedLastBlock []byte
	decryptedFullData, decryptedLastBlock, err = doDecryptionWithPadLastBlock(aesCipher, iv, padder, encryptedFullData, encryptedLastBlock)
	if err != nil {
		fmt.Printf(`Decryption failed: %v`, err)
		os.Exit(1)
	}

	// 5. Check result.
	clearFullBlockData := data[:len(encryptedFullData)]
	clearLastBlockData := data[len(encryptedFullData):]
	if bytes.Equal(clearFullBlockData, decryptedFullData) {
		if bytes.Equal(clearLastBlockData, decryptedLastBlock) {
			fmt.Print(`Success!`)
		} else {
			fmt.Printf(`Decrypted last block '%02x' does not match clear data last block '%02x'`, decryptedLastBlock, clearLastBlockData)
			os.Exit(1)
		}
	} else {
		fmt.Printf(`Decrypted full block data '%02x' does not match clear full block data '%02x'`, decryptedFullData, clearFullBlockData)
		os.Exit(1)
	}

	// Output: Success!
}

// ******** Private functions ********

// doEncryptionWithPadLastBlock encrypts a slice of data.
func doEncryptionWithPadLastBlock(
	blockCipher cipher.Block,
	iv []byte,
	padder *blockpad.BlockPad,
	clearData []byte,
) ([]byte, []byte, error) {
	// 1. Create block mode from cipher.
	encrypter := cipher.NewCBCEncrypter(blockCipher, iv)

	// 2. Pad clear data.
	fullBlockData, lastBlock := padder.PadLastBlock(clearData)

	// 3. Encrypt padded data.
	// After this, fullBlockData and lastBlock contain the encrypted padded data.
	encrypter.CryptBlocks(fullBlockData, fullBlockData)
	encrypter.CryptBlocks(lastBlock, lastBlock)

	return fullBlockData, lastBlock, nil
}

// doDecryptionWithPadLastBlock decrypts a slice of data.
func doDecryptionWithPadLastBlock(
	blockCipher cipher.Block,
	iv []byte,
	padder *blockpad.BlockPad,
	encryptedFullBlockData []byte,
	encryptedLastBlock []byte,
) ([]byte, []byte, error) {
	// 1. Create block mode from cipher.
	decrypter := cipher.NewCBCDecrypter(blockCipher, iv)

	decrypter.CryptBlocks(encryptedFullBlockData, encryptedFullBlockData)
	decrypter.CryptBlocks(encryptedLastBlock, encryptedLastBlock)

	// 3. Unpad padded data.
	unpaddedLastBlock, err := padder.Unpad(encryptedLastBlock)
	if err != nil {
		return nil, nil, err
	}

	return encryptedFullBlockData, unpaddedLastBlock, nil
}
