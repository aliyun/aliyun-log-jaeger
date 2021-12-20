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
	"fmt"
	"testing"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/proto"
	"github.com/jaegertracing/jaeger/model"
	"github.com/kr/pretty"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

var (
	someTraceID       = model.TraceID{High: 22222, Low: 44444}
	someSpanID        = model.SpanID(3333)
	someParentSpanID  = model.SpanID(11111)
	someOperationName = "someOperationName"

	someRefs = []model.SpanRef{
		{
			TraceID: someTraceID,
			SpanID:  someParentSpanID,
			RefType: model.ChildOf,
		},
	}
	someRefsValueStr = "[{\"refType\":\"child-of\",\"traceID\":\"56ce000000000000ad9c\",\"spanID\":\"2b67\"}]"

	someStartTime    = model.EpochMicrosecondsAsTime(55555)
	someDuration     = model.MicrosecondsAsDuration(50000)
	someFlags        = model.Flags(1)
	someLogTimestamp = model.EpochMicrosecondsAsTime(12345)
	someServiceName  = "someServiceName"

	someStringTagKey   = "someStringTag"
	someStringTagValue = "someTagValue"

	someBoolTagKey      = "someBoolTag"
	someBoolTagValue    = true
	someBoolTagValueStr = "true"

	someLongTagKey      = "someLongTag"
	someLongTagValue    = int64(123)
	someLongTagValueStr = "123"

	someDoubleTagKey      = "someDoubleTag"
	someDoubleTagValue    = float64(1.4)
	someDoubleTagValueStr = "1.4"

	someBinaryTagKey      = "someBinaryTag"
	someBinaryTagValue    = []byte("someBinaryValue")
	someBinaryTagValueStr = "736f6d6542696e61727956616c7565"

	someTags = model.KeyValues{
		model.String(someStringTagKey, someStringTagValue),
		model.String("db.instance", "db instance"),
		model.Bool(someBoolTagKey, someBoolTagValue),
		model.Int64(someLongTagKey, someLongTagValue),
		model.Float64(someDoubleTagKey, someDoubleTagValue),
		model.Binary(someBinaryTagKey, someBinaryTagValue),
	}
	someTagsValueStr = "[{\"key\":\"someStringTag\",\"vType\":\"string\",\"vStr\":\"someTagValue\"},{\"key\":\"someBoolTag\",\"vType\":\"bool\",\"vNum\":1},{\"key\":\"someLongTag\",\"vType\":\"int64\",\"vNum\":123},{\"key\":\"someDoubleTag\",\"vType\":\"float64\",\"vNum\":4608983858650965606},{\"key\":\"someBinaryTag\",\"vType\":\"binary\",\"vBlob\":\"c29tZUJpbmFyeVZhbHVl\"}]"

	someUnusualTags = model.KeyValues{
		model.String(someStringTagKey, someStringTagValue),
		model.String(someBoolTagKey, someBoolTagValueStr),
		model.String(someLongTagKey, someLongTagValueStr),
		model.String(someDoubleTagKey, someDoubleTagValueStr),
		model.String(someBinaryTagKey, someBinaryTagValueStr),
	}

	someLogs = []model.Log{
		{
			Timestamp: someLogTimestamp,
			Fields:    someTags,
		},
	}
	someLogsValueStr = "[{\"timestamp\":\"1970-01-01T08:00:00.012345+08:00\",\"fields\":[{\"key\":\"someStringTag\",\"vType\":\"string\",\"vStr\":\"someTagValue\"},{\"key\":\"someBoolTag\",\"vType\":\"bool\",\"vNum\":1},{\"key\":\"someLongTag\",\"vType\":\"int64\",\"vNum\":123},{\"key\":\"someDoubleTag\",\"vType\":\"float64\",\"vNum\":4608983858650965606},{\"key\":\"someBinaryTag\",\"vType\":\"binary\",\"vBlob\":\"c29tZUJpbmFyeVZhbHVl\"}]}]"

	someWarnings         = []string{"warning1", "warning2", "warning3"}
	someWarningsValueStr = "[\"warning1\",\"warning2\",\"warning3\"]"
)

func getTestJaegerSpan() *model.Span {
	return &model.Span{
		TraceID:       someTraceID,
		SpanID:        someSpanID,
		ParentSpanID:  someParentSpanID,
		OperationName: someOperationName,
		References:    someRefs,
		Flags:         someFlags,
		StartTime:     someStartTime,
		Duration:      someDuration,
		Tags:          someTags,
		Logs:          someLogs,
		Process:       getTestJaegerProcess(),
		Warnings:      someWarnings,
	}
}

func getTestJaegerProcess() *model.Process {
	return &model.Process{
		ServiceName: someServiceName,
		Tags:        someTags,
	}
}

func getTestUnusualJaegerSpan() *model.Span {
	return &model.Span{
		TraceID:       someTraceID,
		SpanID:        someSpanID,
		ParentSpanID:  someParentSpanID,
		OperationName: someOperationName,
		References:    someRefs,
		Flags:         someFlags,
		StartTime:     someStartTime,
		Duration:      someDuration,
		Tags:          someUnusualTags,
		Logs:          someLogs,
		Process:       getTestUnusualJaegerProcess(),
		Warnings:      someWarnings,
	}
}

func getTestUnusualJaegerProcess() *model.Process {
	return &model.Process{
		ServiceName: someServiceName,
		Tags:        someUnusualTags,
	}
}

func getTestLog() map[string]string {
	return map[string]string{
		traceIDField:                         someTraceID.String(),
		spanIDField:                          someSpanID.String(),
		parentSpanIDField:                    someParentSpanID.String(),
		operationNameField:                   someOperationName,
		referenceField:                       someRefsValueStr,
		flagsField:                           fmt.Sprintf("%d", someFlags),
		startTimeField:                       cast.ToString(someStartTime.UnixNano()),
		durationField:                        cast.ToString(someDuration.Nanoseconds()),
		tagsPrefix + someStringTagKey:        someStringTagValue,
		tagsPrefix + someBoolTagKey:          someBoolTagValueStr,
		tagsPrefix + someLongTagKey:          someLongTagValueStr,
		tagsPrefix + someDoubleTagKey:        someDoubleTagValueStr,
		tagsPrefix + someBinaryTagKey:        someBinaryTagValueStr,
		logsField:                            someLogsValueStr,
		serviceNameField:                     someServiceName,
		processTagsPrefix + someStringTagKey: someStringTagValue,
		processTagsPrefix + someBoolTagKey:   someBoolTagValueStr,
		processTagsPrefix + someLongTagKey:   someLongTagValueStr,
		processTagsPrefix + someDoubleTagKey: someDoubleTagValueStr,
		processTagsPrefix + someBinaryTagKey: someBinaryTagValueStr,
		warningsField:                        someWarningsValueStr,
	}
}

func convertLogGroupToMap(logGroup *sls.LogGroup) map[string]string {
	m := make(map[string]string)
	for _, content := range logGroup.Logs[0].Contents {
		m[*content.Key] = *content.Value
	}
	return m
}

func TestToSpan(t *testing.T) {
	expectedSpan := getTestUnusualJaegerSpan()
	expectedSpan.Tags.Sort()
	expectedSpan.Process.Tags.Sort()
	actualSpan, err := ToSpan(getTestLog())
	actualSpan.Tags.Sort()
	actualSpan.Process.Tags.Sort()
	assert.NoError(t, err)
	if !assert.EqualValues(t, expectedSpan, actualSpan) {
		for _, diff := range pretty.Diff(expectedSpan, actualSpan) {
			t.Log(diff)
		}
	}
}

func TestFromSpan(t *testing.T) {
	span := getTestJaegerSpan()
	logGroup, err := FromSpan(span, "topic", "0.0.0.0", initTagAppendRules("", false), initKindRewriteRules("", false))
	assert.Nil(t, err)
	assert.Equal(t, "topic", *logGroup.Topic)
	assert.Equal(t, "0.0.0.0", *logGroup.Source)
	expectedLog := getTestLog()
	actualLog := convertLogGroupToMap(logGroup)
	if !assert.EqualValues(t, expectedLog, actualLog) {
		for _, diff := range pretty.Diff(expectedLog, actualLog) {
			t.Log(diff)
		}
	}
}

func TestAppendReferences(t *testing.T) {
	contents := make([]*sls.LogContent, 0)
	contents, err := converter{}.appendReferences(contents, someRefs)
	assert.NoError(t, err)
	expectedContents := make([]*sls.LogContent, 0)
	content := sls.LogContent{
		Key:   proto.String(referenceField),
		Value: proto.String(someRefsValueStr),
	}
	expectedContents = append(expectedContents, &content)
	assert.Equal(t, expectedContents, contents)
}

func TestAppendReferences_nil(t *testing.T) {
	contents := make([]*sls.LogContent, 0)
	contents, err := converter{}.appendReferences(contents, nil)
	assert.NoError(t, err)
	expectedContents := make([]*sls.LogContent, 0)
	assert.Equal(t, expectedContents, contents)
}

func TestAppendLogs(t *testing.T) {
	contents := make([]*sls.LogContent, 0)
	contents, err := converter{}.appendLogs(contents, someLogs)
	assert.NoError(t, err)
	expectedContents := make([]*sls.LogContent, 0)
	content := sls.LogContent{
		Key:   proto.String(logsField),
		Value: proto.String(someLogsValueStr),
	}
	expectedContents = append(expectedContents, &content)
	assert.Equal(t, expectedContents, contents)
}

func TestAppendLogs_nil(t *testing.T) {
	contents := make([]*sls.LogContent, 0)
	contents, err := converter{}.appendLogs(contents, nil)
	assert.NoError(t, err)
	expectedContents := make([]*sls.LogContent, 0)
	assert.Equal(t, expectedContents, contents)
}

func TestAppendWarnings(t *testing.T) {
	contents := make([]*sls.LogContent, 0)
	contents, err := converter{}.appendWarnings(contents, someWarnings)
	assert.NoError(t, err)
	expectedContents := make([]*sls.LogContent, 0)
	content := sls.LogContent{
		Key:   proto.String(warningsField),
		Value: proto.String(someWarningsValueStr),
	}
	expectedContents = append(expectedContents, &content)
	assert.Equal(t, expectedContents, contents)
}

func TestAppendWarnings_nil(t *testing.T) {
	contents := make([]*sls.LogContent, 0)
	contents, err := converter{}.appendWarnings(contents, nil)
	assert.NoError(t, err)
	expectedContents := make([]*sls.LogContent, 0)
	assert.Equal(t, expectedContents, contents)
}
