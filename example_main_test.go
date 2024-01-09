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

package main_test

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OpsMx/go-app-base/tracer"
	"github.com/OpsMx/go-app-base/util"
	"github.com/OpsMx/go-app-base/version"
)

const (
	appName = "example-app"
)

var (
	// eg, http://localhost:14268/api/traces
	otlpEndpoint  = flag.String("otlp-endpoint", "", "otlp collector endpoint")
	traceToStdout = flag.Bool("traceToStdout", false, "log traces to stdout")
	traceRatio    = flag.Float64("traceRatio", 0.01, "ratio of traces to create, if incoming request is not traced")
	showversion   = flag.Bool("version", false, "show the version and exit")
)

func exiter() {
	time.Sleep(1 * time.Second)
	err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	if err != nil {
		log.Fatalf("KILL failed: %v", err)
	}
}

func Example() {
	// you would not do this, this is just to make test output match.
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
		fmt.Printf("%s", buf.Bytes())
	}()

	log.Printf("%s", version.VersionString())
	flag.Parse()
	if *showversion {
		os.Exit(0)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if *otlpEndpoint != "" {
		*otlpEndpoint = util.GetEnvar("OTLP_HTTP_URL", "")
	}

	tracerProvider, err := tracer.NewTracerProvider(ctx, *otlpEndpoint, *traceToStdout, version.GitHash(), appName, *traceRatio)
	util.Check(err)
	defer tracerProvider.Shutdown(ctx)

	// other stuff here

	// you would not do this...  this is just to ensure we exit our test
	go exiter()

	sig := <-sigchan
	log.Printf("Exiting cleanly due to a signal: %v", sig)

	// Output:
	// version: dev, hash: dev, buildType: unknown
	// Exiting cleanly due to a signal: interrupt
}
