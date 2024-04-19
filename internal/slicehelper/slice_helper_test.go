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
//    2024-03-17: V1.0.0: Created.
//

package slicehelper

import (
	"testing"
)

func TestFill(t *testing.T) {
	var s []byte
	Fill(s, 0xaa)
	if len(s) != 0 {
		t.Fatal(`Error filling empty slice.`)
	}

	u := make([]string, 7)
	Fill(u, `Empty`)
	for _, e := range u {
		if e != `Empty` {
			t.Fatal(`Error filling string slice.`)
		}
	}

	v := make([][]int, 7)
	w := make([]int, 3)
	Fill(w, -3)
	Fill(v, w)
	for _, e := range v {
		for _, f := range e {
			if f != -3 {
				t.Fatal(`Error filling slice of int slice.`)
			}
		}
	}
}

func TestConcat(t *testing.T) {
	a := make([]uint64, 7)
	Fill(a, 1)
	b := make([]uint64, 11)
	Fill(b, 99)
	c := Concat(a, b)
	for i, e := range c {
		if i < 7 {
			if e != 1 {
				t.Fatal(`Error in 1. part of concatenated slice.`)
			}
		} else {
			if e != 99 {
				t.Fatal(`Error in 2. part of concatenated slice.`)
			}
		}
	}
}
