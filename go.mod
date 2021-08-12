module github.com/qiansheng91/jaeger-sls

go 1.16

require (
	github.com/aliyun/aliyun-log-go-sdk v0.1.21
	github.com/gogo/protobuf v1.3.2
	github.com/hashicorp/go-hclog v0.16.2
	github.com/jaegertracing/jaeger v1.24.0
	github.com/spf13/cast v1.3.1
	github.com/spf13/viper v1.8.1
	google.golang.org/grpc v1.39.1 // indirect
)

replace . => github.com/qiansheng91/jaeger-sls v0.0.0-20210803014446-26eb89a251e1
