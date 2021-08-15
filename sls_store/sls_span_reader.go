package sls_store

import (
	"context"
	"time"

	slsSdk "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/go-hclog"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

type slsSpanReader struct {
	client      *slsSdk.Client
	instance    slsTraceInstance
	maxLookBack time.Duration
	logger      hclog.Logger
}

func (s slsSpanReader) GetServices(ctx context.Context) ([]string, error) {
	defer func() {
		if err := recover(); err != nil {
			s.logger.Error("Failed to GetServices", "Exception", err)
		}
	}()
	from, to := buildSearchingData(s.maxLookBack)

	response, e := s.client.GetLogs(s.instance.project(), s.instance.traceLogStore(), DefaultTopicName, from, to,
		toGetServicesQuery(), DefaultFetchNumber, DefaultOffset, false)

	s.logger.Info("GetServicesList", "Query", toGetServicesQuery(), "StartTime", time.Unix(from, 0), "EndTime", time.Unix(to, 0), "Logstore", s.instance.traceLogStore())

	if e != nil {
		return nil, e
	}

	services := make([]string, response.Count)
	for i, data := range response.Logs {
		services[i] = data[ServiceName]
	}

	return services, nil

}

func (s slsSpanReader) GetOperations(ctx context.Context, query spanstore.OperationQueryParameters) ([]spanstore.Operation, error) {
	defer func() {
		if err := recover(); err != nil {
			s.logger.Error("Failed to get operations", "Exception", err)
		}
	}()

	from, to := buildSearchingData(s.maxLookBack)

	response, e := s.client.GetLogs(s.instance.project(), s.instance.traceLogStore(), DefaultTopicName, from, to,
		toOperationsQuery(query), DefaultFetchNumber, DefaultOffset, false)

	s.logger.Info("GetOperations", "Query", toOperationsQuery(query), "StartTime", time.Unix(from, 0), "EndTime", time.Unix(to, 0), "Logstore", s.instance.traceLogStore())
	if e != nil {
		return nil, e
	}

	operations := make([]spanstore.Operation, response.Count)
	for i, data := range response.Logs {
		operations[i] = spanstore.Operation{
			Name:     data[OperationName],
			SpanKind: data[SpanKind],
		}
	}

	return operations, nil
}

func (s slsSpanReader) FindTraces(ctx context.Context, query *spanstore.TraceQueryParameters) ([]*model.Trace, error) {
	defer func() {
		if err := recover(); err != nil {
			s.logger.Error("Failed to find traces", "Exceptions", err)
		}
	}()

	traceIDs, err := GetTraceIDsWithQuery(s.client, s.instance.project(), s.instance.traceLogStore(), query)
	if err != nil {
		return nil, err
	}

	var result []*model.Trace
	for _, tid := range traceIDs {
		if t, e := GetTraceWithTime(s.client, tid, query.StartTimeMin.Unix(), query.StartTimeMax.Unix(), s.instance.project(),
			s.instance.traceLogStore()); e == nil {
			result = append(result, t)
		} else {
			logger.Warn("Failed to get trace data.", "TID", tid, "Exception", e)
		}
	}

	return result, nil
}

func (s slsSpanReader) FindTraceIDs(ctx context.Context, query *spanstore.TraceQueryParameters) ([]model.TraceID, error) {
	defer func() {
		if err := recover(); err != nil {
			s.logger.Error("Failed to FindTraceIDs", "Exception", err)
		}
	}()

	return GetTraceIDsWithQuery(s.client, s.instance.project(), s.instance.traceLogStore(), query)
}

func (s slsSpanReader) GetTrace(ctx context.Context, traceID model.TraceID) (*model.Trace, error) {
	from, to := buildSearchingData(s.maxLookBack)
	return GetTraceWithTime(s.client, traceID, from, to, s.instance.project(), s.instance.traceLogStore())
}

var logger = hclog.New(&hclog.LoggerOptions{
	Level:      hclog.Info,
	Name:       "aliyun-log-jaeger-plugin",
	JSONFormat: true,
})

func GetTraceIDsWithQuery(client *slsSdk.Client, project, logstore string, query *spanstore.TraceQueryParameters) ([]model.TraceID, error) {
	from, to := query.StartTimeMin.Unix(), query.StartTimeMax.Unix()
	queryString := toFindTraceIdsQuery(query)

	response, e := client.GetLogs(project, logstore, DefaultTopicName, from, to,
		queryString, DefaultFetchNumber, DefaultOffset, false)

	if e != nil {
		return nil, e
	}

	traceIDS := make(map[string]bool)
	for _, log := range response.Logs {
		traceIDS[log[TraceIDField]] = true
	}

	var result []model.TraceID

	for key, _ := range traceIDS {
		traceId, e1 := model.TraceIDFromString(key)
		if e1 != nil {
			logger.Warn("Failed to convert trace ID", "tid", key)
			continue
		}

		result = append(result, traceId)
	}

	return result, nil
}

func GetTraceWithTime(client *slsSdk.Client, traceID model.TraceID, from, to int64, project, logstore string) (*model.Trace, error) {
	if response, e := client.GetLogs(project, logstore, DefaultTopicName, from, to, toGetTraceQuery(traceID),
		DefaultFetchNumber, DefaultOffset, false); e != nil {
		return nil, e
	} else {
		return mappingTraceData(response.Logs)
	}
}

// mappingTraceData the method used to converting sls span data to jaeger span data.
func mappingTraceData(logs []map[string]string) (*model.Trace, error) {
	var processMapping []model.Trace_ProcessMapping
	spans := make([]*model.Span, 0)
	for _, data := range logs {
		if spanData, err := dataConvert.ToJaegerSpan(data); err != nil {
			continue
		} else {
			spans = append(spans, spanData)
		}
	}

	for _, span := range spans {
		processMapping = append(processMapping, model.Trace_ProcessMapping{
			ProcessID: span.ProcessID,
			Process:   *span.Process,
		})
	}

	return &model.Trace{
		Spans:      spans,
		ProcessMap: processMapping,
	}, nil
}
