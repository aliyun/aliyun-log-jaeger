package spanstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	//traceQuery := &spanstore.TraceQueryParameters{
	//	DurationMin:   time.Second,
	//	DurationMax:   time.Second * 2,
	//	StartTimeMin:  time.Time{},
	//	StartTimeMax:  time.Time{}.Add(time.Second),
	//	ServiceName:   "s",
	//	OperationName: "o",
	//}
	//
}
