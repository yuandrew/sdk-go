// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package converter

import (
	"fmt"
	commonpb "go.temporal.io/api/common/v1"
)

type (
	// EncodedValue is used to encapsulate/extract encoded value from workflow/activity.
	EncodedValue interface {
		// HasValue return whether there is value encoded.
		HasValue() bool
		// Get extract the encoded value into strong typed value pointer.
		//
		// Note, values should not be reused for extraction here because merging on
		// top of existing values may result in unexpected behavior similar to
		// json.Unmarshal.
		Get(valuePtr interface{}) error
	}

	// EncodedValues is used to encapsulate/extract encoded one or more values from workflow/activity.
	EncodedValues interface {
		// HasValues return whether there are values encoded.
		HasValues() bool
		// Get extract the encoded values into strong typed value pointers.
		//
		// Note, values should not be reused for extraction here because merging on
		// top of existing values may result in unexpected behavior similar to
		// json.Unmarshal.
		Get(valuePtr ...interface{}) error
	}

	// RawValue is a representation of an unconverted, raw payload.
	//
	// This type can be used as a parameter or return type in workflows and activities to pass through
	// a raw payload. Encoding/decoding of the payload is still done by the system. A RawValue enabled
	// payload converter is required for this.
	RawValue struct {
		payload *commonpb.Payload
	}
)

// NewRawValue creates a new RawValue instance.
func NewRawValue(payload *commonpb.Payload) RawValue {
	return RawValue{payload: payload}
}

func (v RawValue) Payload() *commonpb.Payload {
	return v.payload
}

func (v RawValue) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("RawValue is not JSON serializable")
}

func (v *RawValue) UnmarshalJSON(b []byte) error {
	return fmt.Errorf("RawValue is not JSON serializable")
}
