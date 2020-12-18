# arm平台下构建的镜像entrypoint格式异常的问题

参考文章

1. [一行代码的变更让我陷入无尽加班，Dockerfile的ENTRYPOINT的两种格式](https://www.pkslow.com/archives/docker-entrypoint-issue)
    - `entrypoints`两种格式: `exec`格式与`shell`格式
    - exec格式可以接受参数，而shell格式会忽略参数
    - shell格式相当于在前面还要再添加`/bin/sh -c`，所以app启动的进程ID不是1

参考文章1还是可以读一读的, 但是我遇到的问题有点不太一样.

## 场景描述

在 docker 18.03.6 下, 构建如下 dockerfile.

```dockerfile
FROM registry.cn-hangzhou.aliyuncs.com/generals-space/centos7:latest
ENTRYPOINT ["tail", "-f", "/etc/profile"]
```

使用`docker history`子命令查看镜像层历史信息, 会发现为`ENTRYPOINT ["tail" "-f" "/etc/profile"]`

```console
$ docker history 镜像名 --no-trunc
IMAGE                   CREATED             CREATED BY                                                      SIZE                COMMENT
sha256:4f47f76be97316   16 seconds ago      /bin/sh -c #(nop)  ENTRYPOINT ["tail" "-f" "/etc/profile"]      0B
```

然而我在 arm64v8 平台下(docker 19.03.7), 构建同样的`dockerfile`, 其镜像层历史如下

```console
$ docker history 镜像名 --no-trunc
IMAGE                   CREATED             CREATED BY                                                                              SIZE                COMMENT
sha256:4f47f76be97316   16 seconds ago      /bin/sh -c #(nop)  ENTRYPOINT ["/bin/sh", "-c", "[\"tail\" \"-f\" \"/etc/profile\"]"]    0B
```

然后在启动容器的时候, 报`[tail`命令不存在(有个左括号)...应该是arm平台下 docker 进程的 bug, 接下来就需要想办法规避这个问题了.

## 解决方法

### 1. shell 形式

我找到参考文章1, 尝试使用`shell`的格式重新构建.

```dockerfile
FROM registry.cn-hangzhou.aliyuncs.com/generals-space/centos7:latest
ENTRYPOINT tail -f /etc/profile
```

这次的结果为

```console
$ docker history 镜像名 --no-trunc
IMAGE                   CREATED             CREATED BY                                                                          SIZE                COMMENT
sha256:4f47f76be97316   16 seconds ago      /bin/sh -c #(nop)  ENTRYPOINT ["/bin/sh", "-c", "tail -f /etc/profile"]             0B
```

...这种是可以执行成功的, ~~问题也算是解决了~~.

大意了, 忘记`shell`格式的`ENTRYPOINT`后面没法跟`CMD`参数了(即使有写, 在容器内容查看的时候进程的命令行参数也不会包含`CMD`的内容的, 根本不生效). 

### 2. ENTRYPOINT /docker-entrypoint.sh

`ENTRYPOINT /docker-entrypoint.sh`的形式也不行, 这个脚本根本没办法接收到来自`CMD`, 或是`docker run`中接的命令行参数.

```console
$ ps -ef
UID     PID   PPID  C STIME TTY       TIME CMD
root      1      0  0 22:42 ?     00:00:00 /bin/sh /docker-entrypoint.sh
```

`docker-entrypoint.sh`脚本中, `echo $@`无法得到我想要的信息.

### 3. ENTRYPOINT ["/docker-entrypoint.sh"]

没想到最终用这种方式成功了.

```dockerfile
FROM registry.cn-hangzhou.aliyuncs.com/generals-space/centos7:latest
ENTRYPOINT ["/docker-entrypoint.sh"]
```

这个 dockerfile 构建的镜像内容如下


```console
$ docker history 镜像名 --no-trunc
IMAGE                   CREATED             CREATED BY                                                                          SIZE                COMMENT
sha256:4f47f76be97316   16 seconds ago      /bin/sh -c #(nop)  ENTRYPOINT ["/docker-entrypoint.sh"]             0B
```

倒是正常的, 目前只能怀疑, 是因为`ENTRYPOINT ["/docker-entrypoint.sh"]`后面的命令列表中超过了1个成功, 发生了转义吧...
