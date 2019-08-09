# docker run --label选项应用

参考文章

1. [官方文档 docker run](https://docs.docker.com/engine/reference/commandline/run/)
2. [那些让你看起来很牛逼的Docker使用技巧](https://www.jianshu.com/p/0231568ab335)

`docker run`的`--label`可以实现过滤的效果. 当然ta的本质并不是这样, 而是修改了容器的元数据, 等同于在dockerfile中的LABEL指令.

```
d run -it --name mq01 --label type=middleware generals/centos7 /bin/bash
```

在`inspect`时, 可以得到如下输出

```json
"Labels": {
    "type": "middleware",
    "author": "general",
    "email": "generals.space@gmail.com",
    "org.label-schema.build-date": "20180804",
    "org.label-schema.license": "GPLv2",
    "org.label-schema.name": "CentOS Base Image",
    "org.label-schema.schema-version": "1.0",
    "org.label-schema.vendor": "CentOS"
}
```

可以看到, 在容器启动时添加的`type`标签与在构建镜像时的`author`, `email`同级.

那么怎么使用呢? 在`docker ps`时需要使用`-f`过滤选项. 如下

```
$ d ps
CONTAINER ID        IMAGE               COMMAND                  CREATED             STATUS              PORTS                     NAMES
36e3b598f7e2        generals/centos7    "/bin/bash"              4 minutes ago       Up 4 minutes                                  nifty_hypatia
b286c483cd02        postgres            "docker-entrypoint.s…"   28 hours ago        Up 28 hours         0.0.0.0:15432->5432/tcp   wuhougit_postgres_1
08ebfc7fd4aa        node:8              "docker-entrypoint.s…"   28 hours ago        Up 28 hours                                   wuhougit_frontend_1
d8e8552135b6        generals/golang     "tail -f /etc/profile"   28 hours ago        Up 28 hours         0.0.0.0:5897->5897/tcp    wuhougit_backend_1
```

```
$ d ps -f label=type=middleware
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
36e3b598f7e2        generals/centos7    "/bin/bash"         5 minutes ago       Up 5 minutes                            mq01
```
