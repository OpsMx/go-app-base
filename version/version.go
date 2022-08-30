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

package version

import (
	"fmt"
)

// Versions can be set linker.
//
// For example, to set `gitBranch`, an ldflags of `-X 'github.com/OpsMx/go-app-base/version.gitBranch=v1.4.2'`
// could be used on the 'go build' command line.
//
// These are generally set in the Makefile or Dockerfile from information obtained from the git repo.
var (
	gitBranch = "dev"
	gitHash   = "dev"
	buildType = "unknown"
)

// GitBranch will return the envar, compiled-in var, or "dev" if none set.  Often this is
// a tag, so will have a format such as `v1.0.2` or `v1.0.2-5-g12350123` for changes without
// a specific tag.
//
// This can be used for "short" version strings, while VersionString would be more
// exact.
func GitBranch() string {
	return gitBranch
}

// GitHash will return the envar, compiled-in var, or "dev" if none set.
func GitHash() string {
	return gitHash
}

// BuildType retuns whatever buildTime is set to by the linker.
func BuildType() string {
	return buildType
}

// VersionString returns a formatted version, git hash, and build type.
func VersionString() string {
	return fmt.Sprintf("version: %s, hash: %s, buildType: %s", gitBranch, gitHash, buildType)
}
