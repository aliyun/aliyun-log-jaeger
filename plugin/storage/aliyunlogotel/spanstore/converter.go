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
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/proto"
	"github.com/jaegertracing/jaeger/model"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

// FromSpan converts a model.Span to a log record
func FromSpan(span *model.Span, topic, source string, file TagAppendRules, rule KindRewriteRules) (*sls.LogGroup, error) {
	return converter{}.fromSpan(span, topic, source, file, rule)
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

func (c converter) fromSpan(span *model.Span, topic, source string, file TagAppendRules, rule KindRewriteRules) (*sls.LogGroup, error) {
	logs, err := c.fromSpanToLogs(span, file, rule)
	if err != nil {
		return nil, err
	}
	return &sls.LogGroup{
		Topic:  proto.String(topic),
		Source: proto.String(source),
		Logs:   logs,
	}, nil
}

func (c converter) fromSpanToLogs(span *model.Span, file TagAppendRules, rule KindRewriteRules) ([]*sls.Log, error) {
	contents, err := c.fromSpanToLogContents(span, file, rule)
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

func TraceIDToString(t *model.TraceID) string {
	return fmt.Sprintf("%016x%016x", t.High, t.Low)
}

func (c converter) fromSpanToLogContents(span *model.Span, tagAppendRules TagAppendRules, kindRewriteRules KindRewriteRules) ([]*sls.LogContent, error) {
	contents := make([]*sls.LogContent, 0)
	contents = c.appendContents(contents, traceIDField, TraceIDToString(&span.TraceID))
	contents = c.appendContents(contents, spanIDField, span.SpanID.String())
	contents = c.appendContents(contents, parentSpanIDField, span.ParentSpanID.String())
	contents = c.appendContents(contents, operationNameField, span.OperationName)
	contents = c.appendContents(contents, flagsField, fmt.Sprintf("%d", span.Flags))
	contents = c.appendContents(contents, startTimeField, cast.ToString(span.StartTime.UnixNano()/1000))
	contents = c.appendContents(contents, durationField, cast.ToString(span.Duration.Nanoseconds()/1000))
	contents = c.appendContents(contents, endTimeField, cast.ToString((span.StartTime.UnixNano()+span.Duration.Nanoseconds())/1000))
	contents = c.appendContents(contents, serviceNameField, span.Process.ServiceName)
	contents = c.appendContents(contents, "statusCode", "UNSET")

	attributeMap := make(map[string]string)
	for _, tag := range span.Tags {
		if k, ok := tagAppendRules.SpanTagRules()[tag.Key]; ok {
			attributeMap[k.TagKey] = k.TagValue
		}
		attributeMap[tag.Key] = tag.AsString()

		if k, ok := kindRewriteRules.SpanKindRules()[tag.Key]; ok {
			c.appendContents(contents, spanKindField, k)
		}
	}

	for key, value := range tagAppendRules.OperationPrefixRules() {
		if strings.HasPrefix(span.OperationName, key) {
			attributeMap[value.TagKey] = value.TagValue
		}
	}

	for key, value := range kindRewriteRules.OperationPrefixRules() {
		if strings.HasPrefix(span.OperationName, key) {
			c.appendContents(contents, spanKindField, value)
		}
	}

	tagStr, _ := json.Marshal(attributeMap)
	contents = c.appendContents(contents, "attribute", string(tagStr))

	resourcesMap := make(map[string]string)
	for _, tag := range span.Process.Tags {
		resourcesMap[tag.Key] = tag.AsString()
	}
	resStr, _ := json.Marshal(resourcesMap)
	contents = c.appendContents(contents, "resource", string(resStr))

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

type SLSTraceLogs struct {
	Attribute map[string]string `json:"attribute"`
	Time      int64             `json:"time"` //nano
}

func fieldsToAttribute(fields []model.KeyValue) map[string]string {
	m := make(map[string]string)
	for _, keyVal := range fields {
		m[keyVal.Key] = keyVal.AsString()
	}
	return m
}

func (c converter) appendLogs(contents []*sls.LogContent, logs []model.Log) ([]*sls.LogContent, error) {
	if len(logs) < 1 {
		return contents, nil
	}

	slsLogs := make([]SLSTraceLogs, len(logs))
	for i, l := range logs {
		slsLogs[i] = SLSTraceLogs{
			Time:      l.Timestamp.UnixNano(),
			Attribute: fieldsToAttribute(l.Fields),
		}
	}

	r, err := json.Marshal(slsLogs)
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
			if len(v) == 0 {
				continue
			}
			traceID, err := model.TraceIDFromString(v)
			if err != nil {
				return nil, err
			}
			span.TraceID = traceID
		case spanIDField:
			if len(v) == 0 {
				continue
			}
			spanID, err := model.SpanIDFromString(v)
			if err != nil {
				return nil, err
			}
			span.SpanID = spanID
		case parentSpanIDField:
			if len(v) == 0 {
				continue
			}
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
			span.StartTime = model.EpochMicrosecondsAsTime(cast.ToUint64(v))
		case durationField:
			span.Duration = model.MicrosecondsAsDuration(cast.ToUint64(v))
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
		case "statusMessage":
			if v != "" {
				span.Warnings = append(span.Warnings, v)
			}
		case "attribute":
			tags = c.convertTags(v)
		case "resource":
			process.Tags = c.convertTags(v)
		case "statusCode":
			if v == "ERROR" {
				span.Warnings = append(span.Warnings, v)
			}
		}
	}
	if len(span.References) == 0 && span.ParentSpanID != 0 {
		span.References = append(span.References, model.SpanRef{
			RefType: model.ChildOf,
			TraceID: span.TraceID,
			SpanID:  span.ParentSpanID,
		})
	}
	span.Tags = tags
	span.Process = &process
	return &span, nil
}

func (c converter) convertTags(v string) (kvs model.KeyValues) {
	maps := make(map[string]interface{})
	err := json.Unmarshal([]byte(v), &maps)
	if err != nil {
		kvs = append(kvs, model.String("tags", v))
		return
	}
	for k, v := range maps {
		kvs = append(kvs, model.String(k, fmt.Sprint(v)))
	}
	return kvs
}

func (c converter) unmarshalReferences(s string) (refs []model.SpanRef, err error) {
	err = json.Unmarshal([]byte(s), &refs)
	if err != nil {
		return nil, err
	}
	return refs, nil
}

type otelLog struct {
	Name      string                 `json:"name"`
	Attribute map[string]interface{} `json:"attribute"`
	Time      int64                  `json:"time"`
}

func (c converter) unmarshalLogs(s string) (logs []model.Log, err error) {
	if s == "" {
		return nil, nil
	}
	otelLogs := make([]otelLog, 0)
	err = json.Unmarshal([]byte(s), &otelLogs)
	if err != nil {
		var jaegerLog model.Log
		jaegerLog.Fields = append(jaegerLog.Fields, model.String("logs", s))
		logs = append(logs, jaegerLog)
		return logs, nil
	}
	for _, l := range otelLogs {
		var jaegerLog model.Log
		for k, v := range l.Attribute {
			jaegerLog.Fields = append(jaegerLog.Fields, model.String(k, fmt.Sprint(v)))
		}
		if l.Name != "" {
			jaegerLog.Fields = append(jaegerLog.Fields, model.String("name", l.Name))
		}
		jaegerLog.Timestamp = time.Unix(l.Time/1e9, l.Time%1e9)
		logs = append(logs, jaegerLog)
	}
	return logs, nil
}

func (c converter) unmarshalWarnings(s string) (warnings []string, err error) {
	if s == "" {
		return nil, nil
	}
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
