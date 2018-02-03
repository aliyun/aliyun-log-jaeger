// Copyright (c) 2018 The Jaeger Authors.
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

package config

import (
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
)

// LogstoreType describes the type of a logstore
type LogstoreType int

const (
	// SpanType indicates the logstore is used to store span
	SpanType LogstoreType = iota
	// DependencyType indicates the logstore is used to store dependency
	DependencyType
)

// Configuration describes the configuration properties needed to connect to an AliCloud Log Service cluster
type Configuration struct {
	Project            string
	Endpoint           string
	AccessKeyID        string
	AccessKeySecret    string
	SpanLogstore       string
	DependencyLogstore string
	MaxQueryDuration   time.Duration
}

// LogstoreBuilder creates new sls.Logstore
type LogstoreBuilder interface {
	NewLogstore(logstoreType LogstoreType) (*sls.LogStore, error)
	GetMaxQueryDuration() time.Duration
}

// ApplyDefaults copies settings from source unless its own value is non-zero.
func (c *Configuration) ApplyDefaults(source *Configuration) {
	if c.Project == "" {
		c.Project = source.Project
	}
	if c.Endpoint == "" {
		c.Endpoint = source.Endpoint
	}
	if c.AccessKeyID == "" {
		c.AccessKeyID = source.AccessKeyID
	}
	if c.AccessKeySecret == "" {
		c.AccessKeySecret = source.AccessKeySecret
	}
	if c.SpanLogstore == "" {
		c.SpanLogstore = source.SpanLogstore
	}
	if c.DependencyLogstore == "" {
		c.DependencyLogstore = source.DependencyLogstore
	}
	if c.MaxQueryDuration == 0 {
		c.MaxQueryDuration = source.MaxQueryDuration
	}
}

func (c *Configuration) NewLogstore(logstoreType LogstoreType) (*sls.LogStore, error) {
	p, err := sls.NewLogProject(
		c.Project,
		c.Endpoint,
		c.AccessKeyID,
		c.AccessKeySecret,
	)
	if err != nil {
		return nil, err
	}

	var logstore *sls.LogStore
	if logstoreType == SpanType {
		logstore, err = p.GetLogStore(c.SpanLogstore)
	} else {
		logstore, err = p.GetLogStore(c.DependencyLogstore)
	}
	if err != nil {
		return nil, err
	}

	return logstore, nil
}

func (c *Configuration) GetMaxQueryDuration() time.Duration {
	return c.MaxQueryDuration
}
