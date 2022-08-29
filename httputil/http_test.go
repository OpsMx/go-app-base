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
	"crypto/tls"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_StatusCodeOK(t *testing.T) {
	var tests = []struct {
		code int
		want bool
	}{
		{100, false},
		{199, false},
		{200, true},
		{250, true},
		{299, true},
		{304, false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.code)
		t.Run(testname, func(t *testing.T) {
			ans := StatusCodeOK(tt.code)
			require.Equal(t, tt.want, ans)
		})
	}
}

func TestClientConfig_applyDefaults(t *testing.T) {
	tests := []struct {
		name     string
		provided ClientConfig
		wanted   ClientConfig
	}{
		{
			"all defaults",
			ClientConfig{},
			*defaultClientConfig,
		}, {
			"DialTimeout set",
			ClientConfig{DialTimeout: 1234},
			ClientConfig{
				DialTimeout:           1234,
				ClientTimeout:         defaultClientConfig.DialTimeout,
				TLSHandshakeTimeout:   defaultClientConfig.TLSHandshakeTimeout,
				ResponseHeaderTimeout: defaultClientConfig.ResponseHeaderTimeout,
				MaxIdleConnections:    defaultClientConfig.MaxIdleConnections,
			},
		}, {
			"ClientTimeout set",
			ClientConfig{ClientTimeout: 1234},
			ClientConfig{
				DialTimeout:           defaultClientConfig.DialTimeout,
				ClientTimeout:         1234,
				TLSHandshakeTimeout:   defaultClientConfig.TLSHandshakeTimeout,
				ResponseHeaderTimeout: defaultClientConfig.ResponseHeaderTimeout,
				MaxIdleConnections:    defaultClientConfig.MaxIdleConnections,
			},
		}, {
			"TLSHandshakeTimeout set",
			ClientConfig{TLSHandshakeTimeout: 1234},
			ClientConfig{
				DialTimeout:           defaultClientConfig.DialTimeout,
				ClientTimeout:         defaultClientConfig.DialTimeout,
				TLSHandshakeTimeout:   1234,
				ResponseHeaderTimeout: defaultClientConfig.ResponseHeaderTimeout,
				MaxIdleConnections:    defaultClientConfig.MaxIdleConnections,
			},
		}, {
			"ResponseHeaderTimeout set",
			ClientConfig{ResponseHeaderTimeout: 1234},
			ClientConfig{
				DialTimeout:           defaultClientConfig.DialTimeout,
				ClientTimeout:         defaultClientConfig.DialTimeout,
				TLSHandshakeTimeout:   defaultClientConfig.TLSHandshakeTimeout,
				ResponseHeaderTimeout: 1234,
				MaxIdleConnections:    defaultClientConfig.MaxIdleConnections,
			},
		}, {
			"MaxIdleConnections set",
			ClientConfig{MaxIdleConnections: 1234},
			ClientConfig{
				DialTimeout:           defaultClientConfig.DialTimeout,
				ClientTimeout:         defaultClientConfig.DialTimeout,
				TLSHandshakeTimeout:   defaultClientConfig.TLSHandshakeTimeout,
				ResponseHeaderTimeout: defaultClientConfig.ResponseHeaderTimeout,
				MaxIdleConnections:    1234,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := tt.wanted
			found.applyDefaults()
			require.Equal(t, tt.wanted, found)
		})
	}
}

func Test_SetClientConfig(t *testing.T) {
	t.Run("sets", func(t *testing.T) {
		defaultClientConfig = nil
		c := ClientConfig{DialTimeout: 9876}
		SetClientConfig(c)
		require.Equal(t, 9876, defaultClientConfig.DialTimeout)
	})
}

func Test_SetTLSConfig(t *testing.T) {
	t.Run("sets", func(t *testing.T) {
		defaultTLSConfig = nil
		c := tls.Config{}
		SetTLSConfig(&c)
		require.NotNil(t, defaultTLSConfig)
	})
}
