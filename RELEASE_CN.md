# Aliyun LOG Jaeger 发布流程

## 更新版本
修改 [version.go](pkg/aliyunlog/config/version.go) 中的字段 `version` 至下一个版本。

## 发布 docker 镜像
1. 运行命令 `make docker-images-query-collector DOCKER_NAMESPACE=registry.cn-hangzhou.aliyuncs.com/jaegertracing DOCKER_TAG=<new_version>`，为 `jaeger-query` 和 `jaeger-collector` 构建镜像。
2. 运行命令 `make docker-push-query-collector DOCKER_NAMESPACE=registry.cn-hangzhou.aliyuncs.com/jaegertracing DOCKER_TAG=<new_version>` 
将前一步构建出来的镜像推送至 docker 仓库中。
3. 修改 [aliyunlog-jaeger-docker-compose.yml](docker-compose/aliyunlog-jaeger-docker-compose.yml) 中镜像的字段。

## 发布二进制包
1. 运行 `make build-all-linux`，`make build-all-windows`，`make build-all-darwin` 构建能在不同平台上运行的二进制包。

## 新建 Release
进入 [release](https://github.com/aliyun/aliyun-log-jaeger/releases) 页面，点击`Draft a new release`为该发布的版本创建一个 release。
`Tag version` 为 `<new_version>`，`Release title` 为 `Release <new_version>`。
