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
			"nothing provided applies all defaults",
			Config{},
			defaultConfig,
		}, {
			"URL provided isn't overwritten",
			Config{URL: "abc"},
			Config{
				URL:                    "abc",
				CAPath:                 defaultConfig.CAPath,
				CertificatePath:        defaultConfig.CertificatePath,
				KeyPath:                defaultConfig.KeyPath,
				UpdateFrequencySeconds: defaultConfig.UpdateFrequencySeconds,
			},
		}, {
			"CAPath provided isn't overwritten",
			Config{CAPath: "abc"},
			Config{
				URL:                    defaultConfig.URL,
				CAPath:                 "abc",
				CertificatePath:        defaultConfig.CertificatePath,
				KeyPath:                defaultConfig.KeyPath,
				UpdateFrequencySeconds: defaultConfig.UpdateFrequencySeconds,
			},
		}, {
			"CertificatePath provided isn't overwritten",
			Config{CertificatePath: "abc"},
			Config{
				URL:                    defaultConfig.URL,
				CAPath:                 defaultConfig.CAPath,
				CertificatePath:        "abc",
				KeyPath:                defaultConfig.KeyPath,
				UpdateFrequencySeconds: defaultConfig.UpdateFrequencySeconds,
			},
		}, {
			"KeyPath provided isn't overwritten",
			Config{KeyPath: "abc"},
			Config{
				URL:                    defaultConfig.URL,
				CAPath:                 defaultConfig.CAPath,
				CertificatePath:        defaultConfig.CertificatePath,
				KeyPath:                "abc",
				UpdateFrequencySeconds: defaultConfig.UpdateFrequencySeconds,
			},
		}, {
			"UpdateFrequencySeconds provided isn't overwritten",
			Config{UpdateFrequencySeconds: 1234},
			Config{
				URL:                    defaultConfig.URL,
				CAPath:                 defaultConfig.CAPath,
				CertificatePath:        defaultConfig.CertificatePath,
				KeyPath:                defaultConfig.KeyPath,
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
