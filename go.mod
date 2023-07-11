module github.com/aliyun/aliyun-log-jaeger

go 1.16

require (
	github.com/aliyun/aliyun-log-go-sdk v0.1.21
	github.com/gogo/protobuf v1.3.2
	github.com/hashicorp/go-hclog v1.5.0
	github.com/jaegertracing/jaeger v1.47.0
	github.com/spf13/cast v1.5.1
	github.com/spf13/viper v1.16.0
)

replace github.com/aliyun/aliyun-log-jaeger => ./
