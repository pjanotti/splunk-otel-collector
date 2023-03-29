// Copyright Splunk, Inc.
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

package zeroconfig

import (
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/signalfx/splunk-otel-collector/tests/testutils"
	"github.com/signalfx/splunk-otel-collector/tests/testutils/telemetry"
)

const (
	container_version string = "ltsc2019"

	otlpContainerEndpoint string = "host.docker.internal:" + otlpPort
	otlpHostEndpoint      string = "localhost:" + otlpPort
	otlpPort              string = "23456"

	test_container_name string = "windows-iis-test-container"
	test_image_name     string = "windows-iis-test-image"
)

func TestWindowsIISInstrumentation(t *testing.T) {
	// WARNING:
	// 1. Testcontainers for Go doesn't support Windows containers, see https://github.com/testcontainers/testcontainers-go/issues/948
	//    In light of that will issue docker commands via exec.Command.
	// 2. The test uses the default configuration fo the collector that uses signalfx to export metrics.
	//    To avoid building an (expected to be) short lived signalfx sink the test launches another collector
	//    that receives the signals from the instrumented container. This way the test can leverage the existing
	//    OTLP sink.

	// Set the Docker "context"
	os.Chdir(path.Join(".", "testdata"))
	defer os.Chdir("..")

	expectedResourceMetrics, err := telemetry.LoadResourceMetrics("expected_resource_metrics.yaml")
	require.NoError(t, err)

	cmd := exec.Command(
		"docker", "build", ".",
		"-t", test_image_name,
		"-f", "Dockerfile",
		"--build-arg", "windowscontainer_version="+container_version)
	var out strings.Builder
	cmd.Stdout = &out
	err = cmd.Run()
	t.Log(out.String())
	require.NoError(t, err)

	out.Reset()
	cmd = exec.Command("docker", "run", "--rm", "--detach", "-p", "8000:80", "--name", test_container_name, test_image_name)
	cmd.Stdout = &out
	err = cmd.Run()
	container_id := out.String()
	t.Log(container_id)
	require.NoError(t, err)
	defer func() {
		cmd := exec.Command("docker", "stop", test_container_name)
		var out strings.Builder
		cmd.Stdout = &out
		err := cmd.Run()
		t.Logf("Stopping %s ID %s", test_container_name, out.String())
		require.NoError(t, err)
	}()

	otlp, err := testutils.NewOTLPReceiverSink().WithEndpoint(otlpHostEndpoint).Build()
	require.NoError(t, err)
	defer func() {
		require.Nil(t, otlp.Shutdown())
	}()
	require.NoError(t, otlp.Start())

	httpClient := &http.Client{}

	// GET request to the ASP.NET app
	url := "http://localhost:8000/aspnetfxapp/api/values"
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	resp, err := httpClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	t.Log("ASP.NET HTTP Request succeeded")

	require.NoError(t, otlp.AssertAllMetricsReceived(t, *expectedResourceMetrics, 15*time.Second))
}
