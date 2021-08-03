package sls_store

type slsTraceInstance interface {
	project() string
	traceLogStore() string
	serviceDependencyLogStore() string
}

func newSlsTraceInstance(project, instance string) slsTraceInstance {
	return &slsTraceInstanceImpl{
		projectName:                   project,
		instance:                      instance,
		traceLogStoreName:             instance + "-traces",
		serviceDependencyLogStoreName: instance + "-traces-deps",
	}
}

type slsTraceInstanceImpl struct {
	instance                      string
	traceLogStoreName             string
	serviceDependencyLogStoreName string
	projectName                   string
}

func (s *slsTraceInstanceImpl) project() string {
	return s.projectName
}

func (s *slsTraceInstanceImpl) traceLogStore() string {
	return s.traceLogStoreName
}

func (s *slsTraceInstanceImpl) serviceDependencyLogStore() string {
	return s.serviceDependencyLogStoreName
}
