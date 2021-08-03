package main

import (
	"flag"
	"github.com/jaegertracing/jaeger/plugin/storage/grpc"
	"github.com/jaegertracing/jaeger/plugin/storage/grpc/shared"
	"github.com/qiansheng91/jaeger-sls/sls_store"
	"time"
)

const (
	DefaultLookBack = 7
)

var (
	endpoint     string
	accessKeyID  string
	accessSecret string
	project      string
	instance     string
	lookBack     int64
)

func main() {
	flag.StringVar(&endpoint, "sls-endpoint", "", "The endpoint of Log Service. The format is ${project}.${region-endpoint}")
	flag.StringVar(&accessKeyID, "sls-access_key-id", "", "The AccessKey ID of your Alibaba Cloud account.")
	flag.StringVar(&accessSecret, "sls-access-secret", "", "The AccessKey secret of your Alibaba Cloud account.")
	flag.StringVar(&project, "sls-project", "", "The name of the Log Service project.")
	flag.StringVar(&instance, "sls-instance", "", "The name of the trace instance.")
	flag.Int64Var(&lookBack, "sls-max-look-back", DefaultLookBack, "Maximum time frame for searching data. (Unit: Day)")
	flag.Parse()

	checkParameter(endpoint, accessKeyID, accessSecret, project, instance)

	var plugin = sls_store.NewSLSStorageForJaegerPlugin(
		endpoint,
		accessKeyID,
		accessSecret,
		project,
		instance,
		time.Duration(lookBack)*24*time.Hour,
	)

	grpc.Serve(&shared.PluginServices{
		Store:        plugin,
		ArchiveStore: plugin,
	})
}

func checkParameter(endpoint string, accessKeyID string, accessSecret string, project string, instance string) {
	if endpoint == "" {
		panic("The Endpoint can't be empty")
	}

	if accessKeyID == "" {
		panic("The access key id can't be empty")
	}

	if accessSecret == "" {
		panic("The access secret can't be empty")
	}

	if project == "" {
		panic("The access secret can't be empty")
	}

	if instance == "" {
		panic("The access secret can't be empty")
	}
}
