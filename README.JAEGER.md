# 开放分布式追踪（OpenTracing）入门与 Jaeger 实现

## 分布式系统的运维挑战
容器、Serverless 编程方式的诞生极大提升了软件交付与部署的效率。在架构的演化过程中，可以看到两个变化：

![change.png](/pics/change.png)

* 应用架构开始从单体系统逐步转变微服务，其中的业务逻辑随之而来就会变成微服务之间调用与请求。
* 资源角度来看，传统服务器这个物理单位也逐渐淡化，变成了看不见摸不到的虚拟资源模式。

从以上两个变化可以看到这种弹性、标准化的架构背后，原先运维与诊断的需求也变得越来越复杂。为了应对这种变化趋势，诞生一系列面向 DevOps 的诊断与分析系统，包括集中式日志系统（Logging），集中式度量系统（Metrics）和分布式追踪系统（Tracing）。

## Logging，Metrics 和 Tracing
Logging，Metrics 和 Tracing 有各自专注的部分。
* Logging - 用于记录离散的事件。例如，应用程序的调试信息或错误信息。它是我们诊断问题的依据。
* Metrics - 用于记录可聚合的数据。例如，队列的当前深度可被定义为一个度量值，在元素入队或出队时被更新；HTTP 请求个数可被定义为一个计数器，新请求到来时进行累加。
* Tracing - 用于记录请求范围内的信息。例如，一次远程方法调用的执行过程和耗时。它是我们排查系统性能问题的利器。

这三者也有相互重叠的部分，如下图所示。

![logging_metrics_tracing.png](/pics/logging_metrics_tracing.png)

通过上述信息，我们可以对已有系统进行分类。例如，Zipkin 专注于 tracing 领域；Prometheus 开始专注于 metrics，随着时间推移可能会集成更多的 tracing 功能，但不太可能深入 logging 领域； 阿里云日志服务开始专注于 logging 部分，但同时也不断地集成其他领域的特性到系统中来。

下面我们重点介绍下 tracing。

## Tracing 的诞生
Tracing 是在90年代就已出现的技术。但真正让该领域流行起来的还是源于 Google 的一篇论文"[Dapper, a Large-Scale Distributed Systems Tracing Infrastructure](https://research.google.com/pubs/pub36356.html)"，而另一篇论文"[Uncertainty in Aggregate Estimates from Sampled Distributed Traces](https://research.google.com/pubs/pub40378.html)"中则包含关于采样的更详细分析。论文发表后一批优秀的 Tracing 软件孕育而生，比较流行的有：
* Dapper(Google) : 各 tracer 的基础
* StackDriver Trace (Google)
* Zipkin(twitter)
* Appdash(golang)
* 鹰眼(taobao)
* X-ray(aws)

分布式追踪系统发展很快，种类繁多，但核心步骤一般有三个：代码埋点，数据存储、查询展示。

下图是一个分布式调用的例子，客户端发起请求，请求首先到达负载均衡器，接着经过认证服务，计费服务，然后请求资源，最后返回结果。

![opentracing1.png](/pics/opentracing1.png)

数据被采集存储后，分布式追踪系统一般会选择使用包含时间轴的时序图来呈现这个 Trace。

![opentracing2.png](/pics/opentracing2.png)

但在数据采集过程中，由于需要侵入用户代码，并且不同系统的 API 并不兼容，这就导致了如果您希望切换追踪系统，往往会带来较大改动。

## OpenTracing
为了解决不同的分布式追踪系统 API 不兼容的问题，诞生了 [OpenTracing](http://opentracing.io/) 规范。
OpenTracing 是一个轻量级的标准化层，它位于**应用程序/类库**和**追踪或日志分析程序**之间。
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

### OpenTracing 的优势
* OpenTracing 已进入 CNCF，正在为全球的分布式追踪，提供统一的概念和数据标准。
* OpenTracing 通过提供平台无关、厂商无关的 API，使得开发人员能够方便的添加（或更换）追踪系统的实现。

### OpenTracing 数据模型
OpenTracing 中的 **Trace**（调用链）通过归属于此调用链的 **Span** 来隐性的定义。
特别说明，一条 **Trace**（调用链）可以被认为是一个由多个 **Span** 组成的有向无环图（DAG图），**Span** 与 **Span** 的关系被命名为 **References**。

例如：下面的示例 **Trace** 就是由8个 **Span** 组成：

~~~
单个 Trace 中，span 间的因果关系


        [Span A]  ←←←(the root span)
            |
     +------+------+
     |             |
 [Span B]      [Span C] ←←←(Span C 是 Span A 的孩子节点, ChildOf)
     |             |
 [Span D]      +---+-------+
               |           |
           [Span E]    [Span F] >>> [Span G] >>> [Span H]
                                       ↑
                                       ↑
                                       ↑
                         (Span G 在 Span F 后被调用, FollowsFrom)

~~~

有些时候，使用下面这种，基于时间轴的时序图可以更好的展现 **Trace**（调用链）：

~~~
单个 Trace 中，span 间的时间关系


––|–––––––|–––––––|–––––––|–––––––|–––––––|–––––––|–––––––|–> time

 [Span A···················································]
   [Span B··············································]
      [Span D··········································]
    [Span C········································]
         [Span E·······]        [Span F··] [Span G··] [Span H··]
~~~

每个 **Span** 包含以下的状态:（译者注：由于这些状态会反映在 OpenTracing API 中，所以会保留部分英文说明）

- An operation name，操作名称
- A start timestamp，起始时间
- A finish timestamp，结束时间
- **Span Tag**，一组键值对构成的 Span 标签集合。键值对中，键必须为 string，值可以是字符串，布尔，或者数字类型。
- **Span Log**，一组 span 的日志集合。
  每次 log 操作包含一个键值对，以及一个时间戳。
  键值对中，键必须为 string，值可以是任意类型。
  但是需要注意，不是所有的支持 OpenTracing 的 Tracer，都需要支持所有的值类型。
- **SpanContext**，Span 上下文对象 (下面会详细说明)
- **References**(Span间关系)，相关的零个或者多个 Span（**Span** 间通过 **SpanContext** 建立这种关系）

每一个 **SpanContext** 包含以下状态：

- 任何一个 OpenTracing 的实现，都需要将当前调用链的状态（例如：trace 和 span 的 id），依赖一个独特的 Span 去跨进程边界传输
- **Baggage Items**，Trace 的随行数据，是一个键值对集合，它存在于 trace 中，也需要跨进程边界传输

更多关于 OpenTracing 数据模型的知识，请参考 [OpenTracing语义标准](https://github.com/opentracing-contrib/opentracing-specification-zh/blob/master/specification.md)。

### OpenTracing 实现
这篇[文档](http://opentracing.io/documentation/pages/supported-tracers.html)列出了所有 OpenTracing 实现。在这些实现中，比较流行的为 [Jaeger](http://jaeger.readthedocs.io/en/latest/) 和 [Zipkin](https://zipkin.io/)。

## Jaeger
[Jaeger](http://jaeger.readthedocs.io/en/latest/) 是 Uber 推出的一款开源分布式追踪系统，兼容 OpenTracing API。

### Jaeger 架构

![architecture-jaeger.png](/pics/architecture-jaeger.png)

如上图所示，Jaeger 主要由以下几部分组成。
* Jaeger Client - 为不同语言实现了符合 OpenTracing 标准的 SDK。应用程序通过 API 写入数据，client library 把 trace 信息按照应用程序指定的采样策略传递给 jaeger-agent。
* Agent - 它是一个监听在 UDP 端口上接收 span 数据的网络守护进程，它会将数据批量发送给 collector。它被设计成一个基础组件，部署到所有的宿主机上。Agent 将 client library 和 collector 解耦，为 client library 屏蔽了路由和发现 collector 的细节。
* Collector - 接收 jaeger-agent 发送来的数据，然后将数据写入后端存储。
* Data Store - 后端存储被设计成一个可插拔的组件，支持将数据写入 cassandra、elastic search。
* Query - 接收查询请求，然后从后端存储系统中检索 trace 并通过 UI 进行展示。

### Jaeger 存在的问题
* 需要架设并维护存储。
* UI比较薄弱，有一些个性化的分析需求无法快速满足（例如对比，统计延迟分布等）。

## Jaeger on Aliyun Log Service
[Jaeger on Aliyun Log Service](https://github.com/aliyun/jaeger) 是基于 Jeager 开发的分布式追踪系统，支持将采集到的追踪数据持久化到[日志服务](https://help.aliyun.com/product/28958.html)中，并通过 Jaeger 的原生接口进行查询和展示。

![architecture.png](/pics/architecture.png)

### 优势
* 原生 Jaeger 仅支持将数据持久化到 cassandra 和 elasticsearch 中，用户需要自行维护后端存储系统的稳定性，调节存储容量。Jaeger on Aliyun Log Service 借助阿里云日志服务的海量数据处理能力，让您享受 Jaeger 在分布式追踪领域给您带来便捷的同时无需过多关注后端存储系统的问题。
* Jaeger UI 部分仅提供查询、展示 trace 的功能，对分析问题、排查问题支持不足。使用 Jaeger on Aliyun Log Service，您可以借助日志服务强大的[查询分析](https://help.aliyun.com/document_detail/43772.html)能力，助您更快分析出系统中存在的问题。
* 相对于 Jaeger 使用 elasticsearch 作为后端存储，使用日志服务的好处是支持按量付费，成本仅为 elasticsearch 的13%。参阅[自建ELK vs 日志服务(SLS)全方位对比](https://yq.aliyun.com/articles/213693)

### 配置步骤
参阅：https://github.com/aliyun/jaeger/blob/master/README_CN.md

### 使用实例
[HotROD](https://github.com/aliyun/jaeger/tree/master/examples/hotrod) 是由多个微服务组成的应用程序，它使用了 OpenTracing API 记录 trace 信息。

下面通过一段视频向您展示如何使用 Jaeger on Aliyun Log Service 诊断 HotROD 出现的问题。视频包含以下内容：
1. 如何配置日志服务
2. 如何通过 docker-compose 运行 Jaeger
3. 如何运行 HotROD
4. 通过 Jaeger UI 如何检索特定的 trace
5. 通过 Jaeger UI 如何查看 trace 的详细信息
6. 通过 Jaeger UI 如何定位应用的性能瓶颈
7. 通过日志服务管理控制台，如何定位应用的性能瓶颈
8. 应用程序如何使用 OpenTracing API

<video src="http://cloud.video.taobao.com//play/u/2143829456/p/1/e/6/t/1/50080498316.mp4" controls="true"></video>

[![Watch the video](/pics/jaeger_video.png)](http://cloud.video.taobao.com//play/u/2143829456/p/1/e/6/t/1/50080498316.mp4)

视频中用到的查询分析样例
1. 以分钟为单位统计 `frontend` 服务的 `HTTP GET /dispatch` 操作的平均延迟以及请求个数。
```
process.serviceName: "frontend" and operationName: "HTTP GET /dispatch" |
select from_unixtime( __time__ - __time__ % 60) as time,
truncate(avg(duration)/1000/1000) as avg_duration_ms,
count(1) as count
group by __time__ - __time__ % 60 order by time desc limit 60
```

2. 比较两条 trace 各个操作的耗时
```
traceID: "trace1" or traceID: "trace2" |
select operationName,
(max(duration)-min(duration))/1000/1000 as duration_diff_ms
group by operationName
order by duration_diff_ms desc
```

3. 统计延迟大于 1.5s 的 trace 的 IP 情况
```
process.serviceName: "frontend" and operationName: "HTTP GET /dispatch" and duration > 1500000000 |
select "process.tags.ip" as IP,
truncate(avg(duration)/1000/1000) as avg_duration_ms,
count(1) as count
group by "process.tags.ip"
```

## 参考资料
* Jaeger on Aliyun Log Service - https://github.com/aliyun/jaeger
* OpenTracing 中文文档 - https://wu-sheng.gitbooks.io/opentracing-io/content/
* Jaeger - http://jaeger.readthedocs.io/en/latest/getting_started/
* OpenTracing tutorial - https://github.com/yurishkuro/opentracing-tutorial
* http://peter.bourgon.org/blog/2017/02/21/metrics-tracing-and-logging.html

## 特别感谢
Jaeger on Aliyun Log Service 是由阿里云团队和共创平台上的贡献者共同完成的。感谢 [@WPH95](https://github.com/WPH95) 的杰出工作。

## 技术支持