// Copyright (c) 2018 The Jaeger Authors.
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
	"strconv"
	"testing"
	"time"

	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/jaegertracing/jaeger/storage/spanstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestSpanReader_logsToStringArray(t *testing.T) {
	logs := make([]map[string]string, 3)
	logs[0] = map[string]string{
		traceIDField:       "0",
		operationNameField: "op0",
	}
	logs[1] = map[string]string{
		traceIDField:       "1",
		operationNameField: "op1",
	}
	logs[2] = map[string]string{
		traceIDField:       "2",
		operationNameField: "op2",
	}
	actual, err := logsToStringArray(logs, operationNameField)
	require.NoError(t, err)
	assert.EqualValues(t, []string{"op0", "op1", "op2"}, actual)
}

func TestSpanReader_buildFindTracesQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := `process.serviceName: "s" and operationName: "o" and duration >= 1000000000 and duration <= 2000000000 and tags.http.status_code: "200" | select traceID, max_by("process.serviceName", duration) as "process.serviceName", max_by(operationName, duration) as operationName, max_by(duration, duration) as duration, count(1) as spansCount from (select traceID, "process.serviceName", operationName, duration from log limit 10000) group by traceID limit 30`
	traceQuery := &spanstore.TraceQueryParameters{
		DurationMin:   time.Second,
		DurationMax:   time.Second * 2,
		StartTimeMin:  time.Time{},
		StartTimeMax:  time.Time{}.Add(time.Second),
		ServiceName:   "s",
		OperationName: "o",
		Tags: map[string]string{
			"http.status_code": "200",
		},
		NumTraces: 30,
	}
	actualQuery := r.buildFindTracesQuery(traceQuery)
	assert.Equal(t, expectedStr, actualQuery)
}

func TestSpanReader_buildFindTracesQuery_emptyQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := `| select traceID, max_by("process.serviceName", duration) as "process.serviceName", max_by(operationName, duration) as operationName, max_by(duration, duration) as duration, count(1) as spansCount from (select traceID, "process.serviceName", operationName, duration from log limit 10000) group by traceID limit ` + strconv.Itoa(defaultNumTraces)
	traceQuery := &spanstore.TraceQueryParameters{
		NumTraces: defaultNumTraces,
	}
	actualQuery := r.buildFindTracesQuery(traceQuery)
	assert.Equal(t, expectedStr, actualQuery)
}

func TestSpanReader_buildFindTracesQuery_singleCondition(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := `process.serviceName: "svc1" | select traceID, max_by("process.serviceName", duration) as "process.serviceName", max_by(operationName, duration) as operationName, max_by(duration, duration) as duration, count(1) as spansCount from (select traceID, "process.serviceName", operationName, duration from log limit 10000) group by traceID limit ` + strconv.Itoa(defaultNumTraces)
	traceQuery := &spanstore.TraceQueryParameters{
		ServiceName: "svc1",
		NumTraces: defaultNumTraces,
	}
	actualQuery := r.buildFindTracesQuery(traceQuery)
	assert.Equal(t, expectedStr, actualQuery)
}

func TestSpanReader_buildServiceNameQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := `process.serviceName: "svc1"`
	serviceNameQuery := r.buildServiceNameQuery("svc1")
	assert.Equal(t, expectedStr, serviceNameQuery)
}

func TestSpanReader_buildOperationNameQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := `operationName: "HTTP GET"`
	operationNameQuery := r.buildOperationNameQuery("HTTP GET")
	assert.Equal(t, expectedStr, operationNameQuery)
}

func TestSpanReader_buildDurationQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := "duration >= 1000000000 and duration <= 2000000000"
	durationMin := time.Second
	durationMax := time.Second * 2
	durationQuery := r.buildDurationQuery(durationMin, durationMax)
	assert.Equal(t, expectedStr, durationQuery)

	expectedStr = "duration >= 12000000"
	durationMin = time.Millisecond * 12
	durationQuery = r.buildDurationQuery(durationMin, 0)
	assert.Equal(t, expectedStr, durationQuery)

	expectedStr = "duration <= 18000000000000"
	durationMax = time.Hour * 5
	durationQuery = r.buildDurationQuery(0, durationMax)
	assert.Equal(t, expectedStr, durationQuery)

	expectedStr = ""
	durationQuery = r.buildDurationQuery(0, 0)
	assert.Equal(t, expectedStr, durationQuery)
}

func TestSpanReader_buildTagQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := `tags.http.method: "POST"`
	operationNameQuery := r.buildTagQuery("http.method", "POST")
	assert.Equal(t, expectedStr, operationNameQuery)
}

func TestSpanReader_combineSubQueries(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := ""
	query := r.combineSubQueries(nil)
	assert.Equal(t, expectedStr, query)

	var subQueries []string
	expectedStr = ""
	query = r.combineSubQueries(subQueries)
	assert.Equal(t, expectedStr, query)

	subQueries = []string{"", "", ""}
	expectedStr = ""
	query = r.combineSubQueries(subQueries)
	assert.Equal(t, expectedStr, query)

	subQueries = []string{`tags.http.method: "POST"`}
	expectedStr = `tags.http.method: "POST"`
	query = r.combineSubQueries(subQueries)
	assert.Equal(t, expectedStr, query)

	subQueries = []string{`process.serviceName: "svc1"`, `tags.http.method: "POST"`}
	expectedStr = `process.serviceName: "svc1" and tags.http.method: "POST"`
	query = r.combineSubQueries(subQueries)
	assert.Equal(t, expectedStr, query)

	subQueries = []string{`process.serviceName: "svc1"`, `operationName: "HTTP GET"`, `tags.http.method: "POST"`}
	expectedStr = `process.serviceName: "svc1" and operationName: "HTTP GET" and tags.http.method: "POST"`
	query = r.combineSubQueries(subQueries)
	assert.Equal(t, expectedStr, query)
}

func TestSpanReader_getQuerySuffix(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := `| select traceID, max_by("process.serviceName", duration) as "process.serviceName", max_by(operationName, duration) as operationName, max_by(duration, duration) as duration, count(1) as spansCount from (select traceID, "process.serviceName", operationName, duration from log limit 10000) group by traceID limit 20`
	query := r.getQuerySuffix(10000, 20)
	assert.Equal(t, expectedStr, query)
}
