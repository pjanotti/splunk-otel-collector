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

package telemetry

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v2"
)

// ResourceTraces is a convenience type for test helpers and assertions.
type ResourceTraces struct {
	ResourceSpans []ResourceSpans `yaml:"resource_spans"`
}

// ResourceSpans is the top level grouping of trace data given a Resource (set of attributes)
// and associated ScopeSpans.
type ResourceSpans struct {
	Resource   Resource     `yaml:",inline,omitempty"`
	ScopeSpans []ScopeSpans `yaml:"scope_spans"`
}

// ScopeSpans is the top legel grouping of trace data given InstrumentationScope
// and associated collection of Span instances.
type ScopeSpans struct {
	Scope InstrumentationScope `yaml:"instrumentation_scope,omitempty"`
	Spans []Span               `yaml:"spans,omitempty"`
}

// Span is the trace content, here defined only with the fields required for
// tests. Fields that server as IDs are not present.
type Span struct {
	Name       string          `yaml:"name,omitempty"`
	Attributes *map[string]any `yaml:"attributes,omitempty"`
}

func (rt *ResourceTraces) SaveResourceTraces(path string) error {
	yamlData, err := yaml.Marshal(rt)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, yamlData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func LoadResourceTraces(path string) (*ResourceTraces, error) {
	traceFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer traceFile.Close()

	buffer := new(bytes.Buffer)
	if _, err = buffer.ReadFrom(traceFile); err != nil {
		return nil, err
	}
	by := buffer.Bytes()

	var loaded ResourceTraces
	err = yaml.UnmarshalStrict(by, &loaded)
	if err != nil {
		return nil, err
	}
	// TODO: loaded.FillDefaultValues()
	// TODO: loaded.Validate()
	return &loaded, nil
}
