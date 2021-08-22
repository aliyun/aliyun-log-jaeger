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

	sls "github.com/aliyun/aliyun-log-go-sdk"
)

// LogstoreType describes the type of a logstore
type LogstoreType int

const (
	// SpanType indicates the logstore is used to store span
	SpanType LogstoreType = iota
	// DependencyType indicates the logstore is used to store dependency
	DependencyType
)

const initEcsTokenTryMax = 10
const initEcsTokenSleep = time.Second * 5

// Configuration describes the configuration properties needed to connect to an AliCloud Log Service cluster
type Configuration struct {
	Project            string
	Endpoint           string
	AliCloudK8SFlag    bool
	AccessKeyID        string
	AccessKeySecret    string
	SpanLogstore       string
	SpanAggLogstore    string
	SpanDepLogstore    string
	DependencyLogstore string
	MaxQueryDuration   time.Duration
	InitResourceFlag   bool
	TagAppendRuleFile string
}

// LogstoreBuilder creates new sls.ClientInterface
type LogstoreBuilder interface {
	// NewClient return client, project, logstore, error
	NewClient(logstoreType LogstoreType) (sls.ClientInterface, string, string, string, bool, string, error)
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
	if c.SpanAggLogstore == "" {
		c.SpanAggLogstore = source.SpanAggLogstore
	}
	if c.DependencyLogstore == "" {
		c.DependencyLogstore = source.DependencyLogstore
	}
	if c.MaxQueryDuration == 0 {
		c.MaxQueryDuration = source.MaxQueryDuration
	}
	if c.TagAppendRuleFile == "" {
		c.TagAppendRuleFile = source.TagAppendRuleFile
	}
}

// NewClient return client, project, logstore, error
func (c *Configuration) NewClient(logstoreType LogstoreType) (client sls.ClientInterface, project string, logstore string, aggLogstore string, initResourceFlag bool, tagAppendRuleFile string, err error) {
	if c.AliCloudK8SFlag {
		shutdown := make(chan struct{}, 1)
		for i := 0; i < initEcsTokenTryMax; i++ {
			client, err = sls.CreateTokenAutoUpdateClient(c.Endpoint, UpdateTokenFunction, shutdown)
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, "", "", "", true, "", err
		}

	} else {
		client = sls.CreateNormalInterface(c.Endpoint, c.AccessKeyID, c.AccessKeySecret, "")
	}
	// @todo set user agent
	//p.UserAgent = userAgent
	if logstoreType == SpanType {
		return client, c.Project, c.SpanLogstore, c.SpanAggLogstore, c.InitResourceFlag, c.TagAppendRuleFile, nil
	}
	return client, c.Project, c.DependencyLogstore, c.SpanAggLogstore, c.InitResourceFlag, c.TagAppendRuleFile, nil
}

func (c *Configuration) GetMaxQueryDuration() time.Duration {
	return c.MaxQueryDuration
}
