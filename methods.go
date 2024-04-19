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

package blockpad

import (
	"padding/internal/slicehelper"
)

// ******** This file contains the publicly callable functions, i.e. the interface ********

// ******** Public creation function ********

// NewBlockPadding creates a block padding.
func NewBlockPadding(padAlgorithm PadAlgorithm, blockSize int) (*BlockPad, error) {
	err := checkPadAlgorithmAndBlockSize(padAlgorithm, blockSize)
	if err != nil {
		return nil, err
	}

	return &BlockPad{
		worker:    padImplementation[padAlgorithm],
		blockSize: blockSize,
	}, nil
}

// ******** Public functions ********

// Pad pads a byte slice.

func (pb *BlockPad) Pad(data []byte) []byte {
	dataLen, lastByte := getLenAndLastByte(data)

	pad := pb.worker.filler(lastByte, dataLen, pb.blockSize)

	return slicehelper.Concat(data, pad)
}

// Unpad removes the padding from a byte slice.
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
