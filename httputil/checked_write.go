/*
 * Copyright 2022 OpsMx, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License")
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package httputil

import (
	"io"
	"log"
)

// CheckedWrite will check to ensure the number of bytes intended to be written
// were written, and log a warning if not.  It will also check the error
// returned and log the error if there was one.
//
// It may be better to try to write the data that wasn't successfully
// written.
//
// This is used to log any errors mostly, rather than try to recover.
func CheckedWrite(w io.Writer, d []byte) {
	l, err := w.Write(d)
	if l != len(d) {
		log.Printf("partial write: %d of %d written", l, len(d))
	}
	if err != nil {
		log.Printf("write failed: %v", err)
	}
}
