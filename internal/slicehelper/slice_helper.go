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
// Version: 1.1.0
//
// Change history:
//    2024-02-01: V1.0.0: Created.
//    2024-03-17: V1.1.0: Add FillToCap.
//

// Package slicehelper implements useful helper functions for slices.
package slicehelper

// ******** Public functions ********

// Fill fills a slice with a value in an efficient way up to its length.
func Fill[S ~[]T, T any](s S, v T) {
	sLen := len(s)

	if sLen > 0 {
		doFill(s, v, sLen)
	}
}

// Concat returns a new slice concatenating the passed in slices.
// This is a streamlined version of the slices.Concat function of Go V1.22.
func Concat[S ~[]T, T any](slices ...S) S {
	// 1. Calculate total size.
	size := 0
	for _, s := range slices {
		size += len(s)
	}

	// 2. Make new slice with the total size as the capacity and 0 length.
	result := make(S, 0, size)

	// 3. Append all source slices.
	for _, s := range slices {
		result = append(result, s...)
	}

	return result
}

// ******** Private functions ********

// doFill fills a slice in an efficient way.
func doFill[S ~[]T, T any](s S, v T, l int) {
	// Put the value into the first slice element
	s[0] = v

	// Incrementally duplicate the value into the rest of the slice
	for j := 1; j < l; j <<= 1 {
		copy(s[j:], s[:j])
	}
}
