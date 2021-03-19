// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flatbuf

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

/// ----------------------------------------------------------------------
/// A Buffer represents a single contiguous memory segment
type Buffer struct {
	_tab flatbuffers.Struct
}

func (rcv *Buffer) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Buffer) Table() flatbuffers.Table {
	return rcv._tab.Table
}

/// The relative offset into the shared memory page where the bytes for this
/// buffer starts
func (rcv *Buffer) Offset() int64 {
	return rcv._tab.GetInt64(rcv._tab.Pos + flatbuffers.UOffsetT(0))
}
/// The relative offset into the shared memory page where the bytes for this
/// buffer starts
func (rcv *Buffer) MutateOffset(n int64) bool {
	return rcv._tab.MutateInt64(rcv._tab.Pos+flatbuffers.UOffsetT(0), n)
}

/// The absolute length (in bytes) of the memory buffer. The memory is found
/// from offset (inclusive) to offset + length (non-inclusive). When building
/// messages using the encapsulated IPC message, padding bytes may be written
/// after a buffer, but such padding bytes do not need to be accounted for in
/// the size here.
func (rcv *Buffer) Length() int64 {
	return rcv._tab.GetInt64(rcv._tab.Pos + flatbuffers.UOffsetT(8))
}
/// The absolute length (in bytes) of the memory buffer. The memory is found
/// from offset (inclusive) to offset + length (non-inclusive). When building
/// messages using the encapsulated IPC message, padding bytes may be written
/// after a buffer, but such padding bytes do not need to be accounted for in
/// the size here.
func (rcv *Buffer) MutateLength(n int64) bool {
	return rcv._tab.MutateInt64(rcv._tab.Pos+flatbuffers.UOffsetT(8), n)
}

func CreateBuffer(builder *flatbuffers.Builder, offset int64, length int64) flatbuffers.UOffsetT {
	builder.Prep(8, 16)
	builder.PrependInt64(length)
	builder.PrependInt64(offset)
	return builder.Offset()
}
