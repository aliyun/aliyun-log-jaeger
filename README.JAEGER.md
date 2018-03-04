# 分布式追踪利器 Jaeger 上云实战

## 分布式追踪面临的挑战
近年来，随着微服务架构的兴起，一些在单机环境下容易处理的问题变得越来越困难：
* 如何清晰地展示各服务间的调用关系？
* 一次请求的流量从哪里来，最终落到哪个服务中去？
* 如何找出系统的性能瓶颈？

为了解决这些问题出现了许多分布式追踪系统，包括 Dapper，Zipkin，HTrace，EagleEye 等。但不同系统的 API 并不兼容，这就导致了如果您希望将追踪系统由 Zipkin 替换为 HTrace，往往会带来较大改动。

## OpenTracing
为了解决不同的分布式追踪系统 API 不兼容的问题，诞生了 [OpenTracing](http://opentracing.io/) 规范。
OpenTracing 通过提供平台无关、厂商无关的API，使得开发人员能够方便的添加（或更换）追踪系统的实现。OpenTracing 正在为全球的分布式追踪，提供统一的概念和数据标准。

OpenTracing是一个轻量级的标准化层，它位于**应用程序/类库**和**追踪或日志分析程序**之间。
```
+-------------+  +---------+  +----------+  +------------+
| Application |  | Library |  |   OSS    |  |  RPC/IPC   |
|    Code     |  |  Code   |  | Services |  | Frameworks |
+-------------+  +---------+  +----------+  +------------+
       |              |             |             |
       |              |             |             |
       v              v             v             v
  +------------------------------------------------------+
  |                     OpenTracing                      |
  +------------------------------------------------------+
     |                |                |               |
     |                |                |               |
     v                v                v               v
+-----------+  +-------------+  +-------------+  +-----------+
|  Tracing  |  |   Logging   |  |   Metrics   |  |  Tracing  |
| System A  |  | Framework B |  | Framework C |  | System D  |
+-----------+  +-------------+  +-------------+  +-----------+
```

### OpenTracing 数据模型
一个 trace 代表了一个事务或者流程在（分布式）系统中的执行过程。在 OpenTracing 标准中，trace 是多个 span 组成的一个有向无环图（DAG），每一个 span 代表 trace 中被命名并计时的连续性的执行片段。

下图是一个分布式调用的例子，客户端发起请求，经过认证服务，计费服务，然后请求资源，最后返回结果。

![opentracing1.png](/pics/opentracing1.png)

可以使用包含时间轴的时序图来呈现这个 Trace

![opentracing2.png](/pics/opentracing2.png)

更多关于 OpenTracing 数据模型的知识，请参考 [OpenTracing语义标准](https://github.com/opentracing-contrib/opentracing-specification-zh/blob/master/specification.md)。

## Jaeger
[Jaeger](http://jaeger.readthedocs.io/en/latest/) 是 Uber 推出的一款开源分布式追踪系统，兼容 OpenTracing API。

![architecture-jaeger.png](/pics/architecture-jaeger.png)

如上图所示，Jaeger 主要由以下几部分组成。
* Jaeger Client - 为不同语言实现了符合 OpenTracing 标准的 SDK。应用程序通过 API 写入数据，client library 把 trace 信息按照应用程序指定的采样策略传递给 jaeger-agent。
* Agent - 它是一个监听在 UDP 端口上接收 span 数据的网络守护进程，它会将数据批量发送给 collector。它被设计成一个基础组件，部署到所有的宿主机上。Agent 将 client library 和 collector 解耦，为 client library 屏蔽了路由和发现 collector 的细节。
* Collector - 接收 jaeger-agent 发送来的数据，然后将数据写入后端存储。
* Data Store - 后端存储被设计成一个可插拔的组件，支持将数据写入 cassandra、elastic search。
* Query - 接收查询请求，然后从后端存储系统中检索 trace 并通过 UI 进行展示。

## Jaeger on Aliyun Log Service
[Jaeger on Aliyun Log Service](https://github.com/aliyun/jaeger) 是基于 Jeager 开发的分布式追踪系统，支持将采集到的追踪数据持久化到[日志服务](https://help.aliyun.com/product/28958.html)中，并通过 Jaeger 的原生接口进行查询和展示。

![architecture.png](/pics/architecture.png)

### 功能优势
* 原生 Jaeger 仅支持将数据持久化到 cassandra 和 elasticsearch 中，用户需要自行维护后端存储系统的稳定性，调节存储容量。Jaeger on Aliyun Log Service 借助阿里云日志服务的海量数据处理能力，让您享受 Jaeger 在分布式追踪领域给您带来便捷的同时无需过多关注后端存储系统的问题。
* Jaeger UI 部分仅提供查询、展示 trace 的功能，对分析问题、排查问题支持不足。使用 Jaeger on Aliyun Log Service，您可以借助日志服务强大的[查询分析](https://help.aliyun.com/document_detail/43772.html)能力，助您更快分析出系统中存在的问题。

### 配置步骤
参阅：https://github.com/aliyun/jaeger/blob/master/README_CN.md

### 使用实例
[HotROD](https://github.com/aliyun/jaeger/tree/master/examples/hotrod) 是由多个微服务组成的应用程序，它使用了 OpenTracing API 记录 trace 信息。

下面通过一段视频向您展示如何使用 Jaeger on Aliyun Log Service 诊断 HotROD 出现的问题。视频包含以下内容：
* 如何配置日志服务
* 如何通过 docker-compose 运行 Jaeger
* 如何运行 HotROD
* 如何根据查询条件检索特定的 trace
* 如何查看 trace 的详细信息
* 如何定位应用的性能瓶颈
* 应用程序如何使用 OpenTracing API

<video src="http://cloud.video.taobao.com//play/u/2143829456/p/1/e/6/t/1/50080498316.mp4" controls="true"></video>

[![Watch the video](/pics/jaeger_video.png)](http://cloud.video.taobao.com//play/u/2143829456/p/1/e/6/t/1/50080498316.mp4)

更多关于应用程序如何使用 OpenTracing API 将数据记录到 Jaeger，可参考链接：http://jaeger.readthedocs.io/en/latest/client_libraries/

## 参考资料
* Jaeger on Aliyun Log Service - https://github.com/aliyun/jaeger
* OpenTracing 中文文档 - https://wu-sheng.gitbooks.io/opentracing-io/content/
* Jaeger - http://jaeger.readthedocs.io/en/latest/getting_started/
* OpenTracing tutorial - https://github.com/yurishkuro/opentracing-tutorial

## 技术支持