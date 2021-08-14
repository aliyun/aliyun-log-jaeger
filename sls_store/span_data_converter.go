package sls_store

import (
	"encoding/json"
	"fmt"
	"time"

	slsSdk "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/proto"
	"github.com/jaegertracing/jaeger/model"
	"github.com/spf13/cast"
)

type DataConverter interface {
	ToJaegerSpan(data map[string]string) (*model.Span, error)

	ToSLSSpan(span *model.Span) ([]*slsSdk.LogContent, error)
}

var dataConvert = &dataConverterImpl{}

type dataConverterImpl struct {
}

func (dataConverterImpl) ToJaegerSpan(log map[string]string) (*model.Span, error) {
	span := model.Span{}
	process := model.Process{
		Tags: make([]model.KeyValue, 0),
	}

	for k, v := range log {
		switch k {
		case TraceID:
			traceID, err := model.TraceIDFromString(v)
			if err != nil {
				logger.Warn("Failed to convert traceId", "key", k, "value", v)
				return nil, err
			}
			span.TraceID = traceID
			break
		case SpanID:
			spanID, err := model.SpanIDFromString(v)
			if err != nil {
				logger.Warn("Failed to convert spanID", "key", k, "value", v)
				return nil, err
			}
			span.SpanID = spanID
			break
		case OperationName:
			span.OperationName = v
			break
		case Flags:
			span.Flags = model.Flags(cast.ToUint64(v))
			break
		case StartTime:
			span.StartTime = model.EpochMicrosecondsAsTime(cast.ToUint64(v))
			break
		case Duration:
			span.Duration = model.MicrosecondsAsDuration(cast.ToUint64(v))
			break
		case ServiceName:
			process.ServiceName = v
			break
		case Links:
			refs, err := unmarshalReferences(v)
			if err != nil {
				logger.Warn("Failed to convert links", "key", k, "value", v, "exception", err)
				return nil, err
			}
			span.References = refs
			break
		case Logs:
			logs, err := unmarshalLogs(v)
			if err != nil {
				logger.Warn("Failed to convert logs", "key", k, "value", v, "exception", err)
				return nil, err
			}
			span.Logs = logs
			break
		case StatusMessageField:
			if v != "" {
				span.Warnings = append(span.Warnings, v)
			}
			break
		case Attribute:
			span.Tags = unmarshalTags(v)
			break
		case Resource:
			process.Tags, span.ProcessID = unmarshalResource(v)
			break
		case StatusCodeField:
			if v == "ERROR" {
				span.Warnings = append(span.Warnings, v)
			}
		}
	}

	span.Process = &process
	return &span, nil
}

func (dataConverterImpl) ToSLSSpan(span *model.Span) ([]*slsSdk.LogContent, error) {
	contents := make([]*slsSdk.LogContent, 0)
	contents = appendAttributeToLogContent(contents, TraceID, TraceIDToString(&span.TraceID))
	contents = appendAttributeToLogContent(contents, SpanID, span.SpanID.String())
	contents = appendAttributeToLogContent(contents, ParentSpanID, span.ParentSpanID().String())
	contents = appendAttributeToLogContent(contents, OperationName, span.OperationName)
	contents = appendAttributeToLogContent(contents, Flags, fmt.Sprintf("%d", span.Flags))
	contents = appendAttributeToLogContent(contents, StartTime, cast.ToString(span.StartTime.UnixNano()/1000))
	contents = appendAttributeToLogContent(contents, Duration, cast.ToString(span.Duration.Nanoseconds()/1000))
	contents = appendAttributeToLogContent(contents, EndTime, cast.ToString((span.StartTime.UnixNano()+span.Duration.Nanoseconds())/1000))
	contents = appendAttributeToLogContent(contents, ServiceName, span.Process.ServiceName)
	contents = appendAttributeToLogContent(contents, StatusCode, "UNSET")
	contents = appendAttributeToLogContent(contents, Attribute, marshalTags(span.Tags))
	contents = appendAttributeToLogContent(contents, Resource, marshalResource(span.Process.Tags, span.ProcessID))

	if refStr, err := marshalReferences(span.References); err != nil {
		logger.Warn("Failed to convert references", "spanID", span.SpanID, "reference", span.References, "exception", err)
		return nil, err
	} else {
		contents = appendAttributeToLogContent(contents, Links, refStr)
	}

	if logsStr, err := marshalLogs(span.Logs); err != nil {
		logger.Warn("Failed to convert logs", "spanID", span.SpanID, "logs", span.Logs, "exception", err)
		return nil, err
	} else {
		contents = appendAttributeToLogContent(contents, Logs, logsStr)
	}

	contents, err := appendWarnings(contents, span.Warnings)
	if err != nil {
		logger.Warn("Failed to convert warnings", "spanID", span.SpanID, "warnings", span.Warnings, "exception", err)
		return nil, err
	}

	return contents, nil
}

func appendWarnings(contents []*slsSdk.LogContent, warnings []string) ([]*slsSdk.LogContent, error) {
	if len(warnings) < 1 {
		return contents, nil
	}

	r, err := json.Marshal(warnings)
	if err != nil {
		return nil, err
	}

	return appendAttributeToLogContent(contents, StatusMessage, string(r)), nil
}

func marshalResource(v []model.KeyValue, processID string) string {
	dataMap := keyValueToMap(v)
	dataMap["ProcessID"] = processID

	data, err := json.Marshal(dataMap)
	if err != nil {
		return fmt.Sprintf("%v", string(data))
	}

	return string(data)
}

func unmarshalResource(v string) (kvs []model.KeyValue, processID string) {
	data := make(map[string]string)

	err := json.Unmarshal([]byte(v), &data)
	if err != nil {
		kvs = append(kvs, model.String("tags", v))
		return kvs, ""
	}

	return mapToKeyValue(data), data["ProcessID"]

}

func marshalTags(v []model.KeyValue) string {
	dataMap := keyValueToMap(v)

	data, err := json.Marshal(dataMap)
	if err != nil {
		return fmt.Sprintf("%v", string(data))
	}

	return string(data)
}

func unmarshalTags(v string) (kvs []model.KeyValue) {
	data := make(map[string]string)

	err := json.Unmarshal([]byte(v), &data)
	if err != nil {
		kvs = append(kvs, model.String("tags", v))
		return
	}

	return mapToKeyValue(data)
}

type SpanLog struct {
	Attribute map[string]string `json:"attribute"`
	Time      int64             `json:"time"`
}

func marshalLogs(logs []model.Log) (string, error) {
	if len(logs) <= 0 {
		return "[]", nil
	}

	slsLogs := make([]SpanLog, len(logs))
	for i, l := range logs {
		slsLogs[i] = SpanLog{
			Time:      l.Timestamp.UnixNano(),
			Attribute: keyValueToMap(l.Fields),
		}
	}

	r, err := json.Marshal(slsLogs)
	if err != nil {
		return "", err
	}

	return string(r), nil
}

func unmarshalLogs(s string) ([]model.Log, error) {
	if s == "[]" {
		return nil, nil
	}

	logs := make([]SpanLog, 0)
	if err := json.Unmarshal([]byte(s), &logs); err != nil {
		return nil, err
	}

	result := make([]model.Log, len(logs))
	for i, log := range logs {
		result[i] = model.Log{
			Timestamp: time.Unix(log.Time/1e9, log.Time%1e9),
			Fields:    mapToKeyValue(log.Attribute),
		}
	}
	return result, nil
}

func marshalReferences(refs []model.SpanRef) (string, error) {
	if len(refs) <= 0 {
		return "[]", nil
	}

	rs := make([]map[string]string, 0)

	for _, ref := range refs {
		r := make(map[string]string)
		r["TraceID"] = ref.TraceID.String()
		r["SpanID"] = ref.SpanID.String()
		r["RefType"] = ref.RefType.String()
		rs = append(rs, r)
	}

	r, err := json.Marshal(rs)
	if err != nil {
		return "", err
	}

	return string(r), nil
}

func unmarshalReferences(s string) (refs []model.SpanRef, err error) {
	if s == "[]" {
		return nil, nil
	}

	rs := make([]map[string]string, 0)

	err = json.Unmarshal([]byte(s), &rs)
	if err != nil {
		return nil, err
	}

	for _, r := range rs {
		tid, e1 := model.TraceIDFromString(r["TraceID"])
		if e1 != nil {
			return nil, e1
		}

		spanID, e2 := model.SpanIDFromString(r["SpanID"])
		if e2 != nil {
			return nil, e2
		}

		spanType := model.SpanRefType_value[r["RefType"]]
		refs = append(refs, model.SpanRef{
			TraceID: tid,
			SpanID:  spanID,
			RefType: model.SpanRefType(spanType),
		})
	}

	return refs, nil
}

func mapToKeyValue(data map[string]string) []model.KeyValue {
	result := make([]model.KeyValue, 0)
	for key, value := range data {
		result = append(result, model.String(key, value))
	}

	return result
}

func keyValueToMap(fields []model.KeyValue) map[string]string {
	m := make(map[string]string)
	for _, keyVal := range fields {
		m[keyVal.Key] = keyVal.AsString()
	}
	return m
}

func TraceIDToString(t *model.TraceID) string {
	return t.String()
}

func appendAttributeToLogContent(contents []*slsSdk.LogContent, k, v string) []*slsSdk.LogContent {
	content := slsSdk.LogContent{
		Key:   proto.String(k),
		Value: proto.String(v),
	}
	return append(contents, &content)
}
