`docker buildx build`构建的镜像, 无法通过`docker images`查到, 一般只会直接推送到远程仓库, 然后在`pull`时只拉取客户端所在主机架构的layer层.

## 构建同时推送到远程仓库

```
docker buildx build --platform linux/arm64,linux/amd64 -t hello-go:v0.0.1 --output=type=registry .
```

> `--output=type=registry`与`-push`等价, 构建完成后, 会自动推送到远程仓库, 这要求`-t`指定的镜像地址是合法的, 否则推送会失败.

## 

```
docker buildx build --platform linux/arm64,linux/amd64 -t hello-go:v0.0.1 --output=type=tar,dest=./hello-go.tar .
```

这样生成的tar包并不是常规合法镜像tar包, 无法被`import`

```log
$ docker load < hello-go.tar
open /var/lib/docker/tmp/docker-import-2688680514/linux_amd64/json: no such file or directory
```
