package spanstore

import (
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

func TestSpanReader_buildFindTraceIDsQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := `where "process.serviceName" = 's' and operationName = 'o' and 1000000000 <= duration and duration <= 2000000000 and "tags.http.method" = 'GET' and "tags.http.status_code" = '200'`
	traceQuery := &spanstore.TraceQueryParameters{
		DurationMin:   time.Second,
		DurationMax:   time.Second * 2,
		StartTimeMin:  time.Time{},
		StartTimeMax:  time.Time{}.Add(time.Second),
		ServiceName:   "s",
		OperationName: "o",
		Tags: map[string]string{
			"http.method":      "GET",
			"http.status_code": "200",
		},
	}
	actualQuery := r.buildFindTraceIDsQuery(traceQuery)
	assert.Equal(t, expectedStr, actualQuery)
}

func TestSpanReader_buildServiceNameQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := `"process.serviceName" = 'svc1'`
	serviceNameQuery := r.buildServiceNameQuery("svc1")
	assert.Equal(t, expectedStr, serviceNameQuery)
}

func TestSpanReader_buildOperationNameQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := "operationName = 'op1'"
	operationNameQuery := r.buildOperationNameQuery("op1")
	assert.Equal(t, expectedStr, operationNameQuery)
}

func TestSpanReader_buildDurationQuery(t *testing.T) {
	l := &sls.LogStore{
		Name: "emptyLogStore",
	}
	r := newSpanReader(l, zap.NewNop(), 15*time.Minute)

	expectedStr := "1000000000 <= duration and duration <= 2000000000"
	durationMin := time.Second
	durationMax := time.Second * 2
	durationQuery := r.buildDurationQuery(durationMin, durationMax)
	assert.Equal(t, expectedStr, durationQuery)

	expectedStr = "12000000 <= duration"
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

	expectedStr := `"tags.http.method" = 'POST'`
	operationNameQuery := r.buildTagQuery("http.method", "POST")
	assert.Equal(t, expectedStr, operationNameQuery)
}
