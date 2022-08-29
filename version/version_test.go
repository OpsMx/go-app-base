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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGitBranch(t *testing.T) {
	t.Run("returns gitBranch", func(t *testing.T) {
		gitBranch = "thisIsABranch"
		require.Equal(t, gitBranch, GitBranch())
	})
}

func TestGitHash(t *testing.T) {
	t.Run("returns gitHash", func(t *testing.T) {
		gitHash = "thisIsAHash"
		require.Equal(t, gitHash, GitHash())
	})
}

func TestBuildType(t *testing.T) {
	t.Run("returns buildType", func(t *testing.T) {
		buildType = "thisIsABuildType"
		require.Equal(t, buildType, BuildType())
	})
}

func TestVersionString(t *testing.T) {
	t.Run("returns a string", func(t *testing.T) {
		require.NotEmpty(t, VersionString())
	})
}
