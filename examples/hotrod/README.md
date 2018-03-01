# Hot R.O.D. - Rides on Demand

[中文版README](/examples/hotrod/README_CN.md)

This is a demo application that consists of several microservices and illustrates
the use of the OpenTracing API. It can be run standalone, but requires Jaeger backend
to view the traces. A tutorial / walkthough is available:
* [yunqi article](/README.JAEGER.md)

## Features

* View request timeline & errors, understand how the app works
* Find sources of latency, lack of concurrency

## Running

### Run Jaeger Backend

An all-in-one Jaeger backend is packaged as a Docker container with in-memory storage.

```
docker run -d -p6831:6831/udp -p16686:16686 jaegertracing/all-in-one:latest
```

Jaeger UI can be accessed at http://localhost:16686.

### Run HotROD Application

```
go get github.com/jaegertracing/jaeger
cd $GOPATH/src/github.com/jaegertracing/jaeger
make install
cd examples/hotrod
go run ./main.go all
```

Then open http://127.0.0.1:8080

