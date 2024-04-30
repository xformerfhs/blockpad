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
	"github.com/xformerfhs/blockpad/internal/slicehelper"
)

// ******** This file contains the publicly callable functions, i.e. the interface ********

// ******** Public creation function ********

// NewBlockPadding creates a block padding.
func NewBlockPadding(padAlgorithm PadAlgorithm, blockSize int) (*BlockPad, error) {
	err := checkPadAlgorithmAndBlockSize(padAlgorithm, blockSize)
	if err != nil {
		return nil, err
	}

	worker := padImplementation[padAlgorithm]
	return &BlockPad{
		worker:    worker,
		blockSize: blockSize,
		zeroBlock: make([]byte, blockSize),
	}, nil
}

// ******** Public functions ********

// Pad pads a byte slice.
// It returns a new slice that is a copy of the data with added padding.
// If the data is large this is inefficient.
// PadLastBlock contains a more efficient implementation that avoids
// copying all the data.
func (pb *BlockPad) Pad(data []byte) []byte {
	fullBlockData, lastBlock := pb.PadLastBlock(data)

	return slicehelper.Concat(fullBlockData, lastBlock)
}

// PadLastBlock pads a byte slice.
// It returns a byte slice of the data up to the last block
// and a new slice containing the last block with padding.
// Only the last data that does not fit into a full block is copied.
// This is much more efficient than Pad.
func (pb *BlockPad) PadLastBlock(data []byte) ([]byte, []byte) {
	// 1. Get all kind of lengths.
	dataLen := len(data)
	blockSize := pb.blockSize

	fullBlockDataLen, lastBlockDataLen, padLen := padLengths(dataLen, blockSize)
	lastBlock := make([]byte, blockSize)
	lastData := data[fullBlockDataLen:]

	// There are two copy operations. The first one copies padLen bytes and the second one lastBlockDataLen bytes.
	// padLen + lastBlockDataLen = blockSize, so there are always blockSize bytes copied.

	// 2. Do some additional - functionally unnecessary - copying to achieve constant time.
	copy(lastBlock, pb.zeroBlock[:padLen]) // This copies padLen bytes.

	// 3. Build a full block of filler bytes to help achieve constant-time processing.
	pb.worker.filler(lastBlock, blockSize, lastData, lastBlockDataLen, padLen)

	// 4. Finally, copy last data to last block.
	copy(lastBlock, lastData[:lastBlockDataLen]) // This copies lastBlockDataLen bytes. lastBlockDataLen + padLen = blockSize.

	return data[:fullBlockDataLen], lastBlock
}

// Unpad removes the padding from a byte slice.
// It returns a byte slice into the supplied data and does not allocate a new slice.
// If a last block is unpadded it returns a zero-length slice if that last block contains only padding.
func (pb *BlockPad) Unpad(data []byte) ([]byte, error) {
	dataLen := len(data)
	if dataLen%pb.blockSize != 0 {
		return nil, ErrInvalidPaddedDataLen
	}

	return pb.worker.remover(data, dataLen, pb.blockSize)
}

// String yields the name of the padding algorithm.
// It implements the Stringer interface.
func (pb *BlockPad) String() string {
	return pb.worker.name
}

// ******** Private functions ********

// padLengths calculates the 3 lengths needed for padding.
// It returns the length of full data blocks, the length of the last data block
// and the length of the padding needed.
func padLengths(dataLen int, blockSize int) (int, int, int) {
	fullBlockCount := dataLen / blockSize
	fullBlockDataLen := fullBlockCount * blockSize
	lastBlockDataLen := dataLen - fullBlockDataLen
	padLen := blockSize - lastBlockDataLen

	return fullBlockDataLen, lastBlockDataLen, padLen
}
