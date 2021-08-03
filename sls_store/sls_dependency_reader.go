package sls_store

import (
	"context"
	slsSdk "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/jaegertracing/jaeger/model"
	"strconv"
	"time"
)

const (
	// DependenciesQueryString The query string which calculates the dependency relationship between each service.
	DependenciesQueryString = "* and version: service_name"

	// ParentServiceFieldName The log item key of parent service name
	ParentServiceFieldName = "parent_service"
	// FailureCallingTimesFieldName The key of failure calling times
	FailureCallingTimesFieldName = "n_status_fail"
	// SuccessfulCallingTimesFieldName The key of successful calling times
	SuccessfulCallingTimesFieldName = "n_status_succ"
	// ParentService the key of parent service
	ParentService = "parent_service"
	// ChildService the key of child service
	ChildService = "child_service"
)

type slsDependencyReader struct {
	client   *slsSdk.Client
	instance slsTraceInstance
}

func (s slsDependencyReader) GetDependencies(ctx context.Context, endTs time.Time, lookback time.Duration) ([]model.DependencyLink, error) {
	response, error := s.client.GetLogs(s.instance.project(), s.instance.serviceDependencyLogStore(), "",
		endTs.Unix()-int64(lookback), endTs.Unix(), DependenciesQueryString, 1000, 0, false)

	if error != nil {
		return nil, error
	}

	result := make([]model.DependencyLink, response.Count)
	for i, log := range response.Logs {
		if log[ParentServiceFieldName] == "None" {
			continue
		}

		result[i] = model.DependencyLink{
			Parent:    log[ParentService],
			Child:     log[ChildService],
			CallCount: getCallingTime(log),
		}
	}

	return result, nil
}

func getCallingTime(log map[string]string) uint64 {
	return getCallingTimes(log, SuccessfulCallingTimesFieldName) + getCallingTimes(log, FailureCallingTimesFieldName)
}

func getCallingTimes(log map[string]string, key string) uint64 {
	successCount, e2 := strconv.Atoi(log[key])
	if e2 != nil {
		successCount = 0
	}
	return uint64(successCount)
}
