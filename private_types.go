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

// ******** This file contains the private types ********

// ******** Private types ********

// fillerFunc is the type of a filler function.
type fillerFunc func(byte, int, int) []byte

// removerFunc is the type of a remover function.
type removerFunc func([]byte, int, int) ([]byte, error)

// implementationInfo holds the data necessary for doing padding and unpadding.
type implementationInfo struct {
	name    string
	filler  fillerFunc
	remover removerFunc
}
