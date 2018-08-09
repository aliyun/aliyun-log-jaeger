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
func FromSpan(span *model.Span, topic, source string) (*sls.LogGroup, error) {
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

func (c converter) fromSpan(span *model.Span, topic, source string) (*sls.LogGroup, error) {
	logs, err := c.fromSpanToLogs(span)
	if err != nil {
		return nil, err
	}
	return &sls.LogGroup{
		Topic:  proto.String(topic),
		Source: proto.String(source),
		Logs:   logs,
	}, nil
}

func (c converter) fromSpanToLogs(span *model.Span) ([]*sls.Log, error) {
	contents, err := c.fromSpanToLogContents(span)
	if err != nil {
		return nil, err
	}
	return []*sls.Log{
		{
			Time:     proto.Uint32(uint32(span.StartTime.Unix())),
			Contents: contents,
		},
	}, nil
}

func (c converter) fromSpanToLogContents(span *model.Span) ([]*sls.LogContent, error) {
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

	contents, err := c.appendReferences(contents, span.References)
	if err != nil {
		return nil, err
	}

	contents, err = c.appendLogs(contents, span.Logs)
	if err != nil {
		return nil, err
	}

	contents, err = c.appendWarnings(contents, span.Warnings)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (c converter) appendContents(contents []*sls.LogContent, k, v string) []*sls.LogContent {
	content := sls.LogContent{
		Key:   proto.String(k),
		Value: proto.String(v),
	}
	return append(contents, &content)
}

func (c converter) appendReferences(contents []*sls.LogContent, references []model.SpanRef) ([]*sls.LogContent, error) {
	if len(references) < 1 {
		return contents, nil
	}

	r, err := json.Marshal(references)
	if err != nil {
		return nil, err
	}

	content := sls.LogContent{
		Key:   proto.String(referenceField),
		Value: proto.String(string(r)),
	}

	return append(contents, &content), nil
}

func (c converter) appendLogs(contents []*sls.LogContent, logs []model.Log) ([]*sls.LogContent, error) {
	if len(logs) < 1 {
		return contents, nil
	}

	r, err := json.Marshal(logs)
	if err != nil {
		return nil, err
	}

	content := sls.LogContent{
		Key:   proto.String(logsField),
		Value: proto.String(string(r)),
	}

	return append(contents, &content), nil
}

func (c converter) appendWarnings(contents []*sls.LogContent, warnings []string) ([]*sls.LogContent, error) {
	if len(warnings) < 1 {
		return contents, nil
	}

	r, err := json.Marshal(warnings)
	if err != nil {
		return nil, err
	}

	content := sls.LogContent{
		Key:   proto.String(warningsField),
		Value: proto.String(string(r)),
	}

	return append(contents, &content), nil
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
			parentSpanID, err := model.SpanIDFromString(v)
			if err != nil {
				return nil, err
			}
			span.ParentSpanID = parentSpanID
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
		case referenceField:
			refs, err := c.unmarshalReferences(v)
			if err != nil {
				return nil, err
			}
			span.References = refs
		case logsField:
			logs, err := c.unmarshalLogs(v)
			if err != nil {
				return nil, err
			}
			span.Logs = logs
		case warningsField:
			warnings, err := c.unmarshalWarnings(v)
			if err != nil {
				return nil, err
			}
			span.Warnings = warnings
		}
		tags = c.convertTags(tags, tagsPrefix, k, v)
		process.Tags = c.convertTags(process.Tags, processTagsPrefix, k, v)
	}
	span.Tags = tags
	span.Process = &process
	return &span, nil
}

func (c converter) convertTags(tags model.KeyValues, prefix, k, v string) model.KeyValues {
	if strings.HasPrefix(k, prefix) {
		kv := model.String(strings.TrimPrefix(k, prefix), v)
		return append(tags, kv)
	}
	return tags
}

func (c converter) unmarshalReferences(s string) (refs []model.SpanRef, err error) {
	err = json.Unmarshal([]byte(s), &refs)
	if err != nil {
		return nil, err
	}
	return refs, nil
}

func (c converter) unmarshalLogs(s string) (logs []model.Log, err error) {
	err = json.Unmarshal([]byte(s), &logs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

func (c converter) unmarshalWarnings(s string) (warnings []string, err error) {
	err = json.Unmarshal([]byte(s), &warnings)
	if err != nil {
		return nil, err
	}
	return warnings, nil
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
