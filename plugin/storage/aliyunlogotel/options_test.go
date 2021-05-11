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

package aliyunlogotel

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
	assert.False(t, primary.AliCloudK8SFlag)
	assert.Empty(t, primary.AccessKeyID)
	assert.Empty(t, primary.AccessKeySecret)
	assert.Equal(t, "jaeger-span", primary.SpanLogstore)
	assert.Equal(t, 24*time.Hour, primary.MaxQueryDuration)

	aux := opts.Get("archive")
	assert.Equal(t, primary.Project, aux.Project)
	assert.Equal(t, primary.Endpoint, aux.Endpoint)
	assert.Equal(t, primary.AliCloudK8SFlag, aux.AliCloudK8SFlag)
	assert.Equal(t, primary.AccessKeyID, aux.AccessKeyID)
	assert.Equal(t, primary.AccessKeySecret, aux.AccessKeySecret)
	assert.Equal(t, primary.SpanLogstore, aux.SpanLogstore)
}

func TestOptionsWithFlags(t *testing.T) {
	opts := NewOptions("aliyun-log", "aliyun-log.aux")
	v, command := config.Viperize(opts.AddFlags)
	command.ParseFlags([]string{
		"--aliyun-log.project=my-jaeger-test",
		"--aliyun-log.endpoint=cn-beijing.log.aliyuncs.com",
		"--aliyun-log.alicloud-k8s-flag=false",
		"--aliyun-log.access-key-id=id-xxx",
		"--aliyun-log.access-key-secret=secret-xxx",
		"--aliyun-log.span-logstore=jaeger-span-store",
		"--aliyun-log.max-query-duration=48h",
		// a couple overrides
		"--aliyun-log.aux.project=my-jaeger-test-2",
		"--aliyun-log.aux.alicloud-k8s-flag=true",
		"--aliyun-log.aux.access-key-id=id-yyy",
		"--aliyun-log.aux.access-key-secret=secret-yyy",
		"--aliyun-log.aux.max-query-duration=15m",
	})
	opts.InitFromViper(v)

	primary := opts.GetPrimary()
	assert.Equal(t, "my-jaeger-test", primary.Project)
	assert.Equal(t, "cn-beijing.log.aliyuncs.com", primary.Endpoint)
	assert.Equal(t, "id-xxx", primary.AccessKeyID)
	assert.Equal(t, false, primary.AliCloudK8SFlag)
	assert.Equal(t, "secret-xxx", primary.AccessKeySecret)
	assert.Equal(t, "jaeger-span-store", primary.SpanLogstore)
	assert.Equal(t, 48*time.Hour, primary.MaxQueryDuration)

	aux := opts.Get("aliyun-log.aux")
	assert.Equal(t, "my-jaeger-test-2", aux.Project)
	assert.Equal(t, "cn-beijing.log.aliyuncs.com", aux.Endpoint)
	assert.Equal(t, "id-yyy", aux.AccessKeyID)
	assert.Equal(t, true, aux.AliCloudK8SFlag)
	assert.Equal(t, "secret-yyy", aux.AccessKeySecret)
	assert.Equal(t, "jaeger-span-store", aux.SpanLogstore)
	assert.Equal(t, 15*time.Minute, aux.MaxQueryDuration)

}
