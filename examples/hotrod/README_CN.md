# Hot R.O.D. - 出行需求

这个演示程序由多个微服务组成，它向我们展示了如何使用 OpenTracing API。它可以独立运行，但如果想要观察 trace 必须运行 jaeger。您可以参考下面的教程：
  * [云栖文章](/README.LOCAL.md)

## 功能

* 查看请求时间轴和错误，了解应用程序是如何工作的
* 查找延迟、并发量低的根源

# 运行

## 运行 Jaeger

进入包含文件 [aliyunlog-jaeger-docker-compose.yml](/docker-compose/aliyunlog-jaeger-docker-compose.yml) 的目录，将参数 ${PROJECT}、${ENDPOINT}、${ACCESS_KEY_ID}、${ACCESS_KEY_SECRET}、${SPAN_LOGSTORE} 替换为真实值，然后运行下列命令。

```
docker-compose -f aliyunlog-jaeger-docker-compose.yml up
```

在浏览器中打开 http://127.0.0.1:16686/ 访问 Jaeger UI。

## 运行 HotROD 应用

```
mkdir -p $GOPATH/src/github.com/jaegertracing
cd $GOPATH/src/github.com/jaegertracing
git clone https://github.com/aliyun/jaeger.git jaeger
cd jaeger
make install
cd examples/hotrod
go run ./main.go all
```

在浏览器中打开 http://127.0.0.1:8080/ 。
