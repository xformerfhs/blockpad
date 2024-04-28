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
This library fills that gap and provides an easy-to-use interface for padding and unpadding of data.
Furthermore, padding and unpadding is nearly constant time to thwart [timing attacks]((https://en.wikipedia.org/wiki/Timing_attack)).

## Usage

First one has to create a padder.
This is a data structure that is able to pad or unpad data.
It is created by calling the `NewBlockPad` function:

```
   padder, err := blockpad.NewBlockPad(padAlgorithm, blockSize)
```

`blockSize` is the size of the underlying block cipher's block size.
`padAlgorithm` specifies the pad algorithm to use.
It has one of the following values:

| padAlgorithm        | Meaning                                                                                                                 |
|---------------------|-------------------------------------------------------------------------------------------------------------------------|
| `Zero`              | [Zero padding](https://en.wikipedia.org/wiki/Padding_(cryptography)#Zero_padding) (ISO 10118-1 and ISO 9797-1 method 1) |
| `PKCS7`             | [PKCS#7 padding](https://en.wikipedia.org/wiki/Padding_(cryptography)#PKCS#5_and_PKCS#7) (RFC 5652)                     |
| `X923`              | [ANSI X.923](https://en.wikipedia.org/wiki/Padding_(cryptography)#ANSI_X9.23) padding                                   |
| `ISO10126`          | [ISO 10126](https://en.wikipedia.org/wiki/Padding_(cryptography)#ISO_10126) padding. This is the fastest to unpad.      |
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

This padder has three public function:

| Function                                | Purpose                                                                                                                                                                                                                                                        |
|-----------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `Pad([]byte) []byte`                    | Given a byte slice of data, it returns a new byte slice that contains the data with the padding. The new byte slice has a length that is a multiple of the block size.                                                                                         |
| `PadLastBlock([]byte) ([]byte, []byte)` | Given a byte slice of data, it returns a byte slice of the data up to the last block and a new slice containing the last block with padding. The data slice has a length that is a multiple of the block size. The length of the last block is the block size. |
| `Unpad([]byte) ([]byte, error)`         | Given a byte slice of padded data, it returns a byte slice into the original data with the padding removed. If there is something wrong with the padding, the returned byte slice is `nil` and an error is returned.                                           |

### Rational

One may ask why the padding and unpadding has not been implemented with a more traditional call interface like e.g. `Pad(padAlgorithm, blockSize, data)` and `Unpad(padAlgorithm, blockSize, data)`.
There are two reasons for this:

1. Performance: The block size and the pad algorithm are checked only once, when the padder is created. With the traditional interface they would have to be checked on every call, which slows down processing by about 30%.
2. Simplicity: With the creation of a padder the call interface is not cluttered with parameters.

## Contact

Frank Schwab ([Mail](mailto:github.sfdhi@slmails.com "Mail"))

## License

This source code is published under the [Apache License V2](https://www.apache.org/licenses/LICENSE-2.0.txt).
