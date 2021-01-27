# docker-进入容器的网络空间

参考文章

1. [从宿主机直接进入docker容器的网络空间](https://www.cnblogs.com/549294286/p/10832711.html)
2. [Docker容器进入的4种方式](https://www.cnblogs.com/xhyan/p/6593075.html)
    - `nsenter`的使用

docker有`exec`子命令, 可以直接进入容器, 那为什么还要有这样的需求呢?

主要是因为, 很多超精简的镜像(比如`google_containers/coredns:1.3.1`)内部没有提供`/bin/sh`或`/bin/bash`, 无法进入交互式命令行, 更别说curl, ss等调试工具了. 如果我们的问题是在网络层面的话, 那么进入容器的网络空间, 而直接使用宿主机上的命令会非常方便.

首先查看目标容器在宿主机上映射的pid

```
docker inspect -f '{{.State.Pid}}' 容器名或ID
```

然后根据得到的pid进入容器

```
nsenter -t 容器进程pid --net /bin/sh
```

需要注意的时, 由于进入的是容器的网络空间, 但是进程列表却是与宿主机完全相通的, 在排查进程相关的问题时可能还是会比较棘手.

但是如果再加上`--mount`和`--pid`两个参数后(单纯加`--pid`还是会显示宿主机的进程列表, 而单纯加`--mount`则无法使用`ps -ef`命令, 显示`Error, do this: mount -t proc proc /proc`), 虽然只显示容器内的进程了, 但是又没有办法使用宿主机上的命令了...很难两全.
