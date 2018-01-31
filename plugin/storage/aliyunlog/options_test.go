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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/jaegertracing/jaeger/pkg/config"
)

func TestOptions(t *testing.T) {
	opts := NewOptions("foo")
	primary := opts.GetPrimary()
	assert.Empty(t, primary.Project)
	assert.Empty(t, primary.Endpoint)
	assert.Empty(t, primary.AccessKeyID)
	assert.Empty(t, primary.AccessKeySecret)
	assert.Equal(t, "jaeger-span", primary.SpanLogstore)
	assert.Equal(t, "jaeger-dependency", primary.DependencyLogstore)
	assert.Equal(t, int(2), primary.LogstoreShardCount)
	assert.Equal(t, int(30), primary.LogstoreShardTTL)
	assert.Equal(t, 24*time.Hour, primary.MaxSpanAge)

	aux := opts.Get("archive")
	assert.Equal(t, primary.Project, aux.Project)
	assert.Equal(t, primary.Endpoint, aux.Endpoint)
	assert.Equal(t, primary.AccessKeyID, aux.AccessKeyID)
	assert.Equal(t, primary.AccessKeySecret, aux.AccessKeySecret)
	assert.Equal(t, primary.SpanLogstore, aux.SpanLogstore)
	assert.Equal(t, primary.DependencyLogstore, aux.DependencyLogstore)
}

func TestOptionsWithFlags(t *testing.T) {
	opts := NewOptions("aliyun-log", "aliyun-log.aux")
	v, command := config.Viperize(opts.AddFlags)
	command.ParseFlags([]string{
		"--aliyun-log.project=my-jaeger-test",
		"--aliyun-log.endpoint=cn-beijing.log.aliyuncs.com",
		"--aliyun-log.access-key-id=id-xxx",
		"--aliyun-log.access-key-secret=secret-xxx",
		"--aliyun-log.span-logstore=jaeger-span-store",
		"--aliyun-log.dependency-logstore=jaeger-dependency-store",
		"--aliyun-log.logstore-shard-count=3",
		"--aliyun-log.logstore-shard-ttl=7",
		"--aliyun-log.max-span-age=48h",
		// a couple overrides
		"--aliyun-log.aux.project=my-jaeger-test-2",
		"--aliyun-log.aux.logstore-shard-ttl=1",
		"--aliyun-log.aux.max-span-age=15m",
	})
	opts.InitFromViper(v)

	primary := opts.GetPrimary()
	assert.Equal(t, "my-jaeger-test", primary.Project)
	assert.Equal(t, "cn-beijing.log.aliyuncs.com", primary.Endpoint)
	assert.Equal(t, "id-xxx", primary.AccessKeyID)
	assert.Equal(t, "secret-xxx", primary.AccessKeySecret)
	assert.Equal(t, "jaeger-span-store", primary.SpanLogstore)
	assert.Equal(t, "jaeger-dependency-store", primary.DependencyLogstore)
	assert.Equal(t, int(3), primary.LogstoreShardCount)
	assert.Equal(t, int(7), primary.LogstoreShardTTL)
	assert.Equal(t, 48*time.Hour, primary.MaxSpanAge)

	aux := opts.Get("aliyun-log.aux")
	assert.Equal(t, "my-jaeger-test-2", aux.Project)
	assert.Equal(t, "cn-beijing.log.aliyuncs.com", aux.Endpoint)
	assert.Equal(t, "id-xxx", aux.AccessKeyID)
	assert.Equal(t, "secret-xxx", aux.AccessKeySecret)
	assert.Equal(t, "jaeger-span-store", aux.SpanLogstore)
	assert.Equal(t, "jaeger-dependency-store", aux.DependencyLogstore)
	assert.Equal(t, int(3), aux.LogstoreShardCount)
	assert.Equal(t, int(1), aux.LogstoreShardTTL)
	assert.Equal(t, 15*time.Minute, aux.MaxSpanAge)

}
