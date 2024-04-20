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
	"bytes"
	"crypto/rand"
	"errors"
	"github.com/xformerfhs/blockpad/internal/slicehelper"
	mrand "math/rand"
	"strings"
	"sync"
	"testing"
)

// ******** Private constants ********

// loopCount is the number of times a functional test is to be performed.
const loopCount = 100

// parallelCount is the number of parallel functional tests to be performed.
const parallelCount = 20

// minDataLen is the minimum data length for tests with random data lengths.
const minDataLen = 4

// supDataLen is the supinum of the data length for tests with random data lengths.
// The maximum data length will be supDataLen + minDataLen - 1.
const supDataLen = 411

// testBlockSize is the block size used in the tests.
const testBlockSize = 16

// fixedDataLen is the data length used for tests with fixed lengths.
const fixedDataLen = testBlockSize + 1

// ******** Functional tests ********

func TestPadAllSequential(t *testing.T) {
	for padType := Zero; padType <= maxAlgorithm; padType++ {
		padder, err := NewBlockPadding(padType, testBlockSize)
		if err != nil {
			t.Fatalf(`Error creating BlockPad with pad type %d: %v`, padType, err)
		}

		for i := 0; i < loopCount; i++ {
			dataLen, data := makeZeroSafeRandomLenTestSlice(padType)

			doPadAndUnpad(t, padder, data, dataLen)
		}
	}
}

func TestPadAllParallel(t *testing.T) {
	var wg sync.WaitGroup

	for padType := Zero; padType <= maxAlgorithm; padType++ {
		padder, err := NewBlockPadding(padType, testBlockSize)
		if err != nil {
			t.Fatalf(`Error creating BlockPad with pad type %d: %v`, padType, err)
		}

		for i := 0; i < parallelCount; i++ {
			dataLen, data := makeZeroSafeRandomLenTestSlice(padType)

			wg.Add(1)
			go doPadAndUnpadParallel(t, &wg, padder, data, dataLen)
		}

		wg.Wait()
	}
}

func TestPadZeroLength(t *testing.T) {
	for padType := Zero; padType <= maxAlgorithm; padType++ {
		padder, err := NewBlockPadding(padType, testBlockSize)
		if err != nil {
			t.Fatalf(`Error creating BlockPad with pad type %d: %v`, padType, err)
		}

		var data []byte
		var paddedData []byte
		paddedData = padder.Pad(data)

		var unpaddedData []byte
		unpaddedData, err = padder.Unpad(paddedData)
		if err != nil {
			t.Fatalf(`%s: unpad failed: %v`, padder.String(), err)
		}
		if !bytes.Equal(unpaddedData, data) {
			t.Fatal(padder, `:`, `unpaddedData != data`)
		}
	}
}

func TestNames(t *testing.T) {
	for padType := Zero; padType <= maxAlgorithm; padType++ {
		padder, err := NewBlockPadding(padType, testBlockSize)
		if err != nil {
			t.Fatalf(`Error creating BlockPad with pad type %d: %v`, padType, err)
		}

		name := padder.String()
		if len(name) == 0 {
			t.Fatalf(`Padder %d did not return a name`, padType)
		}
	}
}

// ******** Test invalid paddings ********

func TestPadCrossUnpad(t *testing.T) {
	for padAlgorithm := Zero; padAlgorithm <= maxAlgorithm; padAlgorithm++ {
		padder, err := NewBlockPadding(padAlgorithm, testBlockSize)
		if err != nil {
			t.Fatalf(`Error creating BlockPad with pad type %d: %v`, padAlgorithm, err)
		}

		_, data := makeZeroSafeFixedLenTestSlice(padAlgorithm, fixedDataLen)

		var paddedData []byte
		paddedData = padder.Pad(data)

		for otherPadAlgorithm := Zero; otherPadAlgorithm <= maxAlgorithm; otherPadAlgorithm++ {
			if padAlgorithm != otherPadAlgorithm &&
				!(padAlgorithm == ArbitraryTailByte || otherPadAlgorithm == ArbitraryTailByte) &&
				!((padAlgorithm == PKCS7 || padAlgorithm == X923 || padAlgorithm == RFC4303) && otherPadAlgorithm == ISO10126) &&
				!(padAlgorithm == ISO78164 && otherPadAlgorithm == Zero) {
				var otherPadder *BlockPad
				otherPadder, err = NewBlockPadding(otherPadAlgorithm, testBlockSize)
				if err != nil {
					t.Fatalf(`Error creating other BlockPad with pad type %d: %v`, padAlgorithm, err)
				}

				_, err = otherPadder.Unpad(paddedData)
				if err == nil {
					t.Fatalf(`Cross unpad succeeded between %s and %s`, padder.String(), otherPadder.String())
				}
			}
		}
	}
}

func TestInvalidZeroPadding(t *testing.T) {
	defer func() {
		p := recover()
		if p == nil {
			t.Fatal(`no panic when padding zero data with Zero padding`)
		}
		msg, isString := p.(string)
		if !isString {
			t.Fatal(`panic did not return a string`)
		}
		if !strings.Contains(msg, `must not be 0`) {
			t.Fatalf(`panic with wrong message: %s`, msg)
		}
	}()

	padder, err := NewBlockPadding(Zero, testBlockSize)
	if err != nil {
		t.Fatalf(`Error creating BlockPad with pad type %d: %v`, Zero, err)
	}

	data := make([]byte, 7)
	_ = padder.Pad(data)
}

func TestInvalidPKCS7Padding(t *testing.T) {
	padder, err := NewBlockPadding(PKCS7, testBlockSize)
	if err != nil {
		t.Fatalf(`Error creating BlockPad with pad type %d: %v`, PKCS7, err)
	}

	data := make([]byte, 7)
	slicehelper.Fill(data, 77)

	paddedData := padder.Pad(data)

	paddedData[len(paddedData)-3] = 99

	_, err = padder.Unpad(paddedData)
	if err == nil {
		t.Fatalf(`%s: no error with wrong padding data`, padder.String())
	}
	if !strings.Contains(err.Error(), `invalid padding`) {
		t.Fatalf(`%s: wrong error with wrong padding data: %v`, padder.String(), err)
	}
}

func TestInvalidX923Padding(t *testing.T) {
	padder, err := NewBlockPadding(X923, testBlockSize)
	if err != nil {
		t.Fatalf(`Error creating BlockPad with pad type %d: %v`, X923, err)
	}

	data := make([]byte, 7)
	slicehelper.Fill(data, 77)

	paddedData := padder.Pad(data)

	paddedData[len(paddedData)-3] = 99

	_, err = padder.Unpad(paddedData)
	if err == nil {
		t.Fatalf(`%s: no error with wrong padding data`, padder.String())
	}
	if !errors.Is(err, ErrInvalidPadding) {
		t.Fatalf(`%s: wrong error with wrong padding data: %v`, padder.String(), err)
	}
}

func TestInvalidRFC4303Padding(t *testing.T) {
	padder, err := NewBlockPadding(RFC4303, testBlockSize)
	if err != nil {
		t.Fatalf(`Error creating BlockPad with pad type %d: %v`, RFC4303, err)
	}

	data := make([]byte, 7)
	slicehelper.Fill(data, 77)

	paddedData := padder.Pad(data)

	paddedData[len(paddedData)-3] = 99

	_, err = padder.Unpad(paddedData)
	if err == nil {
		t.Fatalf(`%s: no error with wrong padding data`, padder.String())
	}
	if !errors.Is(err, ErrInvalidPadding) {
		t.Fatalf(`%s: wrong error with wrong padding data: %v`, padder.String(), err)
	}
}

func TestInvalidISO78164Padding(t *testing.T) {
	padder, err := NewBlockPadding(ISO78164, testBlockSize)
	if err != nil {
		t.Fatalf(`Error creating BlockPad with pad type %d: %v`, ISO78164, err)
	}

	data := make([]byte, 7)
	slicehelper.Fill(data, 77)

	paddedData := padder.Pad(data)

	paddedData[len(paddedData)-1] = 99

	_, err = padder.Unpad(paddedData)
	if err == nil {
		t.Fatalf(`%s: no error with invalid last padding byte`, padder.String())
	}
	if !errors.Is(err, ErrInvalidPadding) {
		t.Fatalf(`%s: wrong error with invalid last padding byte: %v`, padder.String(), err)
	}

	paddedData[len(paddedData)-1] = 0
	paddedData[len(data)] = 99
	_, err = padder.Unpad(paddedData)
	if err == nil {
		t.Fatalf(`%s: no error with invalid end marker`, padder.String())
	}
	if !errors.Is(err, ErrInvalidPadding) {
		t.Fatalf(`%s: wrong error invalid end marker: %v`, padder.String(), err)
	}
}

func TestUnpadNoPadding(t *testing.T) {
	data := make([]byte, testBlockSize<<1)
	_, _ = rand.Read(data)
	data[len(data)-1] = 0x5a

	for padType := Zero; padType <= maxAlgorithm; padType++ {
		if padType != ArbitraryTailByte {
			padder, err := NewBlockPadding(padType, testBlockSize)
			if err != nil {
				t.Fatalf(`Error creating BlockPad with pad type %d: %v`, padType, err)
			}
			_, err = padder.Unpad(data)
			if err == nil {
				t.Fatal(padder, `:`, `no error with invalid padded data`)
			}

			if !errors.Is(err, ErrInvalidPadding) {
				t.Fatalf(`%s: wrong error unpadding invalid padded data: %v`, padder.String(), err)
			}
		}
	}
}

func TestUnpadWrongSize(t *testing.T) {
	data := make([]byte, (testBlockSize<<1)-3)

	for padType := Zero; padType <= maxAlgorithm; padType++ {
		padder, err := NewBlockPadding(padType, testBlockSize)
		if err != nil {
			t.Fatalf(`Error creating BlockPad with pad type %d: %v`, padType, err)
		}

		_, err = padder.Unpad(data)
		if err == nil {
			t.Fatalf(`%s: no error with padded data of wrong size`, padder.String())
		}

		if !errors.Is(err, ErrInvalidPaddedDataLen) {
			t.Fatalf(`%s: wrong error with padded data of wrong size: %v`, padder.String(), err)
		}
	}
}

// ******* Invalid parameters ********

func TestTooLargePadType(t *testing.T) {
	_, err := NewBlockPadding(255, testBlockSize)
	if err == nil {
		t.Fatal(`No error creating BlockPad with too large pad type`)
	}
	if !errors.Is(err, ErrInvalidPadAlgorithm) {
		t.Fatalf(`Wrong error creating BlockPad with too large pad algorithm: %v`, err)
	}
}

func TestNegativeBlockSize(t *testing.T) {
	_, err := NewBlockPadding(Zero, -1)
	if err == nil {
		t.Fatal(`No error creating BlockPad with negative block size`)
	}
	if !errors.Is(err, ErrInvalidBlockSize) {
		t.Fatalf(`Wrong error creating BlockPad with negative block size: %v`, err)
	}
}

func TestTooLargeBlockSize(t *testing.T) {
	_, err := NewBlockPadding(Zero, 349127)
	if err == nil {
		t.Fatal(`No error creating BlockPad with too large block size`)
	}
	if !errors.Is(err, ErrInvalidBlockSize) {
		t.Fatalf(`Wrong error creating BlockPad with too large block size: %v`, err)
	}
}

// ******** Private functions ********

// -------- Test slice creation functions ---------

// makeZeroSafeRandomLenTestSlice creates a Zero-safe test slice of random length.
func makeZeroSafeRandomLenTestSlice(padType PadAlgorithm) (int, []byte) {
	dataLen, data := makeRandomLenTestSlice()
	return makeZeroSafeTestSlice(padType, dataLen, data)
}

// makeZeroSafeFixedLenTestSlice creates a Zero-safe test slice of fixed length.
func makeZeroSafeFixedLenTestSlice(padType PadAlgorithm, len int) (int, []byte) {
	dataLen, data := makeFixedLenTestSlice(len)
	return makeZeroSafeTestSlice(padType, dataLen, data)
}

// makeRandomLenTestSlice creates a test slice of random length.
func makeRandomLenTestSlice() (int, []byte) {
	dataLen := mrand.Intn(supDataLen) + minDataLen
	data := makeTestSlice(dataLen)

	return dataLen, data
}

// makeFixedLenTestSlice creates a test slice of fixed length.
func makeFixedLenTestSlice(len int) (int, []byte) {
	data := makeTestSlice(len)

	return len, data
}

// makeFixedLenTestSlice creates a test slice of a given length.
func makeTestSlice(len int) []byte {
	data := make([]byte, len)
	_, _ = rand.Read(data)

	return data
}

// makeZeroSafeTestSlice takes a test slice and makes it Zero-safe, if necessary.
func makeZeroSafeTestSlice(padType PadAlgorithm, dataLen int, data []byte) (int, []byte) {
	// Zero padding will not work if last byte is 0.
	if padType == Zero && data[dataLen-1] == 0 {
		data[dataLen-1] = 0xff
	}

	return dataLen, data
}

// -------- Pad / Unpad runners --------

// doPadAndUnpadParallel runs a pad/unpad test in a Go routine.
func doPadAndUnpadParallel(t *testing.T, wg *sync.WaitGroup, padder *BlockPad, data []byte, dataLen int) {
	defer wg.Done()
	doPadAndUnpad(t, padder, data, dataLen)
}

// doPadAndUnpad runs a pad/unpad test.
func doPadAndUnpad(t *testing.T, padder *BlockPad, data []byte, dataLen int) {
	paddedData := padder.Pad(data)

	unpaddedData, err := padder.Unpad(paddedData)
	if err != nil {
		t.Fatalf(`%s: Unpad failed (dataLen=%d): %v`, padder.String(), dataLen, err)
	}
	if !bytes.Equal(unpaddedData, data) {
		t.Fatalf("%s: unpaddedData != data:\n        data=%02x\n  paddedData=%02x\nunpaddedData=%02x",
			padder.String(),
			data, paddedData, unpaddedData)
	}
}

// ******** Benchmarks ********

func BenchmarkMethodPad(b *testing.B) {
	b.StopTimer()
	_, data := makeFixedLenTestSlice(testBlockSize*5 + 1)
	padder, err := NewBlockPadding(PKCS7, testBlockSize)
	if err != nil {
		b.Fatalf(`padder creation failed: %v`, err)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = padder.Pad(data)
	}
}
