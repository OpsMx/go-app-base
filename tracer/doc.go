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

// The `tracer` package provides a way to get HTTP and other
// tracing via otel, and export to Jaeger and/or stdout.
// A common pattern for using:
//
// jaegerEndpoint = flag.String("jaeger-endpoint", "", "Jaeger collector endpoint")
// traceToStdout  = flag.Bool("traceToStdout", false, "log traces to stdout")
// traceRatio     = flag.Float64("traceRatio", 0.01, "ratio of traces to create, if incoming request is not traced")
//
// ctx, cancel := context.WithCancel(context.Background())
// defer cancel()
//
// if *jaegerEndpoint != "" {
//	*jaegerEndpoint = util.GetEnvar("JAEGER_TRACE_URL", "")
// }
//
// tracerProvider, err := tracer.NewTracerProvider(*jaegerEndpoint, *traceToStdout, version.GitHash(), appName, *traceRatio)
// util.Check(err)
// defer tracerProvider.Shutdown(ctx)
//
// Catching signals would also be wise, so Shutdown() can be called properly.
