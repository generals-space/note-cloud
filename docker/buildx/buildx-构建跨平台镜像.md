参考文章

1. [buildx 构建跨平台镜像](https://zhuanlan.zhihu.com/p/622399482)
    - 一篇足够
    - 「跨平台镜像」这一说法并不十分准确, 实际上, Docker 官方术语叫 Multi-platform images, 即「多平台镜像」, 意思是支持多种不同 CPU 架构的镜像.
2. [Docker buildx question, copy files based on arch](https://www.reddit.com/r/docker/comments/sh5c8v/docker_buildx_question_copy_files_based_on_arch/)
3. [Build multi-arch images with different commands per architecture in Docker file](https://forums.docker.com/t/build-multi-arch-images-with-different-commands-per-architecture-in-docker-file/134795)

`buildx`可以构建跨平台的镜像, 大致原理就是, 在构建时将多个镜像的layer层, 合并到一个镜像里保存及推送, 然后在`docker pull`时, 客户端只拉取当前主机所属架构的layer层.

要安装并使用 buildx，需要 Docker Engine 版本号大于等于 19.03, 启用方法.

```
docker buildx create --name mybuilder
docker buildx use mybuilder
docker buildx inspect --bootstrap mybuilder
```

## 构建方式1 - 常规 Dockerfile

假设如下 golang 程序.

```go
// main.go
package main

import (
    "fmt"
    "runtime"
)

func main() {
    fmt.Printf("Hello, %s/%s!\n", runtime.GOOS, runtime.GOARCH)
}
```

```go
// go.mod
module hello

go 1.20
```

编写 Dockerfile 内容如下

```dockerfile
FROM golang:1.20-alpine AS builder
WORKDIR /app
ADD . .
RUN go build -o hello .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/hello .
CMD ["./hello"]
```

上面的 Dockerfile 与普通场景完全相同, 使用如下命令进行构建

```
docker buildx build --platform linux/arm64,linux/amd64 -t hello-go:v0.0.1 .
```

这样就构建出了兼容2种CPU架构的镜像, 其实相当于分别在arm64/amd64环境中, 都执行了一次`docker build`, 因此`Dockerfile`不显式指定`GOARCH`也没关系.

> 至于模拟arm64/amd64的手段, 是通过`QEMU`虚拟机实现的.

## 构建方式2 - 显示指定平台变量, 条件判断

将上述`Dockerfile`写成下面这种.

```dockerfile
## 拉取适配于当前主机架构的基础镜像
FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
ADD . .
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o hello .

## 跨平台镜像要根据目标平台拉取基础镜像
FROM --platform=$TARGETPLATFORM alpine:latest
WORKDIR /app
COPY --from=builder /app/hello .
CMD ["./hello"]
```

其中`BUILDPLATFORM`、`TARGETOS`、`TARGETARCH`、`TARGETPLATFORM`四个变量是 BuildKit 提供的全局变量，分别表示构建镜像所在平台、操作系统、架构、构建镜像的目标平台。

在构建镜像时，BuildKit 会将当前所在平台信息传递给 Dockerfile 中的 BUILDPLATFORM 参数（如 linux/amd64）。

通过`--platform`参数传递的 linux/arm64,linux/amd64 镜像目标平台列表会依次传递给`TARGETPLATFORM`变量。

而`TARGETOS`、`TARGETARCH`两个变量在使用时则需要先通过`ARG`指令进行声明，BuildKit 会自动为其赋值。

## 3.

在很多场景下, 不同平台的镜像中执行的命令是有区别的. 另外, 对于一个第三方的工具, 需要在不同架构的镜像中, 拷贝对应架构的二进制程序, 这种需求, 如何实现?

```dockerfile
FROM --platform=$TARGETPLATFORM alpine:latest

COPY xxx.amd64 /usr/local/bin/xxx.amd64
COPY xxx.arm64 /usr/local/bin/xxx.arm64

## 接收 buildx 构建时传入的参数
ARG TARGETPLATFORM

RUN chmod 755 -R /usr/local/bin/; \
echo $TARGETPLATFORM; \
if [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
    mv /usr/local/bin/xxx.arm64 /usr/local/bin/xxx; \
elif [ "$TARGETPLATFORM" = "linux/amd64" ]; then \
    mv /usr/local/bin/xxx.amd64 /usr/local/bin/xxx; \
else \
    echo invalid $TARGETPLATFORM; \
fi
```

> 这种情况的构建命令最好加上`--progress`选项, 这可以在构建时打印出标准输出的内容, 这样可以判断是否进入正确的`if`块.
