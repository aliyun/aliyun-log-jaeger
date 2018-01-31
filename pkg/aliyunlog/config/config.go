package config

import "time"

// Configuration describes the configuration properties needed to connect to an AliCloud Log Service cluster
type Configuration struct {
	Project            string
	Endpoint           string
	AccessKeyID        string
	AccessKeySecret    string
	SpanLogstore       string
	DependencyLogstore string
	LogstoreShardCount int
	LogstoreShardTTL   int
	MaxSpanAge         time.Duration
}

// ApplyDefaults copies settings from source unless its own value is non-zero.
func (c *Configuration) ApplyDefaults(source *Configuration) {
	if c.Project == "" {
		c.Project = source.Project
	}
	if c.Endpoint == "" {
		c.Endpoint = source.Endpoint
	}
	if c.AccessKeyID == "" {
		c.AccessKeyID = source.AccessKeyID
	}
	if c.AccessKeySecret == "" {
		c.AccessKeySecret = source.AccessKeySecret
	}
	if c.SpanLogstore == "" {
		c.SpanLogstore = source.SpanLogstore
	}
	if c.DependencyLogstore == "" {
		c.DependencyLogstore = source.DependencyLogstore
	}
	if c.LogstoreShardCount == 0 {
		c.LogstoreShardCount = source.LogstoreShardCount
	}
	if c.LogstoreShardTTL == 0 {
		c.LogstoreShardTTL = source.LogstoreShardTTL
	}
	if c.MaxSpanAge == 0 {
		c.MaxSpanAge = source.MaxSpanAge
	}
}
