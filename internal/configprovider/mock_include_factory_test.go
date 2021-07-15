// Copyright Splunk, Inc.
// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configprovider

import (
	"context"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/experimental/configsource"
)

const (
	// The "type" of file config sources in configuration.
	typeStr = "include"
)

type includeFactory struct{}

func (f *includeFactory) Type() config.Type {
	return typeStr
}

func (f *includeFactory) CreateDefaultConfig() ConfigSettings {
	return &Config{
		Settings: NewSettings(typeStr),
	}
}

func (f *includeFactory) CreateConfigSource(_ context.Context, params CreateParams, cfg ConfigSettings) (configsource.ConfigSource, error) {
	return newIncludeConfigSource(params.Logger, cfg.(*Config))
}

// NewFactory creates a factory for include ConfigSource objects.
func NewFactory() Factory {
	return &includeFactory{}
}
