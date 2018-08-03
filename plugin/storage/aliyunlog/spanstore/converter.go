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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/proto"
	"github.com/jaegertracing/jaeger/model"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// FromSpan converts a model.Span to a log record
func FromSpan(span *model.Span, topic, source string) *sls.LogGroup {
	return converter{}.fromSpan(span, topic, source)
}

// ToSpan converts a log record to a model.Span
func ToSpan(log map[string]string) (*model.Span, error) {
	return converter{}.toSpan(log)
}

// ToTraces converts logs to []*model.Trace
func ToTraces(logs []map[string]string) ([]*model.Trace, error) {
	return converter{}.toTraces(logs)
}

type converter struct{}

func (c converter) fromSpan(span *model.Span, topic, source string) *sls.LogGroup {
	return &sls.LogGroup{
		Topic:  proto.String(topic),
		Source: proto.String(source),
		Logs:   c.fromSpanToLogs(span),
	}
}

func (c converter) fromSpanToLogs(span *model.Span) []*sls.Log {
	return []*sls.Log{
		{
			Time:     proto.Uint32(uint32(span.StartTime.Unix())),
			Contents: c.fromSpanToLogContents(span),
		},
	}
}

func (c converter) fromSpanToLogContents(span *model.Span) []*sls.LogContent {
	contents := make([]*sls.LogContent, 0)
	contents = c.appendContents(contents, traceIDField, span.TraceID.String())
	contents = c.appendContents(contents, spanIDField, span.SpanID.String())
	contents = c.appendContents(contents, parentSpanIDField, span.ParentSpanID.String())
	contents = c.appendContents(contents, operationNameField, span.OperationName)
	contents = c.appendContents(contents, flagsField, fmt.Sprintf("%d", span.Flags))
	contents = c.appendContents(contents, startTimeField, cast.ToString(span.StartTime.UnixNano()))
	contents = c.appendContents(contents, durationField, cast.ToString(span.Duration.Nanoseconds()))
	contents = c.appendContents(contents, serviceNameField, span.Process.ServiceName)
	for _, tag := range span.Tags {
		contents = c.appendContents(contents, tagsPrefix+tag.Key, tag.AsString())
	}
	for _, tag := range span.Process.Tags {
		contents = c.appendContents(contents, processTagsPrefix+tag.Key, tag.AsString())
	}
	contents = c.appendContents(contents, logsField, c.tryMarshalLogs(span.Logs))

	return contents
}

func (c converter) appendContents(contents []*sls.LogContent, k, v string) []*sls.LogContent {
	content := sls.LogContent{
		Key:   proto.String(k),
		Value: proto.String(v),
	}
	return append(contents, &content)
}

func (c converter) toSpan(log map[string]string) (*model.Span, error) {
	span := model.Span{}
	tags := make([]model.KeyValue, 0)
	process := model.Process{
		Tags: make([]model.KeyValue, 0),
	}

	for k, v := range log {
		switch k {
		case traceIDField:
			traceID, err := model.TraceIDFromString(v)
			if err != nil {
				return nil, err
			}
			span.TraceID = traceID
		case spanIDField:
			spanID, err := model.SpanIDFromString(v)
			if err != nil {
				return nil, err
			}
			span.SpanID = spanID
		case parentSpanIDField:
			ParentSpanID, err := model.SpanIDFromString(v)
			if err != nil {
				return nil, err
			}
			span.ParentSpanID = ParentSpanID
		case operationNameField:
			span.OperationName = v
		case flagsField:
			span.Flags = model.Flags(cast.ToUint64(v))
		case startTimeField:
			span.StartTime = model.EpochMicrosecondsAsTime(cast.ToUint64(v) / 1000)
		case durationField:
			span.Duration = model.MicrosecondsAsDuration(cast.ToUint64(v) / 1000)
		case serviceNameField:
			process.ServiceName = v
		case logsField:
			span.Logs = c.tryUnmarshalLogs(v)
		}
		tags = c.appendTags(tags, tagsPrefix, k, v)
		process.Tags = c.appendTags(process.Tags, processTagsPrefix, k, v)
	}
	span.Tags = tags
	span.Process = &process
	return &span, nil
}

func (c converter) appendTags(tags model.KeyValues, prefix, k, v string) model.KeyValues {
	if strings.HasPrefix(k, prefix) {
		kv := model.String(strings.TrimPrefix(k, prefix), v)
		return append(tags, kv)
	}
	return tags
}

func (c converter) tryUnmarshalLogs(s string) (rv []model.Log) {
	err := json.Unmarshal([]byte(s), &rv)
	if err != nil {
		return nil
	}
	return rv
}

func (c converter) tryMarshalLogs(log []model.Log) string {
	if len(log) < 1 {
		return "[]"
	}

	rv, err := json.Marshal(log)
	if err != nil {
		return "[]"
	}

	return string(rv)
}

func (c converter) toTraces(logs []map[string]string) ([]*model.Trace, error) {
	var traces []*model.Trace
	for _, log := range logs {
		trace, err := c.toTrace(log)
		if err != nil {
			return nil, err
		}
		traces = append(traces, trace)
	}
	return traces, nil
}

func (c converter) toTrace(log map[string]string) (*model.Trace, error) {
	span, err := c.toSpan(log)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert log to span")
	}
	return &model.Trace{
		Spans: []*model.Span{span},
	}, nil
}
