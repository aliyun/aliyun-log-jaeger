# Jaeger on Aliyun Log Service

[![License](https://img.shields.io/badge/license-Apache2.0-blue.svg)](/LICENSE)

[README in English](/README.md)

## 内容

* [简介](#简介)
* [架构](#架构)
  * [Jaeger client libraries](#jaeger-client-libraries)
  * [Agent](#agent)
  * [Collector](#collector)
  * [Query](#query)
  * [日志服务](#日志服务)
* [构建](#构建)
  * [Linux](#linux)
  * [Windows](#windows)
  * [MacOS](#macos)
* [配置 &amp; 部署](#配置--部署)
  * [日志服务](#日志服务-1)
  * [Agent](#agent-1)
  * [Collector](#collector-1)
  * [Query Service &amp; UI](#query-service--ui)
  * [Docker Compose](#docker-compose)
* [示例](#示例)
* [错误诊断](#错误诊断)
* [联系我们](#联系我们)
* [贡献者](#贡献者)

## 简介

[Jaeger](http://jaeger.readthedocs.io/en/latest/) 是 Uber 推出的一款开源分布式追踪系统，为微服务场景而生。它主要用于分析多个服务的调用过程，图形化服务调用轨迹，是诊断性能问题、分析系统故障的利器。

Jaeger on Aliyun Log Service 是基于 Jeager 开发的分布式追踪系统，支持将采集到的追踪数据持久化到[日志服务](https://help.aliyun.com/product/28958.html)中，并通过 Jaeger 的原生接口进行查询和展示。

## 架构

![architecture.png](/pics/architecture.png)

### Jaeger client libraries

Jaeger client 为不同语言实现了符合 [OpenTracing](http://opentracing.io/) 标准的 SDK。应用程序通过 API 写入数据，client library 把 trace 信息按照应用程序指定的采样策略传递给 jaeger-agent。数据使用 Thrift 序列化，通过 UDP 进行通信。

### Agent

Agent 是一个监听在 UDP 端口上接收 span 数据的网络守护进程，它会将数据批量发送给 collector。它被设计成一个基础组件，部署到所有的宿主机上。Agent 将 client library 和 collector 解耦，为 client library 屏蔽了路由和发现 collector 的细节。

### Collector

接收 jaeger-agent 发送来的数据，然后将数据写入后端存储。后端存储是一个可插拔的组件，Jaeger on Aliyun Log Service 增加了对阿里云日志服务的支持。

### Query

接收查询请求，从后端存储系统中检索 trace 并通过 UI 进行展示。

### 日志服务

Collector 会将接收到的 span 数据持久化到日志服务中。Query 会从日志服务中检索数据。

## 构建

Jaeger 提供了 docker 镜像能够让您方便地运行各个组件。但是，如果您的环境中无法使用 docker，您也可以直接基于源码构建出能够在相应平台上运行的二进制文件。

开始之前，请确保将该项目克隆到 `$GOPATH` 下的正确位置 `github.com/jaegertracing/jaeger`
```
mkdir -p $GOPATH/src/github.com/jaegertracing
cd $GOPATH/src/github.com/jaegertracing
git clone https://github.com/aliyun/aliyun-log-jaeger.git jaeger
cd jaeger
```

使用如下命令安装依赖
```
git submodule update --init --recursive
make install
```

使用下列命令构建出能够在相应平台上运行的组件：agent、collector 和 query。

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

## 配置 & 部署

### 日志服务

您需要按照以下步骤配置日志服务。

* 登录 [日志服务管理控制台](https://sls.console.aliyun.com/#/)。
* 创建用于存储 span 的 project、logstore。
* 为下列字段创建索引。

| 字段名 | 类型 | 分词符 |
| --- | --- | --- |
| traceID | text | N/A |
| spanID | text | N/A |
| process.serviceName | text | N/A |
| operationName | text | N/A |
| startTime | long | N/A |
| duration | long | N/A |

**注意**：如果查询时需要通过标签进行过滤，还需要为相应的标签字段创建索引。例如，应用程序会生成标签 http.method，http.status_code，并且需要根据这些标签进行查询，可以按下表所示创建索引。

| 字段名 | 类型 | 分词符 |
| --- | --- | --- |
| tags.http.method | text | N/A |
| tags.http.status_code | text | N/A |

### Agent

jaeger-agent 需要运行在包含 jaeger client libraries 应用程序的宿主机上。

agent 暴露如下端口

| 端口号 | 协议 | 功能 |
| --- | --- | --- |
| 5775 | UDP | 通过兼容性 thrift 协议，接收 zipkin thrift 类型的数据 |
| 6831 | UDP | 通过兼容性 thrift 协议，接收 jaeger thrift 类型的数据 |
| 6832 | UDP | 通过二进制 thrift 协议，接收 jaeger thrift 类型的数据 |
| 5778 | HTTP | 可用于配置采样策略 |

如果您的环境中有docker，可以使用如下方式运行 agent
```
docker run \
  --rm \
  -p5775:5775/udp \
  -p6831:6831/udp \
  -p6832:6832/udp \
  -p5778:5778/tcp \
  jaegertracing/jaeger-agent --collector.host-port=<JAEGER_COLLECTOR_HOST>:14267
```

如果您已构建好相应的二进制文件，这里以 macOS 为例，可以使用如下方式运行 agent
```
./cmd/agent/agent-darwin --collector.host-port=localhost:14267
```

### Collector

Collector 是无状态的，因此您可以同时运行任意数量的 jaeger-collector。运行 collector 需要指定用于存储 Span 的存储系统类型。如果指定的存储系统类型为日志服务，您还需要提供连接日志服务所需的相关参数。

参数说明如下

| 参数名 | 参数类型 | 描述 |
| --- | --- | --- |
| SPAN_STORAGE_TYPE | 环境变量 | 指定用于存储 Span 的存储系统类型。例如，aliyun-log |
| aliyun-log.project | 程序参数 | 指定用于存储 Span 的 Project |
| aliyun-log.endpoint | 程序参数 | 指定用于存储 Span 的 Project 所在的 Endpoint |
| aliyun-log.access-key-id | 程序参数 | 指定用户标识 Access Key ID |
| aliyun-log.access-key-secret | 程序参数 | 指定用户标识 Access Key Secret |
| aliyun-log.span-logstore | 程序参数 | 指定用于存储 Span 的 Logstore |

默认情况下，collector 暴露如下端口

| 端口号 | 协议 | 功能 |
| --- | --- | --- |
| 14267 | TChannel | 用于接收  jaeger-agent 发送来的 jaeger.thrift 格式的 span |
| 14268 | HTTP | 能直接接收来自客户端的 jaeger.thrift 格式的 span |
| 9411 | HTTP | 能通过 JSON 或 Thrift 接收 Zipkin spans，默认关闭 |

如果您的环境中有docker，可以使用如下方式运行 collector
```
docker run \
  -it --rm \
  -p14267:14267 -p14268:14268 -p9411:9411 \
  -e SPAN_STORAGE_TYPE=aliyun-log \
  registry.cn-hangzhou.aliyuncs.com/jaegertracing/jaeger-collector:0.0.2 \
  /go/bin/collector-linux \
  --aliyun-log.project=<PROJECT> \
  --aliyun-log.endpoint=<ENDPOINT> \
  --aliyun-log.access-key-id=<ACCESS_KEY_ID> \
  --aliyun-log.access-key-secret=<ACCESS_KEY_SECRET> \
  --aliyun-log.span-logstore=<SPAN_LOGSTORE>
```

如果您已构建好相应的二进制文件，这里以 macOS 为例，可以使用如下方式运行 collector
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

jaeger-query 提供了 API 端口以及 React/Javascript UI。该服务是无状态的，通常情况下它运行在 nginx 这样的负载均衡器后面。和 collector 类似，运行 query 同样需要指定后端存储系统类型。如果指定的存储系统类型为日志服务，您还需要提供连接日志服务所需的相关参数。此外，您还需要通过 query.static-files 参数指定 UI 静态文件的位置。

参数说明如下

| 参数名 | 参数类型 | 描述 |
| --- | --- | --- |
| SPAN_STORAGE_TYPE | 环境变量 | 指定用于存储 Span 的存储系统类型。例如，aliyun-log |
| aliyun-log.project | 程序参数 | 指定用于存储 Span 的 Project |
| aliyun-log.endpoint | 程序参数 | 指定用于存储 Span 的 Project 所在的 Endpoint |
| aliyun-log.access-key-id | 程序参数 | 指定用户标识 Access Key ID |
| aliyun-log.access-key-secret | 程序参数 | 指定用户标识 Access Key Secret |
| aliyun-log.span-logstore | 程序参数 | 指定用于存储 Span 的 Logstore |
| query.static-files | 程序参数 | 指定 UI 静态文件的位置 |

默认情况下，query 暴露如下端口

| 端口号 | 协议 | 功能 |
| --- | --- | --- |
| 16686 | HTTP | 1. /api/* - API 端口路径 </br> 2. / - Jaeger UI 路径 |

如果您的环境中有docker，可以使用如下方式运行 query
```
docker run \
  -it --rm \
  -p16686:16686 \
  -e SPAN_STORAGE_TYPE=aliyun-log \
  registry.cn-hangzhou.aliyuncs.com/jaegertracing/jaeger-query:0.0.2 \
  /go/bin/query-linux \
  --aliyun-log.project=<PROJECT> \
  --aliyun-log.endpoint=<ENDPOINT> \
  --aliyun-log.access-key-id=<ACCESS_KEY_ID> \
  --aliyun-log.access-key-secret=<ACCESS_KEY_SECRET> \
  --aliyun-log.span-logstore=<SPAN_LOGSTORE> \
  --query.static-files=/go/jaeger-ui/
```

如果您已构建好相应的二进制文件，这里以 macOS 为例，可以使用如下方式运行 query
```
export SPAN_STORAGE_TYPE=aliyun-log && \
  ./cmd/query/query-darwin \
  --aliyun-log.project=<PROJECT> \
  --aliyun-log.endpoint=<ENDPOINT> \
  --aliyun-log.access-key-id=<ACCESS_KEY_ID> \
  --aliyun-log.access-key-secret=<ACCESS_KEY_SECRET> \
  --aliyun-log.span-logstore=<SPAN_LOGSTORE> \
  --query.static-files=./jaeger-ui-build/build/
```

### Docker Compose

为了简化部署，我们提供了一个 docker-compose 模板 [aliyunlog-jaeger-docker-compose.yml](/docker-compose/aliyunlog-jaeger-docker-compose.yml)。

您可以通过如下命令将 jaeger-agent，jaeger-collector，jaeger-query 运行起来
```
docker-compose -f aliyunlog-jaeger-docker-compose.yml up
```

您可以通过如下命令将 jaeger-agent，jaeger-collector，jaeger-query 停止
```
docker-compose -f aliyunlog-jaeger-docker-compose.yml stop
```

**注意**：运行该命令之前请替换如下参数为真实值 ${PROJECT}、${ENDPOINT}、${ACCESS_KEY_ID}、${ACCESS_KEY_SECRET}、${SPAN_LOGSTORE}

## 示例

查询 trace

![traces.png](/pics/traces.png)

trace 详细信息

![trace_detail.png](/pics/trace_detail.png)

项目提供了一个名为 hotrod 的演示程序，具体内容请参考此[文档](/examples/hotrod/README_CN.md)。

## 错误诊断

如果您发现数据没有写入日志服务，可通过如下步骤进行错误诊断。

* 追踪数据会首先被宿主机上的 jaeger-agent 收集，请检查 jaeger-agent 是否启动成功，5775、6831、6832这几个用于接收数据的 UDP 端口的连通性。
* 如果 jaeger-agent 启动成功而且相应的端口都可连通，下一步请检查 jaeger-agent 和 jaeger-collector 的连通性。如果jaeger-agent 成功连接 jaeger-collector 会通过标准输出打印出如下信息`"msg":"Connected to peer"`，否则，会持续输出`"msg":"Unable to connect"`，或者在尝试提交数据的时候输出`"msg":"Could not submit jaeger batch","error":"no peers available"`。
* 如果 jaeger-agent 和 jaeger-collector 连接成功，请检查 jaeger-collector 和日志服务的连接问题。检查 jaeger-collector 的标准输出`"msg":"Failed to write span"`打印的错误原因。

## 联系我们

- [阿里云LOG官方网站](https://www.aliyun.com/product/sls/)
- [阿里云LOG官方论坛](https://yq.aliyun.com/groups/50)
- 阿里云官方技术支持：[提交工单](https://workorder.console.aliyun.com/#/ticket/createIndex)

## 贡献者

[@WPH95](https://github.com/WPH95) 对项目作了很大贡献。

感谢 [@WPH95](https://github.com/WPH95) 的杰出工作。
