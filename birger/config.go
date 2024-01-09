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
	"log"
	"os"
)

type Config struct {
	URL                    string `json:"url,omitempty" yaml:"url,omitempty"`
	Token                  string `json:"token,omitempty" yaml:"token,omitempty"`
	UpdateFrequencySeconds int    `json:"updateFrequencySeconds,omitempty" yaml:"updateFrequencySeconds,omitempty"`
}

var defaultConfig = Config{
	UpdateFrequencySeconds: 30,
}

func (cc *Config) applyDefaults() {
	if cc.Token == "" {
		t, found := os.LookupEnv("CONTROLLER_TOKEN")
		if !found {
			log.Fatal("no token in config, nor CONTROLLER_TOKEN envar")
		}
		cc.Token = t
	}
	if cc.UpdateFrequencySeconds == 0 {
		cc.UpdateFrequencySeconds = defaultConfig.UpdateFrequencySeconds
	}
}
