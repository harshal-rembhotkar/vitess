/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package integration

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"vitess.io/vitess/go/hack"
)

func TestTypeSizes(t *testing.T) {
	var PtrSize = hack.RuntimeAllocSize(8)
	var SliceHeaderSize = hack.RuntimeAllocSize(3 * PtrSize)
	var FatPointerSize = hack.RuntimeAllocSize(2 * PtrSize)

	cases := []struct {
		obj  cachedObject
		size int64
	}{
		{&A{}, 16},
		{&C{}, 16},
		{&C{field1: &Bimpl{}}, 24},
		{&D{}, 8},
		{&D{field1: &Bimpl{}}, 16},
		{&Padded{}, 24},

		{&Slice1{}, 24},
		{&Slice1{field1: []A{}}, SliceHeaderSize},
		{&Slice1{field1: []A{{}}}, SliceHeaderSize + 16},
		{&Slice1{field1: []A{{}, {}, {}, {}}}, SliceHeaderSize + 16*4},

		{&Slice2{}, SliceHeaderSize},
		{&Slice2{field1: []B{}}, SliceHeaderSize},
		{&Slice2{field1: []B{&Bimpl{}}}, SliceHeaderSize + FatPointerSize*1 + 8*1},
		{&Slice2{field1: []B{&Bimpl{}, &Bimpl{}, &Bimpl{}, &Bimpl{}}}, SliceHeaderSize + FatPointerSize*4 + 8*4},

		{&Slice3{}, SliceHeaderSize},
		{&Slice3{field1: []*Bimpl{}}, SliceHeaderSize},
		{&Slice3{field1: []*Bimpl{nil}}, SliceHeaderSize + PtrSize*1 + 0},
		{&Slice3{field1: []*Bimpl{nil, nil, nil, nil}}, SliceHeaderSize + PtrSize*4 + 0},
		{&Slice3{field1: []*Bimpl{{}}}, SliceHeaderSize + PtrSize*1 + 8*1},
		{&Slice3{field1: []*Bimpl{{}, {}, {}, {}}}, SliceHeaderSize + PtrSize*4 + 8*4},

		{&Map1{field1: nil}, PtrSize},
		{&Map1{field1: map[uint8]uint8{}}, 56},
		{&Map1{field1: map[uint8]uint8{0: 0}}, 80},

		{&Map2{field1: nil}, PtrSize},
		{&Map2{field1: map[uint64]A{}}, 56},
		{&Map2{field1: map[uint64]A{0: {}}}, 264},

		{&Map3{field1: nil}, PtrSize},
		{&Map3{field1: map[uint64]B{}}, 56},
		{&Map3{field1: map[uint64]B{0: &Bimpl{}}}, 272},
		{&Map3{field1: map[uint64]B{0: nil}}, 264},

		{&String1{}, hack.RuntimeAllocSize(PtrSize*2 + 8)},
		{&String1{field1: "1234"}, hack.RuntimeAllocSize(PtrSize*2+8) + hack.RuntimeAllocSize(4)},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("sizeof(%T)", tt.obj), func(t *testing.T) {
			size := tt.obj.CachedSize(true)
			assert.Equalf(t, tt.size, size, "expected %T to be %d bytes, got %d", tt.obj, tt.size, size)
		})
	}
}
