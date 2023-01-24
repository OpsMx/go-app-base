// Copyright 2022 OpsMx, Inc
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

package util

import "testing"

func TestContains(t *testing.T) {
	type args struct {
		elems []string
		v     string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"Empty list, blank target",
			args{
				elems: []string{},
				v:     "",
			},
			false,
		}, {
			"Empty list, non-blank target",
			args{
				elems: []string{},
				v:     "foo",
			},
			false,
		}, {
			"blank target",
			args{
				elems: []string{"foo", "bar"},
				v:     "",
			},
			false,
		}, {
			"not present",
			args{
				elems: []string{"foo", "bar"},
				v:     "baz",
			},
			false,
		}, {
			"present",
			args{
				elems: []string{"foo", "bar"},
				v:     "foo",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.args.elems, tt.args.v); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
