This is the repository that contains object storage (Alibaba Could log service) plugin for Jaeger.

## About

As a component of an observability/monitoring system, Jaeger is an essential source of data for development and
operations students to locate and find problems and exceptions with the business system. As an SRE, we must ensure that
the monitoring system lives longer than the system. Once the monitoring system is down before the business system,
monitoring is entirely worthless. Monitoring is the last barrier for business exception analysis, and it is more
sensitive to high availability and high performance than other systems.

The [Alibaba Could log service](https://www.alibabacloud.com/product/log-service)(SLS) provides high performance,
resilience, and freedom from operation and maintenance, allowing users to cope with surge traffic or inaccurate size
assessment quickly, and the SLS service itself provides 99.9% availability and 11 out of 9 data reliability.

The Alibab Cloud log service  :heart:  Jaeger

## Build/Compile

In order to compile the plugin from source code you can use `go build`:

```shell
cd /path/to/jaeger-sls
go build
```

## Parameter Flag

(TODO)

## Start

(TODO)

## License

The SLS Storage gRPC Plugin for Jaeger is an [MIT licensed](LICENSE) open source project.