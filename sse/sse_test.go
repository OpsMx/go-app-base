// Copyright 2023 OpsMx, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sse

import (
	"reflect"
	"strings"
	"testing"
)

func TestSSE_SSERead(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Event
	}{
		{
			"Nothing but colons",
			":\n",
			Event{},
		},
		{
			"data with data",
			"data: foo\n\n",
			Event{"data": "foo"},
		},
		{
			"data with data and colons",
			"data: foo\n:\n:\n\n",
			Event{"data": "foo"},
		},
		{
			"multi-line data",
			"data: foo\ndata: bar\n\n",
			Event{"data": "foo\nbar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sse := NewSSE(strings.NewReader(tt.input))
			if got := sse.SSERead(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SSE.SSERead() = %v, want %v", got, tt.want)
			}
		})
	}
}
