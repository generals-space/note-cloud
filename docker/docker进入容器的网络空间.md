# docker进入容器的网络空间

参考文章

1. [从宿主机直接进入docker容器的网络空间](https://www.cnblogs.com/549294286/p/10832711.html)

docker有`exec`子命令, 可以直接进入容器, 那为什么还要有这样的需求呢?

主要是因为, 很多超精简的镜像(比如`google_containers/coredns:1.3.1`)内部没有提供`/bin/sh`或`/bin/bash`, 无法进入交互式命令行, 更别说curl, ss等调试工具了. 如果我们的问题是在网络层面的话, 那么进入容器的网络空间, 而直接使用宿主机上的命令会非常方便.

首先查看目标容器在宿主机上映射的pid

```
docker inspect -f '{{.State.Pid}}' 容器名或ID
```

然后根据得到的pid进入容器

```
nsenter -t 容器进程pid -n /bin/sh
```
