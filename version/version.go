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
var (
	gitBranch = "dev"
	gitHash   = "dev"
	buildType = "unknown"
)

// GitBranch will return the envar, compiled-in var, or "dev" if none set.
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
