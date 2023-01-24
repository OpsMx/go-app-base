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

type Config struct {
	URL                    string `json:"url,omitempty" yaml:"url,omitempty"`
	CAPath                 string `json:"caPath,omitempty" yaml:"caPath,omitempty"`
	CertificatePath        string `json:"certificatePath,omitempty" yaml:"certificatePath,omitempty"`
	KeyPath                string `json:"keyPath,omitempty" yaml:"keyPath,omitempty"`
	UpdateFrequencySeconds int    `json:"updateFrequencySeconds,omitempty" yaml:"updateFrequencySeconds,omitempty"`
}

var defaultConfig = Config{
	CAPath:                 "/app/secrets/controller-ca.crt",
	CertificatePath:        "/app/secrets/controller-control/tls.crt",
	KeyPath:                "/app/secrets/controller-control/tls.key",
	UpdateFrequencySeconds: 30,
}

func (cc *Config) applyDefaults() {
	if cc.CAPath == "" {
		cc.CAPath = defaultConfig.CAPath
	}
	if cc.CertificatePath == "" {
		cc.CertificatePath = defaultConfig.CertificatePath
	}
	if cc.KeyPath == "" {
		cc.KeyPath = defaultConfig.KeyPath
	}
	if cc.UpdateFrequencySeconds == 0 {
		cc.UpdateFrequencySeconds = defaultConfig.UpdateFrequencySeconds
	}
}
