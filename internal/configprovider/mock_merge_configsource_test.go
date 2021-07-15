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
	"fmt"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"go.opentelemetry.io/collector/config/configparser"
	"go.opentelemetry.io/collector/config/experimental/configsource"
)

// mergeConfigSource a ConfigSource to mock a possible "merge" config source.
type mergeConfigSource struct {
	// A list of files right now hard-coded on the test but in practice it would be coming
	// form the config source configuration like other components.
	Files []string
}

var _ configsource.ConfigSource = (*mergeConfigSource)(nil)
var _ configsource.Session = (*mergeConfigSource)(nil)

func (m *mergeConfigSource) NewSession(context.Context) (configsource.Session, error) {
	return m, nil
}

func (m *mergeConfigSource) Retrieve(_ context.Context, _ string, _ interface{}) (configsource.Retrieved, error) {
	mergedCfg := koanf.New(configparser.KeyDelimiter)
	for _, fileName := range m.Files {
		fileCfg := koanf.New(configparser.KeyDelimiter)
		if err := fileCfg.Load(file.Provider(fileName), yaml.Parser()); err != nil {
			return nil, fmt.Errorf("unable to read the file %v: %w", fileName, err)
		}
		if err := mergedCfg.Merge(fileCfg); err != nil {
			return nil, fmt.Errorf("error merging file %v into config: %w", fileName, err)
		}
	}

	watchForUpdateFn := func() error {
		// If it supports watching files it could combine them to get the desired effect.
		return configsource.ErrWatcherNotSupported
	}

	return &retrieved{
		value:            mergedCfg.All(),
		watchForUpdateFn: watchForUpdateFn,
	}, nil
}

func (m *mergeConfigSource) RetrieveEnd(context.Context) error {
	return nil
}

func (m *mergeConfigSource) Close(context.Context) error {
	return nil
}
