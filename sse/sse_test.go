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
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestSSE_Read(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Event
		wantEOF bool
	}{
		{
			"Nothing but colons",
			":\n",
			Event{},
			false,
		},
		{
			"EOF",
			"",
			Event{},
			true,
		},
		{
			"data with data",
			"data: foo\n\n",
			Event{"data": "foo"},
			false,
		},
		{
			"data with data and colons",
			"data: foo\n:\n:\n\n",
			Event{"data": "foo"},
			false,
		},
		{
			"multi-line data",
			"data: foo\ndata: bar\n\n",
			Event{"data": "foo\nbar"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sse := NewSSE(strings.NewReader(tt.input))
			got, gotEOF := sse.Read()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SSE.Read() = %v, want %v", got, tt.want)
			}
			if gotEOF != tt.wantEOF {
				t.Errorf("SSE.Read() EOF == %v, want %v", gotEOF, tt.wantEOF)
			}
		})
	}
}

func TestSSE_KeepAlive(t *testing.T) {
	tests := []struct {
		name    string
		wantW   string
		wantErr bool
	}{
		{
			"works",
			":\n",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sse := NewSSE(strings.NewReader(""))
			w := &bytes.Buffer{}
			if err := sse.KeepAlive(w); (err != nil) != tt.wantErr {
				t.Errorf("SSE.KeepAlive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("SSE.KeepAlive() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestSSE_Write(t *testing.T) {
	tests := []struct {
		name    string
		event   Event
		wantW   string
		wantErr bool
	}{
		{
			"empty event",
			Event{},
			"",
			false,
		},
		{
			"data only",
			Event{
				"data": `{"foo":"bar"}`,
			},
			"data: {\"foo\":\"bar\"}\n\n",
			false,
		},
		{
			"multi-line data",
			Event{
				"data": "foo\nbar",
			},
			"data: foo\ndata: bar\n\n",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sse := NewSSE(strings.NewReader(""))
			w := &bytes.Buffer{}
			if err := sse.Write(w, tt.event); (err != nil) != tt.wantErr {
				t.Errorf("SSE.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("SSE.Write() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
