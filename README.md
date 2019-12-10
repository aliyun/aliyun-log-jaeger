# Jaeger on Aliyun Log Service

[![License](https://img.shields.io/badge/license-Apache2.0-blue.svg)](/LICENSE)

[中文版README](/README_CN.md)

## Content

* [Introduction](#introduction)
* [Architecture](#architecture)
  * [Jaeger client libraries](#jaeger-client-libraries)
  * [Agent](#agent)
  * [Collector](#collector)
  * [Query](#query)
  * [Aliyun Log Service](#aliyun-log-service)
* [Building](#building)
  * [Linux](#linux)
  * [Windows](#windows)
  * [MacOS](#macos)
* [Configure &amp; Deployment](#configure--deployment)
  * [Aliyun Log Service](#aliyun-log-service-1)
  * [Agent](#agent-1)
  * [Collector](#collector-1)
  * [Query Service &amp; UI](#query-service--ui)
  * [Docker Compose](#docker-compose)
* [Example](#example)
* [Contact Us](#contact-us)
* [Contributors ](#contributors)

## Introduction

[Jaeger](http://jaeger.readthedocs.io/en/latest/) is an opensource distributed tracing system developed by Uber, it is mainly used for micro service scenarios. It can be used to analyze the invocation process for multiple services, display the method call trace and the method call relations. It is a useful tool for diagnosing performance problems and analyzing system failures.

`Jaeger on Aliyun Log Service` is a distributed tracing system based on Jaeger which supports persist data into [Aliyun Log Service](https://help.aliyun.com/product/28958.html). What's more you can retrieve them from log service through jaeger-query and display them on Jaeger UI.

## Architecture

![architecture.png](/pics/architecture.png)

### Jaeger client libraries

Jaeger clients are language specific implementations of the [OpenTracing API](http://opentracing.io/). They can be used to instrument applications for distributed tracing either manually or with a variety of existing open source frameworks, such as Flask, Dropwizard, gRPC, and many more, that are already integrated with OpenTracing.

### Agent

A network daemon that listens for spans sent over UDP, which it batches and sends to the collector. It is designed to be deployed to all hosts as an infrastructure component. The agent abstracts the routing and discovery of the collectors away from the client.

### Collector

The collector receives traces from Jaeger agents and runs them through a processing pipeline. The storage is a pluggable component. `Jaeger on Aliyun Log Service` supports use Aliyun Log Service as the backend storage.

### Query

Query is a service that retrieves traces from storage and hosts a UI to display them.

### Aliyun Log Service

The jaeger-collector will persist the received data to the log service. The jaeger-query will retrieve data from the log service.

## Building

Jaeger provides docker images that allows you to run various components in a convenient way. However, if you can't use docker in your environment, you can also build binary files that can run on the corresponding platform based on the source code directly or use the [release packages](https://github.com/aliyun/aliyun-log-jaeger/releases/tag/0.2.4).

To get started, make sure you clone the Git repository into the correct location `github.com/jaegertracing/jaeger` relative to `$GOPATH`:
```
mkdir -p $GOPATH/src/github.com/jaegertracing
cd $GOPATH/src/github.com/jaegertracing
git clone https://github.com/aliyun/aliyun-log-jaeger.git jaeger
cd jaeger
```

Then install dependencies:
```
git submodule update --init --recursive
make install
```
Please use the following commands to build the components that can run on the corresponding platform.

### Linux

```
make build-all-linux
```

### Windows

```
make build-all-windows
```

### MacOS

```
make build-all-darwin
```

## Configure & Deployment

### Aliyun Log Service

Please configure the log service according to the following steps.

* Login on [Aliyun Log Service Web Console](https://sls.console.aliyun.com/#/). 
* Create project, logstore for storing span.
* Create indexes for the following fields.

| Field Name | Type | Token |
| --- | --- | --- |
| traceID | text | N/A |
| spanID | text | N/A |
| process.serviceName | text | N/A |
| operationName | text | N/A |
| startTime | long | N/A |
| duration | long | N/A |

**Note**: if you want to use tags as condition to find traces, you should alse create indexes for the tag fields. For example, the application generate the following tags http.method, http.status_code and you want to use them as condition to find traces, you should create indexes for them.

| Field Name | Type | Token |
| --- | --- | --- |
| tags.http.method | text | N/A |
| tags.http.status_code | text | N/A |

### Agent

Jaeger client libraries expect jaeger-agent process to run locally on each host. The agent exposes the following ports:

| Port | Protocol | Function |
| --- | --- | --- |
| 5775 | UDP | accept zipkin.thrift over compact thrift protocol |
| 6831 | UDP | accept jaeger.thrift over compact thrift protocol |
| 6832 | UDP | accept jaeger.thrift over binary thrift protocol |
| 5778 | HTTP | serve configs, sampling strategies |

If you have already installed docker, you can run agent as follows:
```
docker run \
  --rm \
  -p5775:5775/udp \
  -p6831:6831/udp \
  -p6832:6832/udp \
  -p5778:5778/tcp \
  jaegertracing/jaeger-agent --collector.host-port=<JAEGER_COLLECTOR_HOST>:14267
```

If you have already built the corresponding binary file, take macOS as an example, you can run agent as follows:
```
./cmd/agent/agent-darwin --collector.host-port=localhost:14267
```

### Collector

The collectors are stateless and thus many instances of jaeger-collector can be run in parallel. You need to specify the storage type used to store span. If you specify Aliyun Log Service as your backend storage, you also need to provide the relevant parameters for the log service.

Parameter Description

| Parameter Name | Type | Description |
| --- | --- | --- |
| SPAN_STORAGE_TYPE | environment variable | specify the storage type used to store span |
| aliyun-log.project | program argument | specify the project used to store span |
| aliyun-log.endpoint | program argument | specify the endpoint for your project |
| aliyun-log.access-key-id | program argument | specify the account information for your log services |
| aliyun-log.access-key-secret | program argument | specify the account information for your log services |
| aliyun-log.span-logstore | program argument | specify the logstore used to store span |

At default settings the collector exposes the following ports:

| Port | Protocol | Function |
| --- | --- | --- |
| 14267 | TChannel | used by **jaeger-agent** to send spans in jaeger.thrift format |
| 14268 | HTTP | can accept spans directly from clients in jaeger.thrift format |
| 9411 | HTTP | can accept Zipkin spans in JSON or Thrift (disabled by default) |

If you have already installed docker, you can run collector as follows:
```
docker run \
  -it --rm \
  -p14267:14267 -p14268:14268 -p9411:9411 \
  -e SPAN_STORAGE_TYPE=aliyun-log \
  registry.cn-hangzhou.aliyuncs.com/jaegertracing/jaeger-collector:0.2.4 \
  /go/bin/collector-linux \
  --aliyun-log.project=<PROJECT> \
  --aliyun-log.endpoint=<ENDPOINT> \
  --aliyun-log.access-key-id=<ACCESS_KEY_ID> \
  --aliyun-log.access-key-secret=<ACCESS_KEY_SECRET> \
  --aliyun-log.span-logstore=<SPAN_LOGSTORE>
```

If you have already built the corresponding binary file, take macOS as an example, you can run collector as follows:
```
export SPAN_STORAGE_TYPE=aliyun-log && \
  ./cmd/collector/collector-darwin \
  --aliyun-log.project=<PROJECT> \
  --aliyun-log.endpoint=<ENDPOINT> \
  --aliyun-log.access-key-id=<ACCESS_KEY_ID> \
  --aliyun-log.access-key-secret=<ACCESS_KEY_SECRET> \
  --aliyun-log.span-logstore=<SPAN_LOGSTORE>
```

### Query Service & UI

`jaeger-query` serves the API endpoints and a React/Javascript UI. The service is stateless and is typically run behind a load balancer, e.g. nginx. Similar to collector, if you specify Aliyun Log Service as your backend storage, you also need to provide the relevant parameters for the log service. In addition, you need to specify the location of the UI static file by the parameter `query.static-files`.

Parameters Description

| Parameter Name | Type | Description |
| --- | --- | --- |
| SPAN_STORAGE_TYPE | environment variable | specify the storage type used to store span |
| aliyun-log.project | program argument | specify the project used to store span |
| aliyun-log.endpoint | program argument | specify the endpoint for your project |
| aliyun-log.access-key-id | program argument | specify the account information for your log services |
| aliyun-log.access-key-secret | program argument | specify the account information for your log services |
| aliyun-log.span-logstore | program argument | specify the logstore used to store span |
| aliyun-log.span-agg-logstore | program argument | specify the logstore used to store agg data |
| query.static-files | program argument | Specify the location of the UI static files |

At default settings the query service exposes the following port(s):

| Port | Protocol | Function |
| --- | --- | --- |
| 16686 | HTTP | **/api/*** endpoints and Jaeger UI at / |

If you have already installed docker, you can run query as follows:
```
docker run \
  -it --rm \
  -p16686:16686 \
  -e SPAN_STORAGE_TYPE=aliyun-log \
  registry.cn-hangzhou.aliyuncs.com/jaegertracing/jaeger-query:0.2.4 \
  /go/bin/query-linux \
  --aliyun-log.project=<PROJECT> \
  --aliyun-log.endpoint=<ENDPOINT> \
  --aliyun-log.access-key-id=<ACCESS_KEY_ID> \
  --aliyun-log.access-key-secret=<ACCESS_KEY_SECRET> \
  --aliyun-log.span-logstore=<SPAN_LOGSTORE> \
  --aliyun-log.span-agg-logstore=<SPAN_AGG_LOGSTORE> \
  --query.static-files=/go/jaeger-ui/
```

If you have already built the corresponding binary file, take macOS as an example, you can run query as follows:
```
export SPAN_STORAGE_TYPE=aliyun-log && \
  ./cmd/query/query-darwin \
  --aliyun-log.project=<PROJECT> \
  --aliyun-log.endpoint=<ENDPOINT> \
  --aliyun-log.access-key-id=<ACCESS_KEY_ID> \
  --aliyun-log.access-key-secret=<ACCESS_KEY_SECRET> \
  --aliyun-log.span-logstore=<SPAN_LOGSTORE> \
  --aliyun-log.span-agg-logstore=<SPAN_AGG_LOGSTORE> \
  --query.static-files=./jaeger-ui-build/build/
```

### Docker Compose

To simplify the deployment, we have provided a docker-compose template [aliyunlog-jaeger-docker-compose.yml](/docker-compose/aliyunlog-jaeger-docker-compose.yml).

You can start `jaeger-agent`, `jaeger-collector`, and `jaeger-query` through the following commands
```
docker-compose -f aliyunlog-jaeger-docker-compose.yml up
```

You can stop `jaeger-agent`, `jaeger-collector`, and `jaeger-query` through the following commands
```
docker-compose -f aliyunlog-jaeger-docker-compose.yml stop
```

**Note**: please remember to replace the following parameters with the real value before you run the above commands.


## Example

Find traces

![traces.png](/pics/traces.png)

Display detailed information for trace

![trace_detail.png](/pics/trace_detail.png)

This project provide a demo applicatio named hotrod. Please refer to this [doc](/examples/hotrod/README.md).

## Contact Us
- [Alicloud Log Service homepage](https://www.alibabacloud.com/product/log-service)
- [Alicloud Log Service doc](https://www.alibabacloud.com/help/product/28958.htm)
- [Alicloud Log Servic official forum](https://yq.aliyun.com/groups/50)
- Alicloud Log Servic official technical support: [submit tickets](https://workorder.console.aliyun.com/#/ticket/createIndex)

## Contributors
[@WPH95](https://github.com/WPH95) made a great contribution to this project.

Thanks for the excellent work by [@WPH95](https://github.com/WPH95)
