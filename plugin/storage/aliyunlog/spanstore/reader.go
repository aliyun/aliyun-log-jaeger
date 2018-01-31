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
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	storageMetrics "github.com/jaegertracing/jaeger/storage/spanstore/metrics"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-lib/metrics"
	"go.uber.org/zap"
)

const (
	traceIDField       = "traceID"
	spanIDField        = "spanID"
	parentSpanIDField  = "parentSpanID"
	operationNameField = "operationName"
	flagsField         = "flags"
	startTimeField     = "startTime"
	durationField      = "duration"
	tagsPrefix         = "tags."
	serviceNameField   = "process.serviceName"
	processTagsPrefix  = "process.tags."

	defaultMaxLineNumber = 1000 // the default logstore allowed limit
	defaultNumTraces     = 100
)

var (
	// ErrServiceNameNotSet occurs when attempting to query with an empty service name
	ErrServiceNameNotSet = errors.New("Service Name must be set")

	// ErrStartTimeMinGreaterThanMax occurs when start time min is above start time max
	ErrStartTimeMinGreaterThanMax = errors.New("Start Time Minimum is above Maximum")

	// ErrDurationMinGreaterThanMax occurs when duration min is above duration max
	ErrDurationMinGreaterThanMax = errors.New("Duration Minimum is above Maximum")

	// ErrMalformedRequestObject occurs when a request object is nil
	ErrMalformedRequestObject = errors.New("Malformed request object")

	// ErrStartAndEndTimeNotSet occurs when start time and end time are not set
	ErrStartAndEndTimeNotSet = errors.New("Start and End Time must be set")
)

// SpanReader can query for and load traces from AliCloud Log Service
type SpanReader struct {
	ctx      context.Context
	logstore *sls.LogStore
	logger   *zap.Logger
	// The age of the oldest data we will look for.
	maxLookback time.Duration
}

// NewSpanReader returns a new SpanReader with a metrics.
func NewSpanReader(logstore *sls.LogStore, logger *zap.Logger, maxLookback time.Duration, metricsFactory metrics.Factory) spanstore.Reader {
	return storageMetrics.NewReadMetricsDecorator(newSpanReader(logstore, logger, maxLookback), metricsFactory)
}

func newSpanReader(logstore *sls.LogStore, logger *zap.Logger, maxLookback time.Duration) *SpanReader {
	ctx := context.Background()
	return &SpanReader{
		ctx:         ctx,
		logstore:    logstore,
		logger:      logger,
		maxLookback: maxLookback,
	}
}

// GetTrace takes a traceID and returns a Trace associated with that traceID
func (s *SpanReader) GetTrace(traceID model.TraceID) (*model.Trace, error) {
	return nil, nil
}

// GetServices returns all services traced by Jaeger, ordered by frequency
func (s *SpanReader) GetServices() ([]string, error) {
	currentTime := time.Now()
	resp, err := s.logstore.GetLogs(
		"",
		currentTime.Add(-s.maxLookback).Unix(),
		currentTime.Unix(),
		fmt.Sprintf("| select distinct(%s)", serviceNameField),
		defaultMaxLineNumber,
		0,
		false,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Search service failed")
	}
	return logsToStringArray(resp.Logs, serviceNameField)
}

// GetOperations returns all operations for a specific service traced by Jaeger
func (s *SpanReader) GetOperations(service string) ([]string, error) {
	currentTime := time.Now()
	queryExp := fmt.Sprintf("%s:%s | select distinct(%s)", serviceNameField, service, operationNameField)
	resp, err := s.logstore.GetLogs(
		"",
		currentTime.Add(-s.maxLookback).Unix(),
		currentTime.Unix(),
		queryExp,
		defaultMaxLineNumber,
		0,
		false,
	)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Search operation for %s failed", service))
	}
	return logsToStringArray(resp.Logs, serviceNameField)
}

func logsToStringArray(logs []map[string]string, key string) ([]string, error) {
	strings := make([]string, len(logs))
	for i, log := range logs {
		val, ok := log[key]
		if !ok {
			return nil, errors.New(fmt.Sprintf("Cannot found %s in log", key))
		}
		strings[i] = val
	}
	return strings, nil
}

// FindTraces retrieves traces that match the traceQuery
func (s *SpanReader) FindTraces(traceQuery *spanstore.TraceQueryParameters) ([]*model.Trace, error) {
	return nil, nil
}

func validateQuery(p *spanstore.TraceQueryParameters) error {
	if p == nil {
		return ErrMalformedRequestObject
	}
	if p.ServiceName == "" && len(p.Tags) > 0 {
		return ErrServiceNameNotSet
	}
	if p.StartTimeMin.IsZero() || p.StartTimeMax.IsZero() {
		return ErrStartAndEndTimeNotSet
	}
	if p.StartTimeMax.Before(p.StartTimeMin) {
		return ErrStartTimeMinGreaterThanMax
	}
	if p.DurationMin != 0 && p.DurationMax != 0 && p.DurationMin > p.DurationMax {
		return ErrDurationMinGreaterThanMax
	}
	return nil
}
