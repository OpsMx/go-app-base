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
	"fmt"
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

// Read will return an event, which may be empty if nothing but a keep-alive was
// received thus far.  The boolean flag indicates EOF.  If true, no more reads should
// be performed on this SSE.
func (sse *SSE) Read() (Event, bool) {
	ret := Event{}

	for sse.scanner.Scan() {
		line := sse.scanner.Text()
		if line == "" {
			if len(ret) > 0 {
				return ret, false
			}
		}
		if line == ":" {
			if len(ret) == 0 {
				return ret, false
			}
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
	return Event{}, true
}

func (sse *SSE) Write(w io.Writer, event Event) error {
	if len(event) == 0 {
		return nil
	}

	for k, v := range event {
		for _, vv := range strings.Split(v, "\n") {
			s := k + ": " + vv + "\n"
			n, err := w.Write([]byte(s))
			if err != nil {
				return err
			}
			if n != len(s) {
				return fmt.Errorf("short write: %d of %d", n, len(s))
			}
		}
	}
	n, err := w.Write([]byte("\n"))
	if err != nil {
		return err
	}
	if n != 1 {
		return fmt.Errorf("short write: %d of 1", n)
	}
	return nil
}

func (sse *SSE) KeepAlive(w io.Writer) error {
	n, err := w.Write([]byte(":\n"))
	if err != nil {
		return err
	}
	if n != 2 {
		return fmt.Errorf("short write: %d of 2", n)
	}
	return nil

}
