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

package storage

import (
	"flag"
	"fmt"

	"github.com/spf13/viper"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger/plugin"
	"github.com/jaegertracing/jaeger/plugin/storage/aliyunlog"
	"github.com/jaegertracing/jaeger/plugin/storage/aliyunlogotel"
	"github.com/jaegertracing/jaeger/plugin/storage/cassandra"
	"github.com/jaegertracing/jaeger/plugin/storage/es"
	"github.com/jaegertracing/jaeger/plugin/storage/memory"
	"github.com/jaegertracing/jaeger/storage"
	"github.com/jaegertracing/jaeger/storage/dependencystore"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

const (
	aliyunlogStorageType     = "aliyun-log"
	aliyunlogStorageOtelType = "aliyun-log-otel"
	cassandraStorageType     = "cassandra"
	elasticsearchStorageType = "elasticsearch"
	memoryStorageType        = "memory"
)

var allStorageTypes = []string{aliyunlogStorageType, cassandraStorageType, elasticsearchStorageType, memoryStorageType}

// Factory implements storage.Factory interface as a meta-factory for storage components.
type Factory struct {
	FactoryConfig

	factories map[string]storage.Factory
}

// NewFactory creates the meta-factory.
func NewFactory(config FactoryConfig) (*Factory, error) {
	f := &Factory{FactoryConfig: config}
	uniqueTypes := map[string]struct{}{
		f.SpanStorageType:         {},
		f.DependenciesStorageType: {},
	}
	f.factories = make(map[string]storage.Factory)
	for t := range uniqueTypes {
		ff, err := f.getFactoryOfType(t)
		if err != nil {
			return nil, err
		}
		f.factories[t] = ff
	}
	return f, nil
}

func (f *Factory) getFactoryOfType(factoryType string) (storage.Factory, error) {
	switch factoryType {
	case aliyunlogStorageType:
		return aliyunlog.NewFactory(), nil
	case aliyunlogStorageOtelType:
		return aliyunlogotel.NewFactory(), nil
	case cassandraStorageType:
		return cassandra.NewFactory(), nil
	case elasticsearchStorageType:
		return es.NewFactory(), nil
	case memoryStorageType:
		return memory.NewFactory(), nil
	default:
		return nil, fmt.Errorf("Unknown storage type %s. Valid types are %v", factoryType, allStorageTypes)
	}
}

// Initialize implements storage.Factory
func (f *Factory) Initialize(metricsFactory metrics.Factory, logger *zap.Logger) error {
	for _, factory := range f.factories {
		if err := factory.Initialize(metricsFactory, logger); err != nil {
			return err
		}
	}
	return nil
}

// CreateSpanReader implements storage.Factory
func (f *Factory) CreateSpanReader() (spanstore.Reader, error) {
	factory, ok := f.factories[f.SpanStorageType]
	if !ok {
		return nil, fmt.Errorf("No %s backend registered for span store", f.SpanStorageType)
	}
	return factory.CreateSpanReader()
}

// CreateSpanWriter implements storage.Factory
func (f *Factory) CreateSpanWriter() (spanstore.Writer, error) {
	factory, ok := f.factories[f.SpanStorageType]
	if !ok {
		return nil, fmt.Errorf("No %s backend registered for span store", f.SpanStorageType)
	}
	return factory.CreateSpanWriter()
}

// CreateDependencyReader implements storage.Factory
func (f *Factory) CreateDependencyReader() (dependencystore.Reader, error) {
	factory, ok := f.factories[f.DependenciesStorageType]
	if !ok {
		return nil, fmt.Errorf("No %s backend registered for span store", f.DependenciesStorageType)
	}
	return factory.CreateDependencyReader()
}

// AddFlags implements plugin.Configurable
func (f *Factory) AddFlags(flagSet *flag.FlagSet) {
	for _, factory := range f.factories {
		if conf, ok := factory.(plugin.Configurable); ok {
			conf.AddFlags(flagSet)
		}
	}
}

// InitFromViper implements plugin.Configurable
func (f *Factory) InitFromViper(v *viper.Viper) {
	for _, factory := range f.factories {
		if conf, ok := factory.(plugin.Configurable); ok {
			conf.InitFromViper(v)
		}
	}
}

// CreateArchiveSpanReader implements storage.ArchiveFactory
func (f *Factory) CreateArchiveSpanReader() (spanstore.Reader, error) {
	factory, ok := f.factories[f.SpanStorageType]
	if !ok {
		return nil, fmt.Errorf("No %s backend registered for span store", f.SpanStorageType)
	}
	archive, ok := factory.(storage.ArchiveFactory)
	if !ok {
		return nil, storage.ErrArchiveStorageNotSupported
	}
	return archive.CreateArchiveSpanReader()
}

// CreateArchiveSpanWriter implements storage.ArchiveFactory
func (f *Factory) CreateArchiveSpanWriter() (spanstore.Writer, error) {
	factory, ok := f.factories[f.SpanStorageType]
	if !ok {
		return nil, fmt.Errorf("No %s backend registered for span store", f.SpanStorageType)
	}
	archive, ok := factory.(storage.ArchiveFactory)
	if !ok {
		return nil, storage.ErrArchiveStorageNotSupported
	}
	return archive.CreateArchiveSpanWriter()
}
