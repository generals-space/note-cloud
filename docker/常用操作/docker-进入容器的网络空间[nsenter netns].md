# docker-进入容器的网络空间

参考文章

1. [从宿主机直接进入docker容器的网络空间](https://www.cnblogs.com/549294286/p/10832711.html)
2. [Docker容器进入的4种方式](https://www.cnblogs.com/xhyan/p/6593075.html)
    - `nsenter`的使用
3. [ip netns命令详解](https://blog.csdn.net/supahero/article/details/100606953)
4. [docker networking namespace not visible in ip netns list](https://stackoverflow.com/questions/31265993/docker-networking-namespace-not-visible-in-ip-netns-list)

## 1. nsenter

docker有`exec`子命令, 可以直接进入容器, 那为什么还要有这样的需求呢?

主要是因为, 很多超精简的镜像(比如`google_containers/coredns:1.3.1`)内部没有提供`/bin/sh`或`/bin/bash`, 无法进入交互式命令行, 更别说curl, ss等调试工具了. 如果我们的问题是在网络层面的话, 那么进入容器的网络空间, 而直接使用宿主机上的命令会非常方便.

首先查看目标容器在宿主机上映射的pid

```bash
docker inspect -f '{{.State.Pid}}' 容器名或ID
```

然后根据得到的pid进入容器

```bash
nsenter -t 容器进程pid --net /bin/sh
```

需要注意的时, 由于进入的是容器的网络空间, 但是进程列表却是与宿主机完全相通的, 在排查进程相关的问题时可能还是会比较棘手.

但是如果再加上`--mount`和`--pid`两个参数后(单纯加`--pid`还是会显示宿主机的进程列表, 而单纯加`--mount`则无法使用`ps -ef`命令, 显示`Error, do this: mount -t proc proc /proc`), 虽然只显示容器内的进程了, 但是又没有办法使用宿主机上的命令了...很难两全.

## 2. ip netns

使用`ip netns add ns01`创建的`netns`会出现在`/var/run/netns`目录下(默认这个目录不存在, 在第一次创建`netns`时自动创建).

```log
$ ip netns add ns01
$ cd /var/run/netns/
$ ls
ns01
```

但 docker 创建出来的容器的 netns 并不在`/var/run/netns`, 而是在`/proc/[pid]/ns/net`, 其中容器内部的进程pid可以通过如下命令查看.

```log
$ docker inspect ad7a7a5952fd | grep Pid
            "Pid": 1843,
```

然后可以使用如下命令建立一个软链接再进行访问.

```log
$ ip netns ls
$ ln -s /proc/1843/ns/net /var/run/netns/nginx
$ ip netns ls
nginx
$ ip netns exec nginx sh
sh-4.2# ls
```
