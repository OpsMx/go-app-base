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

package tracer

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
)

// TracerProvider holds the state for a provider, and provides an
// easy way to shut it down without exposing the specific type
// of tracer returned.
type TracerProvider struct {
	Provider *tracesdk.TracerProvider
}

// NewTracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
//
// If no error is returned, `defer provider.Shutdown(ctx)` should be set up
// to ensure flushing occurs.
//
// If the Jaeger URL is empty, the OpenTelemetry
// TracerProvider will be configured to not report to Jaeger.
//
// If traceToStdout is true, traces will be sent to stdout.
func NewTracerProvider(jaegerURL string, traceToStdout bool, githash string, appname string, traceRatio float64) (*TracerProvider, error) {
	res, err := resource.New(context.Background(),
		// add detectors here if needed
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(appname),
			semconv.ServiceVersionKey.String(githash),
		))
	if err != nil {
		log.Fatalf("resource.New: %v", err)
	}

	opts := []tracesdk.TracerProviderOption{
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(traceRatio))),
	}

	exporterCount := 0
	if jaegerURL != "" {
		exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerURL)))
		if err != nil {
			return nil, err
		}
		opts = append(opts, tracesdk.WithBatcher(exp))
		exporterCount++
	}

	if traceToStdout {
		exp, err := newConsoleExporter(os.Stdout)
		if err != nil {
			return nil, err
		}
		opts = append(opts, tracesdk.WithBatcher(exp))
		exporterCount++
	}

	if exporterCount == 0 {
		exp, err := newNullExporter()
		if err != nil {
			return nil, err
		}
		opts = append(opts, tracesdk.WithBatcher(exp))
	}

	tp := tracesdk.NewTracerProvider(opts...)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(propagation.TraceContext{})
	return &TracerProvider{Provider: tp}, nil
}

// Shutdown should be deferred immediately after NewTracerProvider()
// when no error is returned.  This will ensure that on app termination
// it will flush any buffered traces, if possible.  A maximum time
// of 5 seconds will be allowed before we give up, to prevent a hang
// at shutdown.
func (p *TracerProvider) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	if err := p.Provider.Shutdown(ctx); err != nil {
		log.Printf("shutting down tracing: %v", err)
	}
}

func newConsoleExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		stdouttrace.WithPrettyPrint(),
	)
}

// dummy null exporter
func newNullExporter() (trace.SpanExporter, error) {
	return nullExporter{}, nil
}

type nullExporter struct{}

func (nullExporter) ExportSpans(_ context.Context, _ []trace.ReadOnlySpan) error {
	return nil
}

func (nullExporter) Shutdown(_ context.Context) error {
	return nil
}
