package sls_store

import (
	"fmt"
	"time"

	"github.com/jaegertracing/jaeger/model"
	"github.com/jaegertracing/jaeger/storage/spanstore"
)

func toGetTraceQuery(id model.TraceID) string {
	return fmt.Sprintf(GetTraceQueryTemplate, id.String())
}

func toGetServicesQuery() string {
	return GetServiceQueryString
}

var queryBuilder = QueryBuilder{
	query:   "*",
	analyze: " select * from log where 1=1 ",
}

func toOperationsQuery(parameters spanstore.OperationQueryParameters) string {
	return queryBuilder.withSpanKind(parameters.SpanKind).withServiceName(parameters.ServiceName).toString()
}

func toFindTraceIdsQuery(parameters *spanstore.TraceQueryParameters) string {
	return QueryBuilder{
		query:   "*",
		analyze: "select traceid from log where 1=1 ",
	}.withTags(parameters.Tags).
		withDuration(parameters.DurationMin, parameters.DurationMax).
		withServiceName(parameters.ServiceName).
		withOperationName(parameters.OperationName).
		withGroupByTraceID().
		withLimit(parameters.NumTraces).
		toString()
}

type QueryBuilder struct {
	query   string
	analyze string
}

func (o QueryBuilder) withLimit(p int) *QueryBuilder {
	o.analyze += fmt.Sprintf(" limit %d", p)
	return &o
}
func (o QueryBuilder) withGroupByTraceID() *QueryBuilder {
	o.analyze += " group by traceid"
	return &o
}

func (o QueryBuilder) withSpanKind(p string) *QueryBuilder {
	if p != "" {
		o.query += fmt.Sprintf(" and kind: %s", p)
	}

	return &o
}

func (o QueryBuilder) withServiceName(p string) *QueryBuilder {
	if p != "" {
		o.query += fmt.Sprintf(" and service: %s", p)
	}

	return &o
}

func (o QueryBuilder) withOperationName(p string) *QueryBuilder {
	if p != "" {
		o.analyze += fmt.Sprintf(" and name like '%s'", p)
	}

	return &o
}

func (o QueryBuilder) withTags(p map[string]string) QueryBuilder {
	if len(p) == 0 {
		return o
	}

	for key, value := range p {
		o.query += fmt.Sprintf(" and attribute.%s: %s", key, value)
	}

	return o
}

func (o QueryBuilder) withDuration(min, max time.Duration) QueryBuilder {
	if min > 0 {
		o.analyze += fmt.Sprintf(" and duration >= %d", min.Nanoseconds()/1000)
	}

	if max > 0 {
		o.analyze += fmt.Sprintf(" and duration <= %d", max.Nanoseconds()/1000)
	}

	return o
}

func (o QueryBuilder) toString() string {
	return fmt.Sprintf("%s | %s", o.query, o.analyze)
}
