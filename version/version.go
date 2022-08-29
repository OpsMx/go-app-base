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
	"log"

	"github.com/OpsMx/go-app-base/util"
)

// Versions can be set via environment variables at run-time, or
// can be linker-patched directly during the "go build" step.
// Envars will override the variables defined here or patched.
//
// The expectation is that the docker container used in production will
// have these set.  However, if we want to support containers and
// native code, we may want to patch instead.
//
// Patching would also be required if more than one app runs in the same
// container.
var (
	gitBranch = "dev"
	gitHash   = "dev"
)

// FindGitBranch will return the envar, compiled-in var, or "dev" if none set.
func FindGitBranch() string {
	return util.GetEnvar("GIT_BRANCH", gitBranch)
}

// FindGitHash will return the envar, compiled-in var, or "dev" if none set.
func FindGitHash() string {
	return util.GetEnvar("GIT_HASH", gitHash)
}

// ShowGitVersion will log a string using log.Printf() of the found
// git branch and hash.
func ShowGitVersion() {
	log.Printf("GIT Version: %s@%s", FindGitBranch(), FindGitHash())
}
