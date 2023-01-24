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
	"bufio"
	"io"
	"strings"
)

type SSE struct {
	scanner *bufio.Scanner
}

type Event map[string]string

func NewSSE(r io.Reader) *SSE {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	return &SSE{scanner: scanner}
}

func (sse *SSE) SSERead() Event {
	ret := Event{}

	for sse.scanner.Scan() {
		line := sse.scanner.Text()
		if line == "" {
			if len(ret) > 0 {
				return ret
			}
		}
		if line == ":" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		current := ret[parts[0]]
		if current != "" {
			current = current + "\n"
		}
		current = current + strings.TrimSpace(parts[1])
		ret[parts[0]] = current
	}
	return Event{}
}
