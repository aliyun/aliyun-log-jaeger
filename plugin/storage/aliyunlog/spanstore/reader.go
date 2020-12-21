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
	"strconv"
	"strings"
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
	referenceField     = "reference"
	flagsField         = "flags"
	startTimeField     = "startTime"
	durationField      = "duration"
	tagsPrefix         = "tags."
	logsField          = "logs"
	warningsField      = "Warnings"
	serviceNameField   = "process.serviceName"
	processTagsPrefix  = "process.tags."

	defaultServiceLimit    = 1000
	defaultOperationLimit  = 1000
	defaultPageSizeForSpan = 100
	defaultNumTraces       = 100
	defaultMaxNum          = 100000

	emptyTopic = ""

	firstColumn = "_col0"

	progressComplete   = "Complete"
	progressIncomplete = "InComplete"

	querySuffixTemplate = `| select {traceIDField}, max_by("{serviceNameField}", {durationField}) as "{serviceNameField}", max_by({operationNameField}, {durationField}) as {operationNameField}, max_by({startTimeField}, {durationField}) as {startTimeField}, max_by({durationField}, {durationField}) as {durationField} from (select {traceIDField}, "{serviceNameField}", {operationNameField}, {startTimeField}, {durationField} from log limit {maxLineNum}) group by {traceIDField} limit {numTraces}`
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
	ctx         context.Context
	client      sls.ClientInterface
	project     string
	logstore    string
	aggLogstore string
	logger      *zap.Logger
	// The age of the oldest data we will look for.
	maxLookback time.Duration
}

// NewSpanReader returns a new SpanReader with a metrics.
func NewSpanReader(client sls.ClientInterface,
	project string,
	logstore string,
	aggLogstore string,
	logger *zap.Logger,
	maxLookback time.Duration,
	metricsFactory metrics.Factory) spanstore.Reader {
	return storageMetrics.NewReadMetricsDecorator(newSpanReader(client, project, logstore, aggLogstore, logger, maxLookback), metricsFactory)
}

func newSpanReader(client sls.ClientInterface,
	project string,
	logstore string,
	aggLogstore string,
	logger *zap.Logger,
	maxLookback time.Duration) *SpanReader {
	ctx := context.Background()
	return &SpanReader{
		ctx:         ctx,
		client:      client,
		project:     project,
		logstore:    logstore,
		aggLogstore: aggLogstore,
		logger:      logger,
		maxLookback: maxLookback,
	}
}

// GetTrace takes a traceID and returns a Trace associated with that traceID
func (s *SpanReader) GetTrace(traceID model.TraceID) (*model.Trace, error) {
	currentTime := time.Now()
	from := currentTime.Add(-s.maxLookback).Unix()
	to := currentTime.Unix()
	return s.getTrace(traceID.String(), from, to)
}

func (s *SpanReader) getTrace(traceID string, from, to int64) (*model.Trace, error) {
	from = from - 777600
	s.logger.Info(
		"Trying to get trace",
		zap.String("traceID", traceID),
		zap.Int64("from", from),
		zap.Int64("to", to),
	)

	topic := emptyTopic
	queryExp := fmt.Sprintf(`%s: "%s"`, traceIDField, traceID)
	maxLineNum := int64(defaultPageSizeForSpan)
	offset := int64(0)
	reverse := false

	s.logger.Info("Use fix page size")

	count, err := s.getSpansCountForTrace(traceID, topic, from, to)
	if err != nil {
		return nil, err
	}
	if maxLineNum > 100 {
		maxLineNum = 100
	}
	spans := make([]*model.Span, 0)
	for {
		s.logGetLogsParameters(topic, from, to, queryExp, maxLineNum, offset, reverse,
			fmt.Sprintf("Trying to get spans for trace %s", traceID))
		resp, err := func() (*sls.GetLogsResponse, error) {
			var innerErr error
			var lastNotNilResp *sls.GetLogsResponse
			var innerResp *sls.GetLogsResponse
			for i := 0; i < 3; i++ {
				innerResp, innerErr = s.client.GetLogs(s.project, s.logstore, topic, from, to, queryExp, maxLineNum, offset, reverse)
				if innerResp != nil {
					lastNotNilResp = innerResp
				}
				if innerErr == nil && innerResp != nil && innerResp.IsComplete() {
					break
				}
			}
			return lastNotNilResp, innerErr

		}()
		if err != nil && resp == nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Search spans for trace %s failed", traceID))
		}
		for _, log := range resp.Logs {
			span, err := ToSpan(log)
			if err != nil {
				return nil, err
			}
			spans = append(spans, span)
		}
		offset += maxLineNum
		if offset >= count || offset >= 2000 {
			break
		}
	}
	if len(spans) == 0 {
		return nil, spanstore.ErrTraceNotFound
	}
	trace := model.Trace{
		Spans: spans,
	}

	return &trace, nil
}

func (s *SpanReader) getSpansCountForTrace(traceID, topic string, from, to int64) (int64, error) {
	queryExp := fmt.Sprintf(`%s: "%s"`, traceIDField, traceID)

	s.logGetHistograms(topic, from, to, queryExp, fmt.Sprintf("Trying to get count of spans for trace %s", traceID))

	resp, err := s.client.GetHistograms(s.project, s.logstore, topic, from, to, queryExp)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("Failed to get spans count for trace %s", traceID))
	}
	return resp.Count, nil
}

// GetServices returns all services traced by Jaeger, ordered by frequency
func (s *SpanReader) GetServices() ([]string, error) {
	topic := emptyTopic
	currentTime := time.Now()
	from := currentTime.Add(-s.maxLookback).Unix()
	to := currentTime.Unix()
	maxLineNum := int64(0)
	offset := int64(0)
	reverse := false
	aggServiceNames := make([]string, 0)
	if s.aggLogstore != "" {
		queryExp := "__tag__:type:service | select distinct(serviceName) limit 10000"
		s.logGetLogsParameters(topic, from, to, queryExp, maxLineNum, offset, reverse,
			"Trying to get services from agg logstore")
		aggResp, err := s.client.GetLogs(s.project, s.aggLogstore, topic, from, to, queryExp, maxLineNum, offset, reverse)
		if err != nil {
			s.logger.With(zap.Error(err)).Error("Failed to get services from agg logstore")
			return nil, errors.Wrap(err, "Failed to get services from agg logstore")
		}
		s.logProgressIncomplete(topic, from, to, queryExp, maxLineNum, offset, reverse, aggResp.Progress)
		aggServiceNames, err = logsToStringArray(aggResp.Logs, "serviceName")
		if err != nil {
			s.logger.With(zap.Error(err)).Error("Failed to convert logs to string array")
			return nil, errors.Wrap(err, "Failed to convert logs to string array")
		}
		from = to - int64(900)
	}

	queryExp := fmt.Sprintf(
		`| select distinct("%s") from (select "%s" from log limit 100000) limit %d`,
		serviceNameField,
		serviceNameField,
		defaultServiceLimit,
	)
	s.logGetLogsParameters(topic, from, to, queryExp, maxLineNum, offset, reverse,
		"Trying to get services from span logstore")

	resp, err := s.client.GetLogs(s.project, s.logstore, topic, from, to, queryExp, maxLineNum, offset, reverse)
	if err != nil {
		s.logger.With(zap.Error(err)).Error("Failed to get services from span logstore")
		return nil, errors.Wrap(err, "Failed to get services from span logstore")
	}
	s.logProgressIncomplete(topic, from, to, queryExp, maxLineNum, offset, reverse, resp.Progress)

	serviceNames, err := logsToStringArray(resp.Logs, serviceNameField)
	if err != nil {
		s.logger.With(zap.Error(err)).Error("Failed to convert logs to string array")
		return nil, errors.Wrap(err, "Failed to convert logs to string array")
	}
	s.logger.Info(
		"Get services successfully",
		zap.Int("len(aggServiceNames)", len(aggServiceNames)),
		zap.Int("len(serviceNames)", len(serviceNames)),
	)

	return unionStringArray(aggServiceNames, serviceNames)
}

// GetOperations returns all operations for a specific service traced by Jaeger
func (s *SpanReader) GetOperations(service string) ([]string, error) {
	topic := emptyTopic
	currentTime := time.Now()
	from := currentTime.Add(-s.maxLookback).Unix()
	to := currentTime.Unix()
	maxLineNum := int64(0)
	offset := int64(0)
	reverse := false
	aggOperationNames := make([]string, 0)
	if s.aggLogstore != "" {
		queryExp := fmt.Sprintf(
			`__tag__:type:process and serviceName: "%s" | select distinct(operationName) limit 10000 `,
			service)
		s.logGetLogsParameters(topic, from, to, queryExp, maxLineNum, offset, reverse,
			fmt.Sprintf("Trying to get operations for service %s from agg logstore", service))
		aggResp, err := s.client.GetLogs(s.project, s.aggLogstore, topic, from, to, queryExp, maxLineNum, offset, reverse)
		if err != nil {
			s.logger.With(zap.Error(err)).Error(fmt.Sprintf("Failed to get operations for service %s from agg logstore", service))
			return nil, errors.Wrap(err, fmt.Sprintf("Failed to get operations for service %s from agg logstore", service))
		}
		s.logProgressIncomplete(topic, from, to, queryExp, maxLineNum, offset, reverse, aggResp.Progress)
		aggOperationNames, err = logsToStringArray(aggResp.Logs, "operationName")
		if err != nil {
			s.logger.With(zap.Error(err)).Error("Failed to convert logs to string array")
			return nil, errors.Wrap(err, "Failed to convert logs to string array")
		}
		from = to - int64(900)
	}

	queryExp := fmt.Sprintf(
		`%s: "%s" | select distinct(%s) from (select %s from log limit 100000) limit %d`,
		serviceNameField,
		service,
		operationNameField,
		operationNameField,
		defaultOperationLimit,
	)

	s.logGetLogsParameters(topic, from, to, queryExp, maxLineNum, offset, reverse,
		fmt.Sprintf("Trying to get operations for service %s", service))

	resp, err := s.client.GetLogs(s.project, s.logstore, topic, from, to, queryExp, maxLineNum, offset, reverse)
	if err != nil {
		s.logger.With(zap.Error(err)).Error(fmt.Sprintf("Failed to get operations for service %s from span logstore", service))
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to get operations for service %s from span logstore", service))
	}
	s.logProgressIncomplete(topic, from, to, queryExp, maxLineNum, offset, reverse, resp.Progress)
	operationNames, err := logsToStringArray(resp.Logs, operationNameField)
	if err != nil {
		s.logger.With(zap.Error(err)).Error("Failed to convert logs to string array")
		return nil, errors.Wrap(err, "Failed to convert logs to string array")
	}
	s.logger.Info(
		"Get operations successfully",
		zap.Int("len(aggOperationNames)", len(aggOperationNames)),
		zap.Int("len(operationNames)", len(operationNames)),
	)
	return unionStringArray(aggOperationNames, operationNames)
}

func unionStringArray(l, r []string) ([]string, error) {
	if l == nil {
		return r, nil
	}
	if r == nil {
		return l, nil
	}
	m := make(map[string]bool)
	for _, item := range l {
		m[item] = true
	}
	for _, item := range r {
		if _, ok := m[item]; !ok {
			l = append(l, item)
		}
	}
	return l, nil
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
	if err := validateQuery(traceQuery); err != nil {
		return nil, err
	}
	if traceQuery.NumTraces == 0 {
		traceQuery.NumTraces = defaultNumTraces
	}
	return s.findTraces(traceQuery)
}

func (s *SpanReader) multiRead(traceIDs []string, from, to int64) ([]*model.Trace, error) {
	if len(traceIDs) == 0 {
		return []*model.Trace{}, nil
	}

	var traces []*model.Trace
	for _, traceID := range traceIDs {
		trace, err := s.getTrace(traceID, from, to)
		if err != nil {
			return nil, err
		}
		traces = append(traces, trace)
	}
	return traces, nil
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

func (s *SpanReader) findTraceIDs(traceQuery *spanstore.TraceQueryParameters) ([]string, error) {
	query := s.buildFindTracesQuery(traceQuery)

	topic := emptyTopic
	from := traceQuery.StartTimeMin.Unix()
	to := traceQuery.StartTimeMax.Unix() + 1
	queryExp := fmt.Sprintf("| select distinct(%s)", traceIDField)
	if len(query) > 0 {
		queryExp += " " + query
	}
	queryExp += fmt.Sprintf(" limit %d", traceQuery.NumTraces)
	maxLineNum := int64(0)
	offset := int64(0)
	reverse := false

	s.logGetLogsParameters(topic, from, to, queryExp, maxLineNum, offset, reverse, "Trying to find trace ids")

	resp, err := s.client.GetLogs(s.project, s.logstore, topic, from, to, queryExp, maxLineNum, offset, reverse)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to find trace ids")
	}

	return logsToStringArray(resp.Logs, traceIDField)
}

func (s *SpanReader) findTraces(traceQuery *spanstore.TraceQueryParameters) ([]*model.Trace, error) {
	topic := emptyTopic
	from := traceQuery.StartTimeMin.Unix()
	to := traceQuery.StartTimeMax.Unix() + 1
	queryExp := s.buildFindTracesQuery(traceQuery)
	maxLineNum := int64(0)
	offset := int64(0)
	reverse := false

	s.logGetLogsParameters(topic, from, to, queryExp, maxLineNum, offset, reverse, "Trying to find traces")

	resp, err := s.client.GetLogs(s.project, s.logstore, topic, from, to, queryExp, maxLineNum, offset, reverse)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to find traces")
	}

	return ToTraces(resp.Logs)
}

func (s *SpanReader) buildFindTracesQuery(traceQuery *spanstore.TraceQueryParameters) string {
	var subQueries []string

	//add process.serviceName query
	if traceQuery.ServiceName != "" {
		serviceNameQuery := s.buildServiceNameQuery(traceQuery.ServiceName)
		subQueries = append(subQueries, serviceNameQuery)
	}

	//add operationName query
	if traceQuery.OperationName != "" {
		operationNameQuery := s.buildOperationNameQuery(traceQuery.OperationName)
		subQueries = append(subQueries, operationNameQuery)
	}

	//add duration query
	if traceQuery.DurationMax != 0 || traceQuery.DurationMin != 0 {
		durationQuery := s.buildDurationQuery(traceQuery.DurationMin, traceQuery.DurationMax)
		subQueries = append(subQueries, durationQuery)
	}

	for k, v := range traceQuery.Tags {
		tagQuery := s.buildTagQuery(k, v)
		subQueries = append(subQueries, tagQuery)
	}

	query := s.combineSubQueries(subQueries)
	if query != "" {
		query += " "
	}
	query += s.getQuerySuffix(defaultMaxNum, traceQuery.NumTraces)

	return query
}

func (s *SpanReader) buildServiceNameQuery(serviceName string) string {
	return fmt.Sprintf(`%s: "%s"`, serviceNameField, serviceName)
}

func (s *SpanReader) buildOperationNameQuery(operationName string) string {
	return fmt.Sprintf(`%s: "%s"`, operationNameField, operationName)
}

func (s *SpanReader) buildDurationQuery(durationMin time.Duration, durationMax time.Duration) string {
	minDurationMicros := durationMin.Nanoseconds() / 1000
	maxDurationMicros := durationMax.Nanoseconds() / 1000
	if minDurationMicros != 0 && maxDurationMicros != 0 {
		return fmt.Sprintf(
			"%s >= %d and %s <= %d",
			durationField,
			minDurationMicros,
			durationField,
			maxDurationMicros,
		)
	} else if minDurationMicros != 0 {
		return fmt.Sprintf(
			"%s >= %d",
			durationField,
			minDurationMicros,
		)
	} else if maxDurationMicros != 0 {
		return fmt.Sprintf(
			"%s <= %d",
			durationField,
			maxDurationMicros,
		)
	} else {
		return ""
	}
}

func (s *SpanReader) buildTagQuery(k string, v string) string {
	if strings.ContainsAny(v, ": ") && !strings.ContainsAny(v, "'\"") {
		return fmt.Sprintf(`%s: "%s"`, tagsPrefix+k, v)
	} else {
		return fmt.Sprintf(`%s: %s`, tagsPrefix+k, v)
	}
}

func (s *SpanReader) combineSubQueries(subQueries []string) string {
	query := ""
	for _, subQuery := range subQueries {
		if query != "" {
			query += " and "
		}
		if subQuery != "" {
			query += subQuery
		}
	}
	return query
}

func (s *SpanReader) getQuerySuffix(maxLineNum int64, numTraces int) string {
	r := strings.NewReplacer("{traceIDField}", traceIDField,
		"{serviceNameField}", serviceNameField,
		"{operationNameField}", operationNameField,
		"{durationField}", durationField,
		"{startTimeField}", startTimeField,
		"{maxLineNum}", strconv.FormatInt(maxLineNum, 10),
		"{numTraces}", strconv.Itoa(numTraces),
	)
	return r.Replace(querySuffixTemplate)
}

func (s *SpanReader) logGetLogsParameters(topic string, from int64, to int64, queryExp string, maxLineNum int64, offset int64, reverse bool, msg string) {
	s.logger.
		With(zap.String("topic", topic)).
		With(zap.Int64("from", from)).
		With(zap.Int64("to", to)).
		With(zap.String("queryExp", queryExp)).
		With(zap.Int64("maxLineNum", maxLineNum)).
		With(zap.Int64("offset", offset)).
		With(zap.Bool("reverse", reverse)).
		Info(msg)
}

func (s *SpanReader) logGetHistograms(topic string, from int64, to int64, queryExp string, msg string) {
	s.logger.
		With(zap.String("topic", topic)).
		With(zap.Int64("from", from)).
		With(zap.Int64("to", to)).
		With(zap.String("queryExp", queryExp)).
		Info(msg)
}

func (s *SpanReader) logProgressIncomplete(topic string, from int64, to int64, queryExp string, maxLineNum int64, offset int64, reverse bool, progress string) {
	if progress == progressIncomplete {
		s.logger.
			With(zap.String("topic", topic)).
			With(zap.Int64("from", from)).
			With(zap.Int64("to", to)).
			With(zap.String("queryExp", queryExp)).
			With(zap.Int64("maxLineNum", maxLineNum)).
			With(zap.Int64("offset", offset)).
			With(zap.Bool("reverse", reverse)).
			Warn("The response for GetLogs is incomplete")
	}
}
