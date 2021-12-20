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
	"github.com/aliyun/aliyun-log-go-sdk/producer"
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
	ctx               context.Context
	client            sls.ClientInterface
	project           string
	logstore          string
	logger            *zap.Logger
	writerMetrics     spanWriterMetrics
	producer          *producer.Producer
	appendTagRuleFile   TagAppendRules
	rewriteKindRuleFile KindRewriteRules
}

func NewSpanWriter(client sls.ClientInterface, project string, logstore string, initResourceFlag bool, logger *zap.Logger, metricsFactory metrics.Factory, appendTagFile string, rewriteKindFile string, appendTagRuleFileFlag bool, rewriteKindRuleFileFlag bool) (*SpanWriter, error) {
	ctx := context.Background()

	if initResourceFlag {
		logger.Info("Prepare to init span writer resource")
		// init LogService resources
		if initResourceFlag {
			err := InitSpanWriterLogstoreResource(client, project, logstore, logger)
			if err != nil {
				logger.Error("Failed to init span writer resource", zap.Error(err))
				return nil, err
			}
		}
		logger.Info("Init span writer resource successfully")
	}

	newClient, _ := client.(*sls.Client)
	producerConfig := producer.GetDefaultProducerConfig()
	producerConfig.AccessKeySecret = newClient.AccessKeySecret
	producerConfig.AccessKeyID = newClient.AccessKeyID
	producerConfig.Endpoint = newClient.Endpoint
	producerInstance := producer.InitProducer(producerConfig)
	producerInstance.Start()
	appendTagRuleFile := initTagAppendRules(appendTagFile, appendTagRuleFileFlag);
	rewriteKindRuleFile := initKindRewriteRules(rewriteKindFile, rewriteKindRuleFileFlag);

	return &SpanWriter{
		ctx:      ctx,
		client:   client,
		project:  project,
		logstore: logstore,
		logger:   logger,
		producer: producerInstance,
		appendTagRuleFile: appendTagRuleFile,
		rewriteKindRuleFile: rewriteKindRuleFile,
		writerMetrics: spanWriterMetrics{
			putLogs: storageMetrics.NewWriteMetrics(metricsFactory, "putLogs"),
		},
	}, nil
}

func (s *SpanWriter) Close() error {
	s.producer.SafeClose()
	return nil
}

func (s *SpanWriter) WriteSpan(span *model.Span) error {
	logGroup, err := FromSpan(span, "", "0.0.0.0", s.appendTagRuleFile, s.rewriteKindRuleFile)
	if err != nil {
		s.logError(span, err, "Failed to convert span to logGroup", s.logger)
	}
	start := time.Now()
	loglist := logGroup.Logs
	topic := *logGroup.Topic
	source := *logGroup.Source
	err = s.producer.SendLogList(s.project, s.logstore, topic, source, loglist)
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
