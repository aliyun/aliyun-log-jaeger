package sls_store

import (
	"context"
	"time"

	slsSdk "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/jaegertracing/jaeger/model"
)

type slsSpanWriter struct {
	client      *slsSdk.Client
	instance    slsTraceInstance
	maxLookBack time.Duration
	logger      hclog.Logger
}

func (s slsSpanWriter) WriteSpan(ctx context.Context, span *model.Span) error {
	if contents, err := convertToSpanLog(span, "", "0.0.0.0"); err != nil {
		return nil
	} else {
		return s.client.PutLogs(s.instance.project(), s.instance.traceLogStore(), contents)
	}
}

func convertToSpanLog(span *model.Span, topic, source string) (*slsSdk.LogGroup, error) {
	if logs, err := spanToLog(span); err == nil {
		return &slsSdk.LogGroup{
			Topic:  proto.String(topic),
			Source: proto.String(source),
			Logs:   logs,
		}, nil
	} else {
		return nil, err
	}
}

func spanToLog(span *model.Span) ([]*slsSdk.Log, error) {
	contents, err := dataConvert.ToSLSSpan(span)
	if err != nil {
		return nil, err
	}
	return []*slsSdk.Log{
		{
			Time:     proto.Uint32(uint32(span.StartTime.Unix())),
			Contents: contents,
		},
	}, nil
}
