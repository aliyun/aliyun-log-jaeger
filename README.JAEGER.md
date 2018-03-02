# 分布式追踪利器 Jaeger 上云实战

## 分布式追踪面临的挑战
近年来，随着微服务架构的兴起，一些在单机环境下容易处理的问题变得越来越困难：优化接口调用的延迟、分析后端系统的错误根源、找出系统的性能瓶颈、分析分布式系统中各组件的调用情况等。

## Jaeger on Aliyun Log Service
[Jaeger](http://jaeger.readthedocs.io/en/latest/) 是 Uber 推出的一款开源分布式追踪系统，为微服务场景而生。它主要用于分析多个服务的调用过程，图形化服务调用轨迹，是诊断性能问题、分析系统故障的利器。

[Jaeger on Aliyun Log Service](https://github.com/aliyun/jaeger) 是基于 Jeager 开发的分布式追踪系统，支持将采集到的追踪数据持久化到[日志服务](https://help.aliyun.com/product/28958.html)中，并通过 Jaeger 的原生接口进行查询和展示。

![architecture.png](/pics/architecture.png)

## 功能优势
原生 Jaeger 仅支持将数据持久化到 cassandra 和 elasticsearch 中，用户需要自行维护后端存储系统的稳定性，调节存储容量。Jaeger on Aliyun Log Service 借助阿里云日志服务的海量数据处理能力，让您享受 Jaeger 在分布式追踪领域给您带来便捷的同时无需过多关注后端存储系统的问题。

## 配置步骤
参阅：https://github.com/aliyun/jaeger/blob/master/README_CN.md

## 使用实例
[HotROD](https://github.com/aliyun/jaeger/tree/master/examples/hotrod) 是由多个微服务组成的应用程序，它使用了 OpenTracing API 记录 trace 信息。

下面通过一段视频向您展示如何使用 Jaeger on Aliyun Log Service 诊断 HotROD 出现的问题。视频包含以下内容：
* 如何配置日志服务
* 如何通过 docker-compose 运行 Jaeger
* 如何运行 HotROD
* 如何根据查询条件检索特定的 trace
* 如何查看 trace 的详细信息
* 如何定位应用的性能瓶颈
* 应用程序如何使用 OpenTracing API

<video autoplay="autoplay" src="http://cloud.video.taobao.com//play/u/2143829456/p/1/e/6/t/1/50080498316.mp4" controls="true"></video>

[![Watch the video](/pics/jaeger_video.png)](http://cloud.video.taobao.com//play/u/2143829456/p/1/e/6/t/1/50080498316.mp4)

更多关于应用程序如何使用 OpenTracing API 将数据记录到 Jaeger，可参考链接：http://jaeger.readthedocs.io/en/latest/client_libraries/

## 技术支持