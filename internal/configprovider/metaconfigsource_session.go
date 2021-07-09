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

package configprovider

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"go.opentelemetry.io/collector/config/configparser"
	"go.opentelemetry.io/collector/config/experimental/configsource"
)

// Private error types to help with testability.
type (
	errInvalidRetrieveParams struct{ error }
)

type retrieveParams struct {
	// Sources list all the invocations of other config sources that is going to be used to
	// perform the "meta" operation with them.
	Sources []interface{} `mapstructure:"sources"`
}

// metaCfgSrcSession implements the configsource.Session interface.
type metaCfgSrcSession struct {
	manager *Manager
}

var _ configsource.Session = (*metaCfgSrcSession)(nil)

func (m *metaCfgSrcSession) Retrieve(ctx context.Context, selector string, params interface{}) (configsource.Retrieved, error) {
	if selector != "merge" {
		return nil, fmt.Errorf("invalid selector '%s'", selector)
	}

	// TODO: bail on params == nil ?

	actualParams := retrieveParams{}
	paramsParser := configparser.NewParserFromStringMap(cast.ToStringMap(params))
	if err := paramsParser.UnmarshalExact(&actualParams); err != nil {
		return nil, &errInvalidRetrieveParams{fmt.Errorf("failed to unmarshall retrieve params: %w", err)}
	}

	// Resolve any config source usage in the sources using the manager
	// it is expected that they will return map[string]interface{}
	resolvedParser, err := m.manager.Resolve(ctx, paramsParser)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve sources used by meta config source: %w", err)
	}
	resolvedStruct := struct {
		Sources []map[string]interface{} `mapstructure:"sources"`
	}{}
	if err := resolvedParser.UnmarshalExact(&resolvedStruct); err != nil {
		return nil, &errInvalidRetrieveParams{fmt.Errorf("failed to unmarshall config sources used by meta config source: %w", err)}
	}

	// Combine the configurations into a single one and return it.
	res := configparser.NewParser()
	for i, currMap := range resolvedStruct.Sources {
		currParser := configparser.NewParserFromStringMap(currMap)

		// Use AllKeys() so we loop over the "flattened keys", i.e. keys that actually hold value.
		for _, key := range currParser.AllKeys() {
			if res.IsSet(key) {
				// TODO: log here that the previous value is going to be overwritten for now stdout
				fmt.Printf("Config [%s] overwritten by config source %d\n", key, i)
			}
			res.Set(key, currParser.Get(key))
		}
		if err != nil {
			return nil, fmt.Errorf("failed to parse config source on entry %d passed to the meta config source: %w", i, err)
		}
	}

	return NewRetrieved(res.ToStringMap(), m.manager.WatchForUpdate), nil
}

func (m *metaCfgSrcSession) RetrieveEnd(context.Context) error {
	return nil
}

func (m *metaCfgSrcSession) Close(ctx context.Context) error {
	return m.manager.Close(ctx)
}

func newSession(cfgSources map[string]configsource.ConfigSource) *metaCfgSrcSession {
	return &metaCfgSrcSession{
		manager: newManager(cfgSources),
	}
}
