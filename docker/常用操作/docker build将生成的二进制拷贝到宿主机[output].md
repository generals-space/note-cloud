# docker build将生成的二进制拷贝到宿主机

Server Version: 26.1.4

在某种情况下(比如宿主机没有构建环境, 如 golang), 我们希望用官方镜像通过 docker build 构建二进制程序, 直接生成或是拷贝到宿主机.

一般`docker build`都是通过`-t`参数指定生成一个镜像, 我们可以先用这个镜像启动一个容器, 然后用 docker cp 将其拷贝出来.

但还有另一种更方便的方式.

有一个与`-t`平级的参数`--output`, 可以将生成镜像中的所有文件输出到目标路径下(只有文件系统, 不包含层信息, 如 meta 文件, 哈希目录什么的).

而且整个镜像的文件系统还是太多了, 可以通过多阶段构建进行优化, 只输出二进制文件, 不要其他目录.

```dockerfile
FROM golang:1.22 as builder

WORKDIR /mytest
COPY . .
ENV GO111MODULE on
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o mytest ./main.go

################################################################################
FROM scratch AS output

COPY --from=builder /mytest/mytest /
```

> `scratch`是一个空镜像, 第2阶段拷贝到根目录后, 整个文件系统就只有这一个文件, 直接输出即可.

```
docker build --output type=local,dest=bin .
```
