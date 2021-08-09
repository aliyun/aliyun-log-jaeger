package sls_store

import (
	"context"
	"strconv"
	"time"

	slsSdk "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/go-hclog"
	"github.com/jaegertracing/jaeger/model"
)

type slsDependencyReader struct {
	client   *slsSdk.Client
	instance slsTraceInstance
	logger   hclog.Logger
}

func (s slsDependencyReader) GetDependencies(ctx context.Context, endTs time.Time, lookback time.Duration) ([]model.DependencyLink, error) {
	defer func() {
		if err := recover(); err != nil {
			s.logger.Error("Failed to get DependencyLink", "Exception", err)
		}
	}()

	response, error := s.client.GetLogs(s.instance.project(), s.instance.serviceDependencyLogStore(), DefaultTopicName,
		endTs.Add(-1*lookback).Unix(), endTs.Unix(), DependenciesQueryString, DefaultFetchNumber, DefaultOffset, false)

	if error != nil {
		return nil, error
	}

	s.logger.Info("GetDependencies", "Query", DependenciesQueryString, "Logstore", s.instance.serviceDependencyLogStore(), "DependencyLinks", response.Count)
	var result []model.DependencyLink
	for _, log := range response.Logs {
		count, _ := strconv.ParseFloat(log["count"], 0)
		if log[ParentService] == "None" {
			continue
		}

		result = append(result, model.DependencyLink{
			Parent:    log[ParentService],
			Child:     log[ChildService],
			CallCount: uint64(count),
		})
	}

	return result, nil
}
