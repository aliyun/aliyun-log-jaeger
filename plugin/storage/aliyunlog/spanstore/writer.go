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

package spanstore

import (
	"context"
	"fmt"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/davecgh/go-spew/spew"
	"github.com/gogo/protobuf/proto"
	"github.com/jaegertracing/jaeger/model"
	storageMetrics "github.com/jaegertracing/jaeger/storage/spanstore/metrics"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

type spanWriterMetrics struct {
	putLogs *storageMetrics.WriteMetrics
}

type SpanWriter struct {
	ctx           context.Context
	logstore      *sls.LogStore
	logger        *zap.Logger
	writerMetrics spanWriterMetrics
}

func NewSpanWriter(logstore *sls.LogStore, logger *zap.Logger, metricsFactory metrics.Factory) *SpanWriter {
	ctx := context.Background()

	return &SpanWriter{
		ctx:      ctx,
		logstore: logstore,
		logger:   logger,
		writerMetrics: spanWriterMetrics{
			putLogs: storageMetrics.NewWriteMetrics(metricsFactory, "putLogs"),
		},
	}
}

func (s *SpanWriter) WriteSpan(span *model.Span) error {
	content := s.buildContent(span)
	s.logger.Info(spew.Sdump(content.Content))

	logs := []*sls.Log{
		{
			Time:     proto.Uint32(uint32(span.StartTime.Unix())),
			Contents: content.Content,
		},
	}

	start := time.Now()
	err := s.logstore.PutLogs(&sls.LogGroup{
		Topic:  proto.String("xxx"),
		Source: proto.String("0.0.0.0"),
		Logs:   logs,
	})
	s.writerMetrics.putLogs.Emit(err, time.Since(start))

	if err != nil {
		s.logError(span, err, "send log failed", s.logger)
	}
	return err
}

func (s *SpanWriter) buildContent(span *model.Span) *Content {
	content := newContent()
	content.Add(traceIDField, span.TraceID.String())
	content.Add(spanIDField, span.SpanID.String())
	content.Add(parentSpanIDField, span.ParentSpanID.String())
	content.Add(operationNameField, span.OperationName)
	content.Add(flagsField, fmt.Sprintf("%d", span.Flags))
	content.Add(startTimeField, cast.ToString(span.StartTime.UnixNano()))
	content.Add(durationField, cast.ToString(span.Duration.Nanoseconds()))
	content.Add(serviceNameField, span.Process.ServiceName)

	for _, tag := range span.Tags {
		content.Add(tagsPrefix+tag.Key, tag.AsString())
	}

	for _, tag := range span.Process.Tags {
		content.Add(processTagsPrefix+tag.Key, tag.AsString())
	}
	return content
}

type Content struct {
	Content []*sls.LogContent
}

func newContent() *Content {
	return &Content{
		Content: []*sls.LogContent{},
	}
}

func (c *Content) Add(key string, value string) {
	content := sls.LogContent{
		Key:   proto.String(key),
		Value: proto.String(value),
	}
	c.Content = append(c.Content, &content)
}

func (s *SpanWriter) logError(span *model.Span, err error, msg string, logger *zap.Logger) error {
	logger.
		With(zap.String("traceID", span.TraceID.String())).
		With(zap.String("spanID", span.SpanID.String())).
		With(zap.Error(err)).
		Error(msg)
	return errors.Wrap(err, msg)
}
