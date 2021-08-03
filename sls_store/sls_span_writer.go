package sls_store

import (
	"context"
	slsSdk "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/jaegertracing/jaeger/model"
	"time"
)

type slsSpanWriter struct {
	client      *slsSdk.Client
	instance    slsTraceInstance
	maxLookBack time.Duration
}

func (s slsSpanWriter) WriteSpan(ctx context.Context, span *model.Span) error {
	return nil
}
