package sls_store

import (
	slsSdk "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/jaegertracing/jaeger/storage/dependencystore"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"time"
)

const (
	DefaultRetryTimeOut   = 30
	DefaultRequestTimeOut = 30
)

type SlsJaegerStoragePlugin struct {
	endpoint     string
	accessKeyID  string
	accessSecret string
	project      string
	instance     slsTraceInstance
	maxLookBack  time.Duration
}

func NewSLSStorageForJaegerPlugin(endpoint string, accessKeyID string, accessSecret string,
	project string, instance string, maxLookBack time.Duration) *SlsJaegerStoragePlugin {
	return &SlsJaegerStoragePlugin{
		endpoint:     endpoint,
		accessKeyID:  accessKeyID,
		accessSecret: accessSecret,
		project:      project,
		instance:     newSlsTraceInstance(project, instance),
		maxLookBack:  maxLookBack,
	}
}

func (s SlsJaegerStoragePlugin) ArchiveSpanReader() spanstore.Reader {
	return &slsSpanReader{
		client:      buildSLSSdkClient(s),
		instance:    s.instance,
		maxLookBack: s.maxLookBack,
	}
}

func (s SlsJaegerStoragePlugin) ArchiveSpanWriter() spanstore.Writer {
	return &slsSpanWriter{
		client:      buildSLSSdkClient(s),
		instance:    s.instance,
		maxLookBack: s.maxLookBack,
	}
}

func (s SlsJaegerStoragePlugin) SpanReader() spanstore.Reader {
	return &slsSpanReader{
		client:   buildSLSSdkClient(s),
		instance: s.instance,
	}
}

func (s SlsJaegerStoragePlugin) SpanWriter() spanstore.Writer {
	return &slsSpanWriter{
		client:   buildSLSSdkClient(s),
		instance: s.instance,
	}
}

func (s SlsJaegerStoragePlugin) DependencyReader() dependencystore.Reader {
	return &slsDependencyReader{
		client:   buildSLSSdkClient(s),
		instance: s.instance,
	}
}

func buildSLSSdkClient(s SlsJaegerStoragePlugin) *slsSdk.Client {
	return &slsSdk.Client{
		Endpoint:        s.endpoint,
		AccessKeyID:     s.accessKeyID,
		AccessKeySecret: s.accessSecret,
		SecurityToken:   s.accessSecret,
		RequestTimeOut:  DefaultRequestTimeOut,
		RetryTimeOut:    DefaultRetryTimeOut,
	}
}
