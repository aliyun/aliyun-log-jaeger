module github.com/aliyun/aliyun-log-jaeger

go 1.16

require (
	github.com/aliyun/aliyun-log-go-sdk v0.1.21
	github.com/gogo/protobuf v1.3.2
	github.com/hashicorp/go-hclog v0.16.2
	github.com/jaegertracing/jaeger v1.24.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/viper v1.8.1
)

replace github.com/aliyun/aliyun-log-jaeger  => ./
