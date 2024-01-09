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

package birger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_applyDefaults(t *testing.T) {
	tests := []struct {
		name     string
		provided Config
		want     Config
	}{
		{
			"URL provided isn't overwritten",
			Config{URL: "abc", Token: "abc"},
			Config{
				URL:                    "abc",
				Token:                  "abc",
				UpdateFrequencySeconds: defaultConfig.UpdateFrequencySeconds,
			},
		}, {
			"token isn't overwritten",
			Config{Token: "xyz"},
			Config{
				URL:                    defaultConfig.URL,
				Token:                  "xyz",
				UpdateFrequencySeconds: defaultConfig.UpdateFrequencySeconds,
			},
		}, {
			"UpdateFrequencySeconds provided isn't overwritten",
			Config{UpdateFrequencySeconds: 1234, Token: "abc"},
			Config{
				URL:                    defaultConfig.URL,
				Token:                  "abc",
				UpdateFrequencySeconds: 1234,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.provided.applyDefaults()
			require.Equal(t, tt.want, tt.provided)
		})
	}
}
