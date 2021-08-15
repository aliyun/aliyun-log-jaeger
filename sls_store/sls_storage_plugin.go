package sls_store

import (
	"time"

	slsSdk "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/go-hclog"
	"github.com/jaegertracing/jaeger/storage/dependencystore"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

type SlsJaegerStoragePlugin struct {
	endpoint     string
	accessKeyID  string
	accessSecret string
	project      string
	instance     slsTraceInstance
	maxLookBack  time.Duration
	logger       hclog.Logger
}

func NewSLSStorageForJaegerPlugin(endpoint string, accessKeyID string, accessSecret string,
	project string, instance string, maxLookBack time.Duration, logger hclog.Logger) *SlsJaegerStoragePlugin {
	return &SlsJaegerStoragePlugin{
		endpoint:     endpoint,
		accessKeyID:  accessKeyID,
		accessSecret: accessSecret,
		project:      project,
		instance:     newSlsTraceInstance(project, instance),
		maxLookBack:  maxLookBack,
		logger:       logger,
	}
}

func (s SlsJaegerStoragePlugin) ArchiveSpanReader() spanstore.Reader {
	return &slsSpanReader{
		client:      buildSLSSdkClient(s),
		instance:    s.instance,
		maxLookBack: s.maxLookBack,
		logger:      s.logger,
	}
}

func (s SlsJaegerStoragePlugin) ArchiveSpanWriter() spanstore.Writer {
	return &slsSpanWriter{
		client:      buildSLSSdkClient(s),
		instance:    s.instance,
		maxLookBack: s.maxLookBack,
		logger:      s.logger,
	}
}

func (s SlsJaegerStoragePlugin) SpanReader() spanstore.Reader {
	return &slsSpanReader{
		client:      buildSLSSdkClient(s),
		instance:    s.instance,
		maxLookBack: s.maxLookBack,
		logger:      s.logger,
	}
}

func (s SlsJaegerStoragePlugin) SpanWriter() spanstore.Writer {
	return &slsSpanWriter{
		client:      buildSLSSdkClient(s),
		instance:    s.instance,
		maxLookBack: s.maxLookBack,
		logger:      s.logger,
	}
}

func (s SlsJaegerStoragePlugin) DependencyReader() dependencystore.Reader {
	return &slsDependencyReader{
		client:   buildSLSSdkClient(s),
		instance: s.instance,
		logger:   s.logger,
	}
}

func buildSLSSdkClient(s SlsJaegerStoragePlugin) *slsSdk.Client {
	return &slsSdk.Client{
		Endpoint:        s.endpoint,
		AccessKeyID:     s.accessKeyID,
		AccessKeySecret: s.accessSecret,
		RequestTimeOut:  DefaultRequestTimeOut,
		RetryTimeOut:    DefaultRetryTimeOut,
	}
}
