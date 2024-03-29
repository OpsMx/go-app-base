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

package httputil

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// ClientConfig defines various timeouts we will want to change.
// All times are in seconds.  If 0, a default will be used.
type ClientConfig struct {
	DialTimeout           int `json:"dialTimeout,omitempty" yaml:"dialTimeout,omitempty"`
	ClientTimeout         int `json:"clientTimeout,omitempty" yaml:"clientTimeout,omitempty"`
	TLSHandshakeTimeout   int `json:"tlsHandshakeTimeout,omitempty" yaml:"tlsHandshakeTimeout,omitempty"`
	ResponseHeaderTimeout int `json:"responseHeaderTimeout,omitempty" yaml:"responseHeaderTimeout,omitempty"`
	MaxIdleConnections    int `json:"maxIdleConnections,omitempty" yaml:"maxIdleConnections,omitempty"`
}

var defaultTLSConfig *tls.Config

var defaultClientConfig = &ClientConfig{
	DialTimeout:           15,
	ClientTimeout:         60,
	TLSHandshakeTimeout:   15,
	ResponseHeaderTimeout: 60,
	MaxIdleConnections:    5,
}

func (c *ClientConfig) applyDefaults() {
	if c.DialTimeout == 0 {
		c.DialTimeout = defaultClientConfig.DialTimeout
	}
	if c.ClientTimeout == 0 {
		c.ClientTimeout = defaultClientConfig.ClientTimeout
	}
	if c.TLSHandshakeTimeout == 0 {
		c.TLSHandshakeTimeout = defaultClientConfig.TLSHandshakeTimeout
	}
	if c.ResponseHeaderTimeout == 0 {
		c.ResponseHeaderTimeout = defaultClientConfig.ResponseHeaderTimeout
	}
	if c.MaxIdleConnections == 0 {
		c.MaxIdleConnections = defaultClientConfig.MaxIdleConnections
	}
}

// SetClientConfig will replace the current clientConfig for all future clients
// returned by NewHTTPClient().  Generally, this will be set once, and probably
// not changed per connection.  It is not going to be thread-safe, in that
// setting the config and then calling NewHTTPClient() could be a race.
func SetClientConfig(c ClientConfig) {
	if defaultClientConfig == nil {
		defaultClientConfig = &ClientConfig{}
	}
	*defaultClientConfig = c
	defaultClientConfig.applyDefaults()
}

// SetTLSConfig sets the default TLS configuration used by NewHTTPClient().
// This will generally be set once for adding custom CA roots or other
// configuration used throughout the application.
//
// NewHTTPClient() also allows per-client TLS configuration, if desired.
func SetTLSConfig(tlsconfig *tls.Config) {
	defaultTLSConfig = tlsconfig
}

// NewHTTPClient returns a new http.Client that is configured with
// sane timeouts, a global TLS configuration, and optionally a per-client
// TLS config.
//
// Generally, the global config will have things like custom CA roots,
// and we will want to trust those for every outgoing conneciton.
// A per-client TLS config would be used where we are talking to a
// specific API, and want to insert our certificates or a custom
// CA root for just that connection.
//
// Future changes should allow merging tls configs, so we can add to
// the base default rather than replace it entirely.
func NewHTTPClient(tlsConfig *tls.Config) *http.Client {
	if tlsConfig == nil {
		tlsConfig = defaultTLSConfig
	}
	dialer := net.Dialer{Timeout: time.Duration(defaultClientConfig.DialTimeout) * time.Second}
	client := &http.Client{
		Timeout: time.Duration(defaultClientConfig.ClientTimeout) * time.Second,
		Transport: otelhttp.NewTransport(&http.Transport{
			Dial:                  dialer.Dial,
			DialContext:           dialer.DialContext,
			TLSHandshakeTimeout:   time.Duration(defaultClientConfig.TLSHandshakeTimeout) * time.Second,
			TLSClientConfig:       tlsConfig,
			ResponseHeaderTimeout: time.Duration(defaultClientConfig.ResponseHeaderTimeout) * time.Second,
			ExpectContinueTimeout: time.Second,
			MaxIdleConns:          defaultClientConfig.MaxIdleConnections,
		}),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	return client
}
