// Copyright Splunk, Inc.
// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build windows
// +build windows

package main

import (
	"fmt"
	"syscall"

	"go.opentelemetry.io/collector/otelcol"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
)

func run(params otelcol.CollectorSettings) error {
	// do not need to supply service name when startup is invoked through Service Control Manager directly
	if err := svc.Run("", otelcol.NewSvcHandler(params)); err != nil {
		errno, ok := err.(syscall.Errno)
		if ok && errno == windows.ERROR_FAILED_SERVICE_CONTROLLER_CONNECT {
			return runInteractive(params)
		}

		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}
