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

package dependencystore

import (
	"context"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/jaegertracing/jaeger/model"
	"go.uber.org/zap"
)

// DependencyStore handles all queries and insertions to AliCloud Log Service dependencies
type DependencyStore struct {
	ctx      context.Context
	client   sls.ClientInterface
	project  string
	logstore string
	logger   *zap.Logger
}

// NewDependencyStore returns a DependencyStore
func NewDependencyStore(client sls.ClientInterface, project string, logstore string, logger *zap.Logger) *DependencyStore {
	return &DependencyStore{
		ctx:      context.Background(),
		client:   client,
		project:  project,
		logstore: logstore,
		logger:   logger,
	}
}

// WriteDependencies implements dependencystore.Writer#WriteDependencies.
func (s *DependencyStore) WriteDependencies(ts time.Time, dependencies []model.DependencyLink) error {
	return nil
}

// GetDependencies returns all interservice dependencies
func (s *DependencyStore) GetDependencies(endTs time.Time, lookback time.Duration) ([]model.DependencyLink, error) {
	var retDependencies []model.DependencyLink
	return retDependencies, nil
}
