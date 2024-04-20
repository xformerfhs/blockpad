# blockpad

A [Go](https://go.dev) package for block cipher paddings.

## Introduction

One type of encryption is the so-called [block cipher](https://en.wikipedia.org/wiki/Block_cipher).
This means that data is processed in blocks of bytes and not byte by byte.
Block ciphers are substitution ciphers.
One block of data is substituted by another block.
The advantage of blocks is that they allow a large mapping space.

However, this type of processing leads to some practical challenges.
Block ciphers need to be used in a specific [mode](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation).
Some of these modes require that the data are processed in blocks.
These "classic" modes are known by their abbreviations [ECB](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Electronic_codebook_(ECB)), [CBC](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Cipher_block_chaining_(CBC)) or [PCBC](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Propagating_cipher_block_chaining_(PCBC)).

When the last bytes of the clear data are processed a challenge appears:
Most of the time the last clear data will not fill a block completely.
So there has to be a way to know where the real data in the last block ends.
This process is called [padding](https://en.wikipedia.org/wiki/Padding_(cryptography)).

Nowadays, usually [AEAD](https://en.wikipedia.org/wiki/Authenticated_encryption) modes are used, which have built-in integrity protection.
The corresponding modes are referred to as [CCM](https://en.wikipedia.org/wiki/CCM_mode), [EAX](https://en.wikipedia.org/wiki/EAX_mode) , [GCM](https://en.wikipedia.org/wiki/Galois/Counter_Mode), [OCB](https://en.wikipedia.org/wiki/OCB_mode) or [SIV](https://en.wikipedia.org/wiki/Block_cipher_mode_of_operation#Synthetic_initialization_vector_(SIV)).
These have the advantage that they do not need padding and have a built-in integrity protection.
Attacks such as a [padding oracle](https://en.wikipedia.org/wiki/Padding_oracle_attack) are not possible with these modes.
However, the classic modes with additional integrity protection, e.g. through HMACs, are still justified as they are cryptographically more secure than the AEAD modes.
Also, the AEAD modi are sometimes just not feasible, as they can not chain data.

Go has an excellent library of cryptographic primitives.
Strangely, however, it lacks any support for padding.
This library fills that gap and provides an easy to use interface for padding and unpadding of data.

## Usage

First one has to create a padder.
This is a data structure that is able to pad or unpad data.
It is created by calling the `NewBlockPad` function:

```
   padder, err := NewBlockPad(padAlgorithm, blockSize)
```

`blockSize` is the size of the underlying block cipher's block size.
`padAlgorithm` specifies the pad algorithm to use.
It has one of the following values:

| padAlgorithm        | Meaning                                                                                                                 |
|---------------------|-------------------------------------------------------------------------------------------------------------------------|
| `Zero`              | [Zero padding](https://en.wikipedia.org/wiki/Padding_(cryptography)#Zero_padding) (ISO 10118-1 and ISO 9797-1 method 1) |
| `PKCS7`             | [PKCS#7 padding](https://en.wikipedia.org/wiki/Padding_(cryptography)#PKCS#5_and_PKCS#7) (RFC 5652)                     |
| `X923`              | [ANSI X.923](https://en.wikipedia.org/wiki/Padding_(cryptography)#ANSI_X9.23) padding                                   |
| `ISO10126`          | [ISO 10126](https://en.wikipedia.org/wiki/Padding_(cryptography)#ISO_10126) padding                                     |
| `RFC4303`           | [RFC 4303](https://datatracker.ietf.org/doc/html/rfc4303#section-2.4) padding                                           |
| `ISO78164`          | [ISO 7816-4](https://en.wikipedia.org/wiki/Padding_(cryptography)#ISO/IEC_7816-4) padding (ISO 9797-1 method 2)         |
| `ArbitraryTailByte` | [Arbitrary tail byte padding](https://eprint.iacr.org/2003/098.pdf)                                                     |

> [!CAUTION]
> With CBC mode, nearly all the padding methods enable a very dangerous attack, the so-called padding oracle.
> They must only be used with integrity protection, e.g. by an [HMAC](https://en.wikipedia.org/wiki/HMAC).
> Only arbitrary tail byte padding is not susceptible to this attack.
> An integrity protection is **always** advisable.

> [!CAUTION]
> When using Zero padding the clear data **must not** end with a 0 byte.
> Zero padding panics if the clear data ends with a 0 byte.

This padder has two public function:

| Function                      | Purpose                                                                                                                                                                                                              |
|-------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `Pad(data) []byte`            | Given a byte slice of data, it returns a new byte slice that contains the data with the padding. The new byte slice has a length that is a multiple of the block size.                                               |
| `Unpad(data) ([]byte, error)` | Given a byte slice of padded data, it returns a byte slice into the original data with the padding removed. If there is something wrong with the padding, the returned byte slice is `nil` and an error is returned. |

### Rational

One may ask why the padding and unpadding has not been implemented with a more traditional call interface like e.g. `Pad(padAlgorithm, blockSize, data)` and `Unpad(padAlgorithm, blockSize, data)`.
There are two reasons for this:

1. Performance: The block size and the pad algorithm are checked only once, when the padder is created. With the traditional interface they would have to be checked on every call, which slows down processing by about 30%.
2. Simplicity: With the creation of a padder the call interface is not cluttered with parameters.

## Example

In this example a very simple main program calls one encryption and one decryption function.
It shows how a block cipher and a padder would typically be used for encryption and decryption.

```go
package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"github.com/xformerfhs/blockpad"
	"log"
)

func main() {
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
		log.Fatalf(`Could not create AES cipher: %v`, err)
	}

	// 2. Create arbitrary tail byte padder.
	var padder *blockpad.BlockPad
	padder, err = blockpad.NewBlockPadding(blockpad.ArbitraryTailByte, aes.BlockSize)
	if err != nil {
		log.Fatalf(`Could not create padder: %v`, err)
	}

	// 3. Encrypt the data with a unique iv for every encryption.
	var encryptedData []byte
	encryptedData, err = doEncryption(aesCipher, iv, padder, data)
	if err != nil {
		log.Fatalf(`Encryption failed: %v`, err)
	}

	// 4. Decrypt the encrypted data.
	var decryptedData []byte
	decryptedData, err = doDecryption(aesCipher, iv, padder, encryptedData)

	// 5. Check result.
	if bytes.Compare(data, decryptedData) == 0 {
		log.Print(`Success!`)
	} else {
		log.Fatalf(`Decrypted data '%02x' does not match clear data '%02x'`, decryptedData, data)
	}
}

// doEncryption encrypts a slice of data.
func doEncryption(blockCipher cipher.Block, iv []byte, padder *blockpad.BlockPad, clearData []byte) ([]byte, error) {
	// 1. Create block mode from cipher.
	encrypter := cipher.NewCBCEncrypter(blockCipher, iv)

	// 2. Pad clear data.
	paddedData := padder.Pad(clearData)

	// 3. Encrypt padded data.
	// After this, paddedData contains the encrypted padded data.
	encrypter.CryptBlocks(paddedData, paddedData)

	return paddedData, nil
}

// doDecryption decrypts a slice of data.
func doDecryption(blockCipher cipher.Block, iv []byte, padder *blockpad.BlockPad, encryptedData []byte) ([]byte, error) {
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
```

## Contact

Frank Schwab ([Mail](mailto:github.sfdhi@slmails.com "Mail"))

## License

This source code is published under the [Apache License V2](https://www.apache.org/licenses/LICENSE-2.0.txt).
