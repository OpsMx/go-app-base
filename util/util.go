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

import (
	"log"
	"os"
)

// Check will log an error using log.Fatal() if err is not nil, otherwise
// nothing happens.
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// GetEnvar will return the envar if set, otherwise the default string provided.
func GetEnvar(name string, defaultValue string) string {
	value, found := os.LookupEnv(name)
	if !found {
		return defaultValue
	}
	return value
}

// Contains looks inside a slice to see if a specific element exists.
// This is not fast, and intended only when the length of the slice
// is small as a linear search is used.
func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
