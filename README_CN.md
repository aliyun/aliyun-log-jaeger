# Jaeger on Aliyun Log Service

[README in English](/README.md)

## 简介

[Jaeger](http://jaeger.readthedocs.io/en/latest/) 是 Uber 推出的一款开源分布式追踪系统，为微服务场景而生。它主要用于分析多个服务的调用过程，图形化服务调用轨迹，是诊断性能问题、分析系统故障的利器。

Jaeger on Aliyun Log Service 是基于 Jeager 开发的分布式追踪系统，支持将采集到的追踪数据持久化到[日志服务](https://help.aliyun.com/product/28958.html)中，并通过 Jaeger 的原生接口进行查询和展示。

## 构建

Jaeger 提供了 docker 镜像能够让您方便地运行各个组件。但是，如果您的环境中无法使用 docker，您也可以直接基于源码构建出能够在相应平台上运行的二进制文件。

开始之前，请确保将该项目克隆到 `$GOPATH` 下的正确位置 `github.com/jaegertracing/jaeger`
```
mkdir -p $GOPATH/src/github.com/jaegertracing
cd $GOPATH/src/github.com/jaegertracing
git clone git@github.com:jaegertracing/jaeger.git jaeger
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

## 部署

Jaeger 后端组件分为 jaeger-agent，jaeger-collector 和 jaeger-query。

### Agent

由于我们并未对 agent 作任何修改，关于 agent 的部署方式请参考[原始文档](http://jaeger.readthedocs.io/en/latest/deployment/#agent)。

### Collectors

Collector 是无状态的，因此您可以同时运行任意数量的 jaeger-collector。运行 collector 需要指定用于存储 Span 的存储系统类型。如果指定的存储系统类型为日志服务，您还需要提供连接日志服务所需的相关参数。

参数说明如下

| 参数名 | 参数类型 | 描述 |
| --- | --- | --- |
| SPAN_STORAGE_TYPE | 环境变量 | 指定用于存储 Span 的存储系统类型。例如，aliyun-log |
| aliyun-log.project | 程序参数 | 指定用于存储 Span 的 Project |
| aliyun-log.endpoint | 程序参数 | 指定用于存储 Span 的 Project 所在的 Endpoint |
| aliyun-log.access-key-id | 程序参数 | 指定用户标识 Access Key ID |
| aliyun-log.access-key-secret | 程序参数 | 指定用户标识 Access Key Secret |
| aliyun-log.logstore | 程序参数 | 指定用于存储 Span 的 Logstore |

默认情况下，collector 暴露如下端口

| 端口号 | 协议 | 功能 |
| --- | --- | --- |
| 14267 | TChannel | 用于接收  jaeger-agent 发送来的 jaeger.thrift 格式的 span |
| 14268 | HTTP | 能直接接收来自客户端的 jaeger.thrift 格式的 span |
| 9411 | HTTP | 能通过 JSON 或 Thrift 接收 Zipkin spans |

如果您的环境中有docker，可以使用如下方式运行 collector
```
docker run -it --rm -e SPAN_STORAGE_TYPE=aliyun-log registry.cn-hangzhou.aliyuncs.com/jaegertracing/jaeger-collector:0.0.1 /go/bin/collector-linux --aliyun-log.project=<PROJECT> --aliyun-log.endpoint=<ENDPOINT> --aliyun-log.access-key-id=<ACCESS_KEY_ID> --aliyun-log.access-key-secret=<ACCESS_KEY_SECRET> --aliyun-log.span-logstore=<SPAN_LOGSTORE>
```

如果您以构建好相应的二进制文件，可以使用如下方式运行 collector
```
export SPAN_STORAGE_TYPE=aliyun-log && ./cmd/collector/collector-darwin --aliyun-log.project=<PROJECT> --aliyun-log.endpoint=<ENDPOINT> --aliyun-log.access-key-id=<ACCESS_KEY_ID> --aliyun-log.access-key-secret=<ACCESS_KEY_SECRET> --aliyun-log.span-logstore=<SPAN_LOGSTORE>
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
| aliyun-log.logstore | 程序参数 | 指定用于存储 Span 的 Logstore |
| query.static-files | 程序参数 | 指定 UI 静态文件的位置 |

默认情况下，query 暴露如下端口

| 端口号 | 协议 | 功能 |
| --- | --- | --- |
| 16686 | HTTP | 1. /api/* - API 端口路径 </br> 2. / - Jaeger UI 路径 |

如果您的环境中有docker，可以使用如下方式运行 query
```
docker run -it --rm -e SPAN_STORAGE_TYPE=aliyun-log registry.cn-hangzhou.aliyuncs.com/jaegertracing/jaeger-query:0.0.1 /go/bin/collector-linux --aliyun-log.project=<PROJECT> --aliyun-log.endpoint=<ENDPOINT> --aliyun-log.access-key-id=<ACCESS_KEY_ID> --aliyun-log.access-key-secret=<ACCESS_KEY_SECRET> --aliyun-log.span-logstore=<SPAN_LOGSTORE> --query.static-files=/go/jaeger-ui/
```

如果您以构建好相应的二进制文件，可以使用如下方式运行 query
```
export SPAN_STORAGE_TYPE=aliyun-log && ./cmd/query/query-darwin --aliyun-log.project=<PROJECT> --aliyun-log.endpoint=<ENDPOINT> --aliyun-log.access-key-id=<ACCESS_KEY_ID> --aliyun-log.access-key-secret=<ACCESS_KEY_SECRET> --aliyun-log.span-logstore=<SPAN_LOGSTORE> --query.static-files=./jaeger-ui-build/build/
```

### Docker Compose

为了简化部署，我们提供了一个 docker-compose 模板 [aliyun-jaeger-docker-compose.yml](/docker-compose/aliyun-jaeger-docker-compose.yml)。

您可以通过如下命令将 jaeger-agent，jaeger-collector，jaeger-query 运行起来
```
docker-compose -f aliyunlog-jaeger-docker-compose.yml up
```

您可以通过如下命令将 jaeger-agent，jaeger-collector，jaeger-query 停止
```
docker-compose -f aliyunlog-jaeger-docker-compose.yml stop
```

**注意**：运行该命令之前请替换如下参数为真实值 ${PROJECT}、${ENDPOINT}、${ACCESS_KEY_ID}、${ACCESS_KEY_SECRET}、${LOGSTORE}