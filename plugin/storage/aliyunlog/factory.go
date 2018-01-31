// Copyright (c) 2017 Uber Technologies, Inc.
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

package aliyunlog

import (
	"flag"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/jaegertracing/jaeger/pkg/aliyunlog/config"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	logDepStore "github.com/jaegertracing/jaeger/plugin/storage/aliyunlog/dependencystore"
	logSpanStore "github.com/jaegertracing/jaeger/plugin/storage/aliyunlog/spanstore"
	"github.com/jaegertracing/jaeger/storage/dependencystore"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

// Factory implements storage.Factory for Elasticsearch backend.
type Factory struct {
	Options *Options

	metricsFactory metrics.Factory
	logger         *zap.Logger

	primaryConfig config.LogstoreBuilder
	spanLogstore  *sls.LogStore
	depLogstore   *sls.LogStore
}

// NewFactory creates a new Factory.
func NewFactory() *Factory {
	return &Factory{
		Options: NewOptions("aliyun-log"),
	}
}

// AddFlags implements plugin.Configurable
func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
	f.Options.AddFlags(flagSet)
}

// InitFromViper implements plugin.Configurable
func (f *Factory) InitFromViper(v *viper.Viper) {
	f.Options.InitFromViper(v)
	f.primaryConfig = f.Options.GetPrimary()
}

// Initialize implements storage.Factory
func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	f.metricsFactory, f.logger = metricsFactory, logger

	spanLogstore, err := f.primaryConfig.NewLogstore(config.SpanType)
	if err != nil {
		return err
	}
	f.spanLogstore = spanLogstore

	depLogstore, err := f.primaryConfig.NewLogstore(config.DependencyType)
	if err != nil {
		return err
	}
	f.depLogstore = depLogstore

	return nil
}

// CreateSpanReader implements storage.Factory
func (f *Factory) CreateSpanReader() (spanstore.Reader, error) {
	cfg := f.primaryConfig
	return logSpanStore.NewSpanReader(
		f.spanLogstore,
		f.logger,
		cfg.GetMaxQueryDuration(),
		f.metricsFactory,
	), nil
}

// CreateSpanWriter implements storage.Factory
func (f *Factory) CreateSpanWriter() (spanstore.Writer, error) {
	return logSpanStore.NewSpanWriter(f.spanLogstore, f.logger, f.metricsFactory), nil
}

// CreateDependencyReader implements storage.Factory
func (f *Factory) CreateDependencyReader() (dependencystore.Reader, error) {
	return logDepStore.NewDependencyStore(f.depLogstore, f.logger), nil
}
