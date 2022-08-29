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
	"encoding/json"
	"net/http"

	"log"
)

// StatusCodeOK returns true if the error code provided is between
// 200 and 2999, inclusive.
func StatusCodeOK(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

// HTTPError defines a simple struct to return JSON formatted error
// messages.
type HttpError struct {
	Status string      `json:"status,omitempty" yaml:"status,omitempty"`
	Code   int         `json:"code,omitempty" yaml:"code,omitempty"`
	Error  interface{} `json:"error,omitempty" yaml:"error,omitempty"`
}

// SetError returns a JSON error message with a Status field set to 'error',
// a 'code' set to the provided code, and the 'error' set to the provided
// content.
//
// Prior to calling, nothing should be written to the writer, and afterwards
// nothing should be written.
//
// The content-type will be set to application/json.
func SetError(w http.ResponseWriter, statusCode int, message interface{}) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	m := HttpError{Status: "error", Code: statusCode, Error: message}
	d, err := json.Marshal(m)
	if err != nil {
		log.Printf("marshalling error json: %v", err)
	}
	l, err := w.Write(d)
	if l != len(d) {
		log.Printf("writing error json: %d of %d bytes written", l, len(d))
		return
	}
	if err != nil {
		log.Printf("writing error json: %v", err)
	}
}
