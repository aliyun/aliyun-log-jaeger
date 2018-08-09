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
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/jaegertracing/jaeger/model"
	storageMetrics "github.com/jaegertracing/jaeger/storage/spanstore/metrics"
	"github.com/pkg/errors"
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
	logGroup, err := FromSpan(span, "", "0.0.0.0")
	if err != nil {
		s.logError(span, err, "Failed to convert span to logGroup", s.logger)
	}

	start := time.Now()
	err = s.logstore.PutLogs(logGroup)
	s.writerMetrics.putLogs.Emit(err, time.Since(start))

	if err != nil {
		s.logError(span, err, "Failed to write span", s.logger)
	}
	return err
}

func (s *SpanWriter) logError(span *model.Span, err error, msg string, logger *zap.Logger) error {
	logger.
		With(zap.String("traceID", span.TraceID.String())).
		With(zap.String("spanID", span.SpanID.String())).
		With(zap.Error(err)).
		Error(msg)
	return errors.Wrap(err, msg)
}
