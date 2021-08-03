package sls_store

import (
	"context"
	slsSdk "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"time"
)

type slsSpanReader struct {
	client      *slsSdk.Client
	instance    slsTraceInstance
	maxLookBack time.Duration
}

func (s slsSpanReader) GetTrace(ctx context.Context, traceID model.TraceID) (*model.Trace, error) {
	// traceID: xxx
	return nil, nil
}

func (s slsSpanReader) GetServices(ctx context.Context) ([]string, error) {
	// * | select DISTINCT service as service from log
	return nil, nil
}

func (s slsSpanReader) GetOperations(ctx context.Context, query spanstore.OperationQueryParameters) ([]spanstore.Operation, error) {
	// * and kind: client and service: front-end
	return nil, nil
}

func (s slsSpanReader) FindTraces(ctx context.Context, query *spanstore.TraceQueryParameters) ([]*model.Trace, error) {
	//
	return nil, nil
}

func (s slsSpanReader) FindTraceIDs(ctx context.Context, query *spanstore.TraceQueryParameters) ([]model.TraceID, error) {
	return nil, nil
}
