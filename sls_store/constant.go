package sls_store

import "time"

// Span Attribute name List
const (
	// ParentService The log item key of parent service name
	ParentService = "parent_service"
	// ChildService the field name of child service
	ChildService = "child_service"
	// ServiceName the field name of service
	ServiceName = "service"
	// OperationName the field name of operation name
	OperationName = "name"
	// SpanKind  the field name of span kind
	SpanKind = "kind"
	// TraceID the field name of trace id
	TraceID = "traceID"
	// TraceIDField
	TraceIDField = "traceid"
	// SpanID the field name of span id
	SpanID = "spanID"
	// SpanIDField
	SpanIDField = "spanID"
	// ParentSpanID the field name of parent span id
	ParentSpanID = "parentSpanID"
	// StartTime the field name of start time
	StartTime = "start"
	// Duration the field name of duration
	Duration = "duration"
	// Attribute the field name of span tags
	Attribute = "attribute"
	// Resource the field name of span process tag
	Resource = "resource"
	// Logs the field name of span log
	Logs = "logs"
	// Links the field name of span reference
	Links = "links"
	// StatusMessage the field name of warning message of span
	StatusMessage = "statusMessage"
	//StatusMessageField
	StatusMessageField = "statusmessage"
	// Flags the field name of flags
	Flags = "flags"
	// EndTime the field name of end time
	EndTime = "end"
	// StatusCode the field name of status code
	StatusCode = "statusCode"
	// StatusCodeField
	StatusCodeField = "statuscode"
)

// Query template List
const (
	// DependenciesQueryString The query string which calculates the dependency relationship between each service.
	DependenciesQueryString = "* and version: service_name | SELECT  parent_service, child_service , sum(n_status_fail + n_status_succ) as count from log group by parent_service, child_service"
	// GetTraceQueryTemplate The template query string which selects trace by trace id
	GetTraceQueryTemplate = "traceID: %s"
	// GetServiceQueryString the query string which queries all service name
	GetServiceQueryString = "* | select DISTINCT service"
)

// query operation values
const (
	// DefaultFetchNumber the max fetching number of span
	DefaultFetchNumber = 1000
	// DefaultOffset  default offset
	DefaultOffset = 0
	// DefaultTopicName default topic name
	DefaultTopicName = ""
	// DefaultRetryTimeOut the default value of retry timeout
	DefaultRetryTimeOut = 2 * time.Minute
	// DefaultRequestTimeOut the default value of request timeout
	DefaultRequestTimeOut = 2 * time.Minute
)
