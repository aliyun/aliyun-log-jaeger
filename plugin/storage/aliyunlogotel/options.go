// Copyright (c) 2017 Uber Technologies, Inc.
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

package aliyunlogotel

import (
	"flag"
	"time"

	"github.com/jaegertracing/jaeger/pkg/aliyunlog/config"
	"github.com/spf13/viper"
)

const (
	suffixProject          = ".project"
	suffixEndpoint         = ".endpoint"
	suffixAliCloudK8S      = ".alicloud-k8s-flag"
	suffixAccessKeyID      = ".access-key-id"
	suffixAccessKeySecret  = ".access-key-secret"
	suffixSpanLogstore     = ".span-logstore"
	suffixSpanAggLogstore  = ".span-agg-logstore"
	suffixSpanDepLogstore  = ".span-dep-logstore"
	suffixMaxQueryDuration = ".max-query-duration"
	suffixInitResourceFlag = ".init-resource-flag"
	suffixTagAppenderRule  = ".tag-appender-rules"
	suffixKindRewriteRule  = ".kind-rewrite-rules"
)

// Options contains various type of AliCloud Log Service configs and provides the ability
// to bind them to command line flag and apply overlays, so that some configurations
// (e.g. archive) may be underspecified and infer the rest of its parameters from primary.
type Options struct {
	primary *namespaceConfig

	others map[string]*namespaceConfig
}

type namespaceConfig struct {
	config.Configuration
	namespace string
}

// NewOptions creates a new Options struct.
func NewOptions(primaryNamespace string, otherNamespaces ...string) *Options {
	// TODO all default values should be defined via cobra flags
	options := &Options{
		primary: &namespaceConfig{
			Configuration: config.Configuration{
				Project:          "",
				Endpoint:         "",
				AliCloudK8SFlag:  false,
				AccessKeyID:      "",
				AccessKeySecret:  "",
				SpanLogstore:     "jaeger-span",
				SpanAggLogstore:  "",
				MaxQueryDuration: 24 * time.Hour,
				InitResourceFlag: true,
				TagAppendRuleFile: "",
				KindRewriteRuleFile: "",
			},
			namespace: primaryNamespace,
		},
		others: make(map[string]*namespaceConfig, len(otherNamespaces)),
	}

	for _, namespace := range otherNamespaces {
		options.others[namespace] = &namespaceConfig{namespace: namespace}
	}

	return options
}

// AddFlags adds flags for Options
func (opt *Options) AddFlags(flagSet *flag.FlagSet) {
	addFlags(flagSet, opt.primary)
	for _, cfg := range opt.others {
		addFlags(flagSet, cfg)
	}
}

func addFlags(flagSet *flag.FlagSet, nsConfig *namespaceConfig) {
	flagSet.String(
		nsConfig.namespace+suffixProject,
		nsConfig.Project,
		"The project required by AliCloud Log Service")
	flagSet.String(
		nsConfig.namespace+suffixEndpoint,
		nsConfig.Endpoint,
		"The endpoint required by AliCloud Log Service i.e cn-hangzhou.log.aliyuncs.com")
	flagSet.Bool(
		nsConfig.namespace+suffixAliCloudK8S,
		nsConfig.AliCloudK8SFlag,
		"Set this flag true if jaeger deploy in AliCloud kubernetes cluster, and you don't need to set access key pair")
	flagSet.String(
		nsConfig.namespace+suffixAccessKeyID,
		nsConfig.AccessKeyID,
		"The access-key-id required by AliCloud Log Service")
	flagSet.String(
		nsConfig.namespace+suffixAccessKeySecret,
		nsConfig.AccessKeySecret,
		"The access-key-secret required by AliCloud Log Service")
	flagSet.String(
		nsConfig.namespace+suffixSpanLogstore,
		nsConfig.SpanLogstore,
		"The logstore to save span data in AliCloud Log Service")
	flagSet.String(
		nsConfig.namespace+suffixSpanAggLogstore,
		nsConfig.SpanAggLogstore,
		"The agg logstore to save span data in AliCloud Log Service")
	flagSet.String(
		nsConfig.namespace+suffixSpanDepLogstore,
		nsConfig.SpanDepLogstore,
		"The dependency logstore to save span data in AliCloud Log Service")
	flagSet.Duration(
		nsConfig.namespace+suffixMaxQueryDuration,
		nsConfig.MaxQueryDuration,
		"The maximum query duration for logstore in AliCloud Log Service")
	flagSet.Bool(
		nsConfig.namespace+suffixInitResourceFlag,
		nsConfig.InitResourceFlag,
		"The flag to specify whether to init resource in AliCloud Log Service")
	flagSet.String(
		nsConfig.namespace+suffixTagAppenderRule,
		nsConfig.TagAppendRuleFile,
		"The file of rule which appending tag to span.")

	flagSet.String(
		nsConfig.namespace+suffixKindRewriteRule,
		nsConfig.KindRewriteRuleFile,
		"The file of rule which rewrite span kind.")

}

// InitFromViper initializes Options with properties from viper
func (opt *Options) InitFromViper(v *viper.Viper) {
	initFromViper(opt.primary, v)
	for _, cfg := range opt.others {
		initFromViper(cfg, v)
	}
}

func initFromViper(cfg *namespaceConfig, v *viper.Viper) {
	cfg.Project = v.GetString(cfg.namespace + suffixProject)
	cfg.Endpoint = v.GetString(cfg.namespace + suffixEndpoint)
	cfg.AliCloudK8SFlag = v.GetBool(cfg.namespace + suffixAliCloudK8S)
	cfg.AccessKeyID = v.GetString(cfg.namespace + suffixAccessKeyID)
	cfg.AccessKeySecret = v.GetString(cfg.namespace + suffixAccessKeySecret)
	cfg.SpanLogstore = v.GetString(cfg.namespace + suffixSpanLogstore)
	cfg.SpanAggLogstore = v.GetString(cfg.namespace + suffixSpanAggLogstore)
	cfg.DependencyLogstore = v.GetString(cfg.namespace + suffixSpanDepLogstore)
	cfg.MaxQueryDuration = v.GetDuration(cfg.namespace + suffixMaxQueryDuration)
	cfg.InitResourceFlag = v.GetBool(cfg.namespace + suffixInitResourceFlag)
	cfg.TagAppendRuleFile = v.GetString(cfg.namespace + suffixTagAppenderRule)
	cfg.KindRewriteRuleFile = v.GetString(cfg.namespace + suffixKindRewriteRule)
}

// GetPrimary returns primary configuration.
func (opt *Options) GetPrimary() *config.Configuration {
	return &opt.primary.Configuration
}

// Get returns auxiliary named configuration.
func (opt *Options) Get(namespace string) *config.Configuration {
	nsCfg, ok := opt.others[namespace]
	if !ok {
		nsCfg = &namespaceConfig{}
		opt.others[namespace] = nsCfg
	}
	nsCfg.Configuration.ApplyDefaults(&opt.primary.Configuration)
	return &nsCfg.Configuration
}
