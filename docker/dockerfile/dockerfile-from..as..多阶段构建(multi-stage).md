# dockerfile-from..as..多阶段构建(multi-stage)

参考文章

1. [generals-space/cni-terway](https://github.com/generals-space/cni-terway/blob/41d2d806e0ef0749f5a091cb70b8b32cfe1676ac/dockerfile)

以[generals-space/cni-terway](https://github.com/generals-space/cni-terway/blob/master/dockerfile)项目的 dockerfile 为例.

```dockerfile
## docker build --no-cache=true -f dockerfile -t registry.cn-hangzhou.aliyuncs.com/generals-kuber/cni-terway:1.1 .
########################################################
FROM golang:1.12 as builder
## docker镜像通用设置
LABEL author=general
LABEL email="generals.space@gmail.com"
## 环境变量, 使docker容器支持中文
ENV LANG C.UTF-8

WORKDIR /cni-terway
COPY . .
ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o terway ./cmd/pod
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cni-terway ./cmd/cni

########################################################
FROM generals/alpine
## docker镜像通用设置
LABEL author=general
LABEL email="generals.space@gmail.com"
## 环境变量, 使docker容器支持中文
ENV LANG C.UTF-8

COPY --from=builder /cni-terway/terway /
COPY --from=builder /cni-terway/cni-terway /
CMD ["/terway"]
```
